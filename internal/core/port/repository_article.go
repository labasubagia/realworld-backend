package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type FilterArticlePayload struct {
	Slugs     []string
	IDs       []domain.ID
	AuthorIDs []domain.ID
	Limit     int
	Offset    int
}

type AddTagsPayload struct {
	Tags []string
}

type AssignTagPayload struct {
	ArticleID domain.ID
	TagIDs    []domain.ID
}

type FilterTagPayload struct {
	IDs   []domain.ID
	Names []string
}

type FilterArticleTagPayload struct {
	ArticleIDs []domain.ID
	TagIDs     []domain.ID
}

type FilterFavoritePayload struct {
	UserIDs    []domain.ID
	ArticleIDs []domain.ID
}

type FilterCommentPayload struct {
	ArticleIDs []domain.ID
	AuthorIDs  []domain.ID
}

type ArticleRepository interface {
	CreateArticle(context.Context, domain.Article) (domain.Article, error)
	UpdateArticle(context.Context, domain.Article) (domain.Article, error)
	DeleteArticle(context.Context, domain.Article) error
	FilterArticle(context.Context, FilterArticlePayload) ([]domain.Article, error)
	FindOneArticle(context.Context, FilterArticlePayload) (domain.Article, error)

	FilterTags(context.Context, FilterTagPayload) ([]domain.Tag, error)
	AddTags(context.Context, AddTagsPayload) ([]domain.Tag, error)

	FilterArticleTags(context.Context, FilterArticleTagPayload) ([]domain.ArticleTag, error)
	AssignArticleTags(context.Context, AssignTagPayload) ([]domain.ArticleTag, error)

	AddFavorite(context.Context, domain.ArticleFavorite) (domain.ArticleFavorite, error)
	RemoveFavorite(context.Context, domain.ArticleFavorite) (domain.ArticleFavorite, error)
	FilterFavorite(context.Context, FilterFavoritePayload) ([]domain.ArticleFavorite, error)
	FilterFavoriteCount(context.Context, FilterFavoritePayload) ([]domain.ArticleFavoriteCount, error)

	AddComment(context.Context, domain.Comment) (domain.Comment, error)
	DeleteComment(context.Context, domain.Comment) error
	FilterComment(context.Context, FilterCommentPayload) ([]domain.Comment, error)
}
