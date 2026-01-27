package users

import (
	"context"
)

type Service interface {
	ListUsers(ctx context.Context) ([]string, error)
}

func NewService() Service { //
	return &svc{}
}

type svc struct {
	// Repository
}

func (s *svc) ListUsers(ctx context.Context) ([]string, error) {
	// Success case
	users := []string{"user1", "user2"}
	return users, nil
}
