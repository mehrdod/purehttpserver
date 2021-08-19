package repository

import (
	"errors"
	"github.com/mehrdod/purehttpserver/domain"
	"time"
)

type RequestInfoRepo struct {
	db *SliceDb
}

func NewRequestInfoRepo(db *SliceDb) *RequestInfoRepo {
	return &RequestInfoRepo{db: db}
}

func (r *RequestInfoRepo) Get() (domain.RequestInfo, error) {
	var reqInfo domain.RequestInfo
	if len(r.db.lastRequests) != r.db.counterNum {
		return reqInfo, errors.New("something wrong with slice db")
	}
	reqInfo.LastRequestsNum = 0
	for i := 0; i < len(r.db.lastRequests); i++ {
		reqInfo.LastRequestsNum += r.db.lastRequests[i]
	}

	return reqInfo, nil
}

func (r *RequestInfoRepo) Increment(curReqTime time.Time) error {
	if len(r.db.lastRequests) != r.db.counterNum {
		return errors.New("something wrong with slice db")
	}

	curReqTime = curReqTime.Round(time.Second)
	lastReqTime := r.db.firstReqTime.Add(time.Duration(r.db.counterNum) * time.Second)

	if lastReqTime.Before(curReqTime) {
		for i := 0; i < len(r.db.lastRequests); i++ {
			r.db.lastRequests[i] = 0
		}
		r.db.lastRequests[0] = 1
		r.db.firstReqTime = curReqTime
		return nil
	}

	for !r.db.firstReqTime.Equal(curReqTime) {
		r.db.firstReqTime = r.db.firstReqTime.Add(time.Second)
		r.db.lastRequests = r.db.lastRequests[1:]
		r.db.lastRequests = append(r.db.lastRequests, 0)
	}
	length := len(r.db.lastRequests)
	r.db.lastRequests[length-1]++
	return nil
}
