package port

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
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

type ArticleRepository interface {
	CreateArticle(context.Context, CreateArticleParams) (domain.Article, error)
	AddTagsIfNotExists(context.Context, AddTagsParams) ([]domain.Tag, error)
	AssignTags(context.Context, AssignTags) ([]domain.ArticleTag, error)
}
