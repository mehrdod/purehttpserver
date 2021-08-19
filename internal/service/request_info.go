package service

import (
	"github.com/mehrdod/purehttpserver/domain"
	"github.com/mehrdod/purehttpserver/internal/repository"
	"time"
)

type RequestInfoService struct {
	repo repository.RequestInfo
}

func NewRequestInfoService(repo repository.RequestInfo) *RequestInfoService {
	return &RequestInfoService{repo: repo}
}

func (s *RequestInfoService) Get() (domain.RequestInfo, error) {
	err := s.repo.Increment(time.Now())
	if err != nil {
		return domain.RequestInfo{}, err
	}
	return s.repo.Get()
}

