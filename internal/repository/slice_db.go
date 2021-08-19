package repository

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mehrdod/purehttpserver/pkg/logger"
	"os"
	"time"
)

type SliceDb struct {
	lastRequests []int
	firstReqTime time.Time
	counterNum   int

	backUpFileName string
}

func NewSliceDb(counterNum int, backupFileName string) (*SliceDb, error) {
	var db SliceDb
	db.backUpFileName = backupFileName
	err := db.RecoverState(counterNum)
	if err == nil {
		logger.Info("DB restored")
		return &db, nil
	}
	logger.Info("new db created, reason: %v", err)
	db.lastRequests = make([]int, counterNum)
	db.firstReqTime = time.Now().Add(time.Duration(-2*counterNum) * time.Second)
	db.counterNum = counterNum
	return &db, nil
}

func (db *SliceDb) RecoverState(counterNum int) error {
	file, err := os.Open(db.backUpFileName)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	var fistReqTime int64

	_, err = fmt.Fscanf(reader, "%d\n", &fistReqTime)
	if err != nil {
		return err
	}
	db.firstReqTime = time.Unix(fistReqTime, 0)
	_, err = fmt.Fscanf(reader, "%d", &db.counterNum)
	if err != nil {
		return err
	}

	if counterNum != db.counterNum {
		return errors.New("config counterNum and backuped should match")
	}
	db.lastRequests = make([]int, counterNum)

	for i := 0; i < counterNum; i++ {
		_, err = fmt.Fscanf(reader, "%d", &db.lastRequests[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *SliceDb) SaveState() error {
	file, err := os.Create(db.backUpFileName)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "%d\n", db.firstReqTime.Unix())
	fmt.Fprintf(writer, "%d ", db.counterNum)
	for _, v := range db.lastRequests {
		fmt.Fprintf(writer, "%d ", v)
	}

	return nil
}
