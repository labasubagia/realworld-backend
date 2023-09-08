package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type CreateArticleTxParams struct {
	Article domain.Article
	Tags    []string
}

type CreateArticleTxResult struct {
	Article domain.Article
	Tags    []domain.Tag
}

type ArticleService interface {
	Create(context.Context, CreateArticleTxParams) (CreateArticleTxResult, error)
}
