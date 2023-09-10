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

type ListArticleParams struct {
	AuthArg        AuthParams
	IDs            []domain.ID
	Tags           []string
	AuthorNames    []string
	FavoritedNames []string
	Limit          int
	Offset         int
}

type ListArticleResult struct {
	Articles []domain.Article
	Count    int
}

type AddFavoriteParams struct {
	AuthArg AuthParams
	Slug    string
	UserID  domain.ID
}

type AddFavoriteResult struct {
	Article domain.Article
}

type ArticleService interface {
	Create(context.Context, CreateArticleTxParams) (CreateArticleTxResult, error)
	List(context.Context, ListArticleParams) (ListArticleResult, error)
	AddFavorite(context.Context, AddFavoriteParams) (AddFavoriteResult, error)
}
