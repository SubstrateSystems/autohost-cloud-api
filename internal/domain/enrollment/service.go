package enrollment

import (
	"errors"
	"time"
)

var (
	ErrEnrollTokenNotFound    = errors.New("enroll token not found")
	ErrInvalidEnrollTokenData = errors.New("invalid enroll token data")
	ErrUnauthorized           = errors.New("unauthorized")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEnrollToken(token string, userID string, expiresAt time.Time) error {
	if token == "" || userID == "" {
		return ErrInvalidEnrollTokenData
	}
	return s.repo.CreateEnrollToken(token, userID, expiresAt)
}
