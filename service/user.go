package service

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/port"
)

type userService struct {
	repo port.Repository
}

func NewUserService(repo port.Repository) port.UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, req port.CreateUserTxParams) (result port.CreateUserTxResult, err error) {
	result.User, err = s.repo.User().CreateUser(ctx, port.CreateUserParams{User: req.User})
	if err != nil {
		return port.CreateUserTxResult{}, err
	}
	return result, nil
}
