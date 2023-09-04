package port

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
)

type Service interface {
	User(repo Repository) UserService
	Article(repo Repository) ArticleService
}

type CreateUserTxParams struct {
	User      domain.User
	AfterFunc func(domain.User) error
}

type CreateUserTxResult struct {
	User domain.User
}

type UserService interface {
	Create(context.Context, CreateArticleParams) (CreateArticleTxResult, error)
}

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
