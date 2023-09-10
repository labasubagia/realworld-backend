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

func (r *articleRepo) CreateArticle(ctx context.Context, arg port.CreateArticlePayload) (domain.Article, error) {
	article := arg.Article
	_, err := r.db.NewInsert().Model(&article).Exec(ctx)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	return article, nil
}

func (r *articleRepo) FilterArticle(ctx context.Context, filter port.FilterArticlePayload) ([]domain.Article, error) {
	articles := []domain.Article{}
	query := r.db.NewSelect().Model(&articles)
	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.IDs))
	}
	if len(filter.Slugs) > 0 {
		query = query.Where("slug IN (?)", bun.In(filter.Slugs))
	}
	if len(filter.AuthorIDs) > 0 {
		query = query.Where("author_id IN (?)", bun.In(filter.AuthorIDs))
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	query = query.Offset(filter.Offset)
	query = query.Order("created_at DESC")
	err := query.Scan(ctx)
	if err != nil {
		return []domain.Article{}, nil
	}
	return articles, nil
}

func (r *articleRepo) FindOneArticle(ctx context.Context, filter port.FilterArticlePayload) (domain.Article, error) {
	articles, err := r.FilterArticle(ctx, filter)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	if len(articles) < 1 {
		return domain.Article{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}
	return articles[0], nil
}

func (r *articleRepo) FilterTags(ctx context.Context, filter port.FilterTagPayload) ([]domain.Tag, error) {
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

func (r *articleRepo) AddTagsIfNotExists(ctx context.Context, arg port.AddTagsPayload) ([]domain.Tag, error) {
	if len(arg.Tags) == 0 {
		return []domain.Tag{}, exception.Validation().AddError("tags", "empty")
	}

	existing, err := r.FilterTags(ctx, port.FilterTagPayload{Names: arg.Tags})
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

func (r *articleRepo) AssignArticleTags(ctx context.Context, arg port.AssignTagPayload) ([]domain.ArticleTag, error) {
	if len(arg.TagIDs) == 0 {
		return []domain.ArticleTag{}, exception.Validation().AddError("tags", "empty")
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

func (r *articleRepo) FilterArticleTags(ctx context.Context, arg port.FilterArticleTagPayload) ([]domain.ArticleTag, error) {
	articleTags := []domain.ArticleTag{}
	query := r.db.NewSelect().Model(&articleTags)
	if len(arg.TagIDs) > 0 {
		query = query.Where("tag_id IN (?)", bun.In(arg.TagIDs))
	}
	if len(arg.ArticleIDs) > 0 {
		query = query.Where("article_id IN (?)", bun.In(arg.ArticleIDs))
	}
	err := query.Scan(ctx)
	if err != nil {
		return []domain.ArticleTag{}, intoException(err)
	}
	return articleTags, nil
}

func (r *articleRepo) AddFavorite(ctx context.Context, arg domain.ArticleFavorite) (domain.ArticleFavorite, error) {
	favorite := arg
	_, err := r.db.NewInsert().Model(&favorite).Exec(ctx)
	if err != nil {
		return domain.ArticleFavorite{}, intoException(err)
	}
	return favorite, nil
}

func (r *articleRepo) FilterFavorite(ctx context.Context, arg port.FilterFavoritePayload) ([]domain.ArticleFavorite, error) {
	articleFavorites := []domain.ArticleFavorite{}
	query := r.db.NewSelect().Model(&articleFavorites)
	if len(arg.UserIDs) > 0 {
		query = query.Where("user_id IN (?)", bun.In(arg.UserIDs))
	}
	if len(arg.ArticleIDs) > 0 {
		query = query.Where("article_id IN (?)", bun.In(arg.ArticleIDs))
	}
	err := query.Scan(ctx)
	if err != nil {
		return []domain.ArticleFavorite{}, intoException(err)
	}
	return articleFavorites, nil
}

func (r *articleRepo) FilterFavoriteCount(ctx context.Context, filter port.FilterFavoritePayload) ([]domain.ArticleFavoriteCount, error) {
	counts := []domain.ArticleFavoriteCount{}
	err := r.db.NewSelect().
		Model(&counts).
		Column("article_id").
		ColumnExpr("count(article_id) as favorite_count").
		Group("article_id").
		Scan(ctx)
	if err != nil {
		return []domain.ArticleFavoriteCount{}, intoException(err)
	}
	return counts, nil
}
