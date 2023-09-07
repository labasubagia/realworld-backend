package repository

import (
	"context"
	"errors"

	"github.com/labasubagia/realworld-backend/domain"
	"github.com/labasubagia/realworld-backend/port"
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
		return domain.Article{}, AsServiceError(err)
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
		return []domain.Tag{}, AsServiceError(err)
	}
	return result, nil
}

func (r *articleRepo) AddTagsIfNotExists(ctx context.Context, arg port.AddTagsParams) ([]domain.Tag, error) {
	if len(arg.Tags) == 0 {
		return []domain.Tag{}, errors.New("tags cannot be empty")
	}

	existing, err := r.FilterTags(ctx, port.FilterTagParams{Names: arg.Tags})
	if err != nil {
		return []domain.Tag{}, AsServiceError(err)
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
			return []domain.Tag{}, AsServiceError(err)
		}
	}

	// return existing and new tags
	return append(existing, newTags...), nil
}

func (r *articleRepo) AssignTags(ctx context.Context, arg port.AssignTags) ([]domain.ArticleTag, error) {
	if len(arg.TagIDs) == 0 {
		return []domain.ArticleTag{}, errors.New("tag ids cannot be empty")
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
		return []domain.ArticleTag{}, AsServiceError(err)
	}

	return articleTags, nil
}
