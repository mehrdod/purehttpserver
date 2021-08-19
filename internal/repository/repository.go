package repository

import (
	"github.com/mehrdod/purehttpserver/domain"
	"time"
)

type RequestInfo interface {
	Get() (domain.RequestInfo, error)
	Increment(reqTime time.Time) error
}

type Repositories struct {
	RequestInfo RequestInfo
}

func NewRepositories(db *SliceDb) *Repositories {
	return &Repositories{RequestInfo: NewRequestInfoRepo(db)}
}
