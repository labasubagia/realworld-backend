package service

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/port"
)

type articleService struct {
	repo port.Repository
}

func NewArticleService(repo port.Repository) port.ArticleService {
	return &articleService{
		repo: repo,
	}
}

// Create implements port.ArticleService.
func (s *articleService) Create(context.Context, port.CreateArticleTxParams) (port.CreateArticleTxResult, error) {

	return port.CreateArticleTxResult{}, nil
}
