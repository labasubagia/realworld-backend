package repository

import (
	"context"
	"errors"
	"strings"

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

func (r *articleRepo) CreateArticle(ctx context.Context, req port.CreateArticleParams) (domain.Article, error) {
	article := req.Article
	_, err := r.db.NewInsert().Model(&article).Exec(ctx)
	if err != nil {
		return domain.Article{}, err
	}
	return article, nil
}

func (r *articleRepo) AddTagsIfNotExists(ctx context.Context, arg port.AddTagsParams) ([]domain.Tag, error) {
	if len(arg.Tags) == 0 {
		return []domain.Tag{}, errors.New("tags cannot be empty")
	}

	existing := []domain.Tag{}
	err := r.db.NewSelect().Model(&existing).Where("name IN (?)", bun.In(arg.Tags)).Scan(ctx)
	if err != nil {
		return []domain.Tag{}, err
	}
	existMap := map[string]domain.Tag{}
	for _, tag := range existing {
		existMap[strings.ToLower(tag.Name)] = tag
	}

	newTags := []domain.Tag{}
	for _, tag := range arg.Tags {
		name := strings.ToLower(tag)
		if _, exist := existMap[name]; exist {
			continue
		}
		newTags = append(newTags, domain.Tag{Name: name})
	}
	if len(newTags) > 0 {
		_, err = r.db.NewInsert().Model(&newTags).Exec(ctx)
		if err != nil {
			return []domain.Tag{}, err
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
		return []domain.ArticleTag{}, err
	}

	return articleTags, nil
}
