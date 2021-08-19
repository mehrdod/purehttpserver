package app

import (
	"context"
	"errors"
	delivery "github.com/mehrdod/purehttpserver/internal/delivery/http"
	"github.com/mehrdod/purehttpserver/internal/repository"
	"github.com/mehrdod/purehttpserver/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mehrdod/purehttpserver/internal/config"

	"github.com/mehrdod/purehttpserver/internal/server"
	"github.com/mehrdod/purehttpserver/pkg/logger"
)

func Run(configPath string) {
	//Initialize configs
	cfg, err := config.Init(configPath, "dev")
	if err != nil {
		logger.Errorf("failed to init config %v", err)
		return
	}
	db, err := repository.NewSliceDb(cfg.Other.RequestCounterTTL, cfg.Other.BackUpFileName)
	if err != nil {
		logger.Errorf("failed to init db %v", err)
		return
	}
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos)

	handlers := delivery.NewHandler(services, cfg.Other.RequestCounterTTL)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init())

	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
	if err := db.SaveState(); err != nil {
		logger.Errorf("failed to save state of db: %v", err)
	}
	logger.Info("Server stopped and saved the state")
}
