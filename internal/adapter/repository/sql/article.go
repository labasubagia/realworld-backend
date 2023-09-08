package sql

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
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

func (r *articleRepo) CreateArticle(ctx context.Context, req port.CreateArticleParams) (domain.Article, error) {
	article := req.Article
	_, err := r.db.NewInsert().Model(&article).Exec(ctx)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	return article, nil
}

func (r *articleRepo) FilterTags(ctx context.Context, filter port.FilterTagParams) ([]domain.Tag, error) {
	result := []domain.Tag{}
	query := r.db.NewSelect().Model(&result)
	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.IDs))
	}
	if len(filter.Names) > 0 {
		query = query.Where("name IN (?)", bun.In(filter.Names))
	}
	err := query.Scan(ctx)
	if err != nil {
		return []domain.Tag{}, intoException(err)
	}
	return result, nil
}

func (r *articleRepo) AddTagsIfNotExists(ctx context.Context, arg port.AddTagsParams) ([]domain.Tag, error) {
	if len(arg.Tags) == 0 {
		return []domain.Tag{}, exception.New(exception.TypeValidation, "tags empty", nil)
	}

	existing, err := r.FilterTags(ctx, port.FilterTagParams{Names: arg.Tags})
	if err != nil {
		return []domain.Tag{}, intoException(err)
	}
	existMap := map[string]domain.Tag{}
	for _, tag := range existing {
		existMap[tag.Name] = tag
	}

	newTags := []domain.Tag{}
	for _, tag := range arg.Tags {
		if _, exist := existMap[tag]; exist {
			continue
		}
		newTags = append(newTags, domain.Tag{Name: tag})
	}
	if len(newTags) > 0 {
		_, err = r.db.NewInsert().Model(&newTags).Returning("*").Exec(ctx)
		if err != nil {
			return []domain.Tag{}, intoException(err)
		}
	}

	// return existing and new tags
	return append(existing, newTags...), nil
}

func (r *articleRepo) AssignTags(ctx context.Context, arg port.AssignTags) ([]domain.ArticleTag, error) {
	if len(arg.TagIDs) == 0 {
		return []domain.ArticleTag{}, exception.New(exception.TypeValidation, "tags empty", nil)
	}

	articleTags := make([]domain.ArticleTag, len(arg.TagIDs))
	for i, tagID := range arg.TagIDs {
		articleTags[i] = domain.ArticleTag{
			ArticleID: arg.ArticleID,
			TagID:     tagID,
		}
	}
	_, err := r.db.NewInsert().Model(&articleTags).Exec(ctx)
	if err != nil {
		return []domain.ArticleTag{}, intoException(err)
	}

	return articleTags, nil
}
