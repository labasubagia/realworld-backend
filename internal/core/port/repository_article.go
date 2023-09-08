package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type CreateArticleParams struct {
	Article domain.Article
}

type AddTagsParams struct {
	Tags []string
}

type AssignTags struct {
	ArticleID int64
	TagIDs    []int64
}

type FilterTagParams struct {
	IDs   []int64
	Names []string
}

type ArticleRepository interface {
	CreateArticle(context.Context, CreateArticleParams) (domain.Article, error)
	FilterTags(context.Context, FilterTagParams) ([]domain.Tag, error)
	AddTagsIfNotExists(context.Context, AddTagsParams) ([]domain.Tag, error)
	AssignTags(context.Context, AssignTags) ([]domain.ArticleTag, error)
}
