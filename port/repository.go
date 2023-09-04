package port

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
)

type RepositoryAtomicCallback func(r Repository) error

type Repository interface {
	Atomic(context.Context, RepositoryAtomicCallback) error
	User() UserRepository
	Article() ArticleRepository
}

type CreateUserParams struct {
	User domain.User
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserParams) (domain.User, error)
}

type CreateArticleParams struct {
	Article domain.Article
}

type CreateTagParams struct {
	ArticleID string
	Tag       string
}

type ArticleRepository interface {
	CreateArticle(context.Context, CreateArticleParams) (domain.Article, error)
	CreateTag(context.Context, CreateTagParams) (domain.Tag, error)
}
