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

type UpdateArticleParams struct {
	AuthArg AuthParams
	Slug    string
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

type AddFavoriteParams struct {
	AuthArg AuthParams
	Slug    string
	UserID  domain.ID
}

type RemoveFavoriteParams AddFavoriteParams

type GetArticleParams struct {
	AuthArg AuthParams
	Slug    string
}

type AddCommentParams struct {
	AuthArg AuthParams
	Slug    string
	Comment domain.Comment
}

type ListCommentParams struct {
	AuthArg AuthParams
	Slug    string
}

type DeleteCommentParams struct {
	AuthArg   AuthParams
	Slug      string
	CommentID domain.ID
}

type ArticleService interface {
	Create(context.Context, CreateArticleTxParams) (domain.Article, error)
	Update(context.Context, UpdateArticleParams) (domain.Article, error)
	Delete(context.Context, DeleteArticleParams) error
	List(context.Context, ListArticleParams) ([]domain.Article, error)
	Feed(context.Context, ListArticleParams) ([]domain.Article, error)
	Get(context.Context, GetArticleParams) (domain.Article, error)

	AddComment(context.Context, AddCommentParams) (domain.Comment, error)
	ListComments(context.Context, ListCommentParams) ([]domain.Comment, error)
	DeleteComment(context.Context, DeleteCommentParams) error

	AddFavorite(context.Context, AddFavoriteParams) (domain.Article, error)
	RemoveFavorite(context.Context, RemoveFavoriteParams) (domain.Article, error)

	ListTags(context.Context) ([]string, error)
}
