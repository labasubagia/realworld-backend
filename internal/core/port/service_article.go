package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type CreateArticleTxParams struct {
	AuthArg AuthParams
	Article domain.Article
	Tags    []string
}

type CreateArticleTxResult struct {
	Article domain.Article
	Tags    []domain.Tag
}

type UpdateArticleParams struct {
	AuthArg AuthParams
	Slug    string
	Article domain.Article
}

type UpdateArticleResult struct {
	Article domain.Article
}

type DeleteArticleParams struct {
	AuthArg AuthParams
	Slug    string
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

type GetArticleParams struct {
	AuthArg AuthParams
	Slug    string
}

type GetArticleResult struct {
	Article domain.Article
}

type AddCommentParams struct {
	AuthArg AuthParams
	Slug    string
	Comment domain.Comment
}

type AddCommentResult struct {
	Comment domain.Comment
}

type ListCommentParams struct {
	AuthArg AuthParams
	Slug    string
}

type ListCommentResult struct {
	Comments []domain.Comment
}

type ArticleService interface {
	Create(context.Context, CreateArticleTxParams) (CreateArticleTxResult, error)
	Update(context.Context, UpdateArticleParams) (UpdateArticleResult, error)
	Delete(context.Context, DeleteArticleParams) error
	List(context.Context, ListArticleParams) (ListArticleResult, error)
	Feed(context.Context, ListArticleParams) (ListArticleResult, error)
	Get(context.Context, GetArticleParams) (GetArticleResult, error)

	AddComment(context.Context, AddCommentParams) (AddCommentResult, error)
	ListComments(context.Context, ListCommentParams) (ListCommentResult, error)

	AddFavorite(context.Context, AddFavoriteParams) (AddFavoriteResult, error)
}
