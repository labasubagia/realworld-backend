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

func (*userService) Create(context.Context, port.CreateArticleParams) (port.CreateArticleTxResult, error) {
	return port.CreateArticleTxResult{}, nil
}
