package service

import (
	"github.com/mehrdod/purehttpserver/domain"
	"github.com/mehrdod/purehttpserver/internal/repository"
)

type RequestInfo interface {
	Get() (domain.RequestInfo, error)
}

type Services struct {
	RequestInfo RequestInfo
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		RequestInfo: NewRequestInfoService(repos.RequestInfo),
	}
}

