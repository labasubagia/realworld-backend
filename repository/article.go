package repository

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
	"github.com/labasubagia/go-backend-realworld/port"
	"github.com/uptrace/bun"
)

type articleRepo struct {
	db bun.IDB
}

func NewArticleRepository(db bun.IDB) port.ArticleRepository {
	return &articleRepo{
		db: db,
	}
}

func (r *articleRepo) CreateArticle(context.Context, port.CreateArticleParams) (domain.Article, error) {
	return domain.Article{}, nil
}

func (r *articleRepo) CreateTag(context.Context, port.CreateTagParams) (domain.Tag, error) {
	return domain.Tag{}, nil
}
