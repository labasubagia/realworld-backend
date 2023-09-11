package sql

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql/model"
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

func (r *articleRepo) CreateArticle(ctx context.Context, arg domain.Article) (domain.Article, error) {
	article := model.AsArticle(arg)
	_, err := r.db.NewInsert().Model(&article).Exec(ctx)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	return article.ToDomain(), nil
}

func (r *articleRepo) UpdateArticle(ctx context.Context, arg domain.Article) (domain.Article, error) {
	if arg.Title != "" {
		arg.SetTitle(arg.Title)
	}
	article := model.AsArticle(arg)

	// update
	_, err := r.db.NewUpdate().Model(&article).OmitZero().Where("id = ?", article.ID).Exec(ctx)
	if err != nil {
		return domain.Article{}, intoException(err)
	}

	// find updated
	updated, err := r.FindOneArticle(ctx, port.FilterArticlePayload{IDs: []domain.ID{article.ID}})
	if err != nil {
		return domain.Article{}, intoException(err)
	}

	return updated, err
}

func (r *articleRepo) DeleteArticle(ctx context.Context, arg domain.Article) error {
	article := model.AsArticle(arg)
	_, err := r.db.NewDelete().
		Model(&article).
		Where("id = ?", article.ID).
		Where("slug = ?", article.Slug).
		Exec(ctx)
	if err != nil {
		return intoException(err)
	}
	return nil
}

func (r *articleRepo) FilterArticle(ctx context.Context, filter port.FilterArticlePayload) ([]domain.Article, error) {
	articles := []model.Article{}
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
	result := []domain.Article{}
	for _, article := range articles {
		result = append(result, article.ToDomain())
	}
	return result, nil
}

func (r *articleRepo) FindOneArticle(ctx context.Context, filter port.FilterArticlePayload) (domain.Article, error) {
	articles, err := r.FilterArticle(ctx, filter)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	if len(articles) == 0 {
		return domain.Article{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}
	return articles[0], nil
}

func (r *articleRepo) FilterTags(ctx context.Context, filter port.FilterTagPayload) ([]domain.Tag, error) {
	tags := []model.Tag{}
	query := r.db.NewSelect().Model(&tags)
	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.IDs))
	}
	if len(filter.Names) > 0 {
		query = query.Where("name IN (?)", bun.In(filter.Names))
	}
	query = query.Order("name ASC")
	err := query.Scan(ctx)
	if err != nil {
		return []domain.Tag{}, intoException(err)
	}

	result := []domain.Tag{}
	for _, tag := range tags {
		result = append(result, tag.ToDomain())
	}
	return result, nil
}

func (r *articleRepo) AddTags(ctx context.Context, arg port.AddTagsPayload) ([]domain.Tag, error) {
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

	newTags := []model.Tag{}
	for _, tag := range arg.Tags {
		if _, exist := existMap[tag]; exist {
			continue
		}
		newTag := domain.NewTag(domain.Tag{Name: tag})
		newTags = append(newTags, model.AsTag(newTag))
	}
	if len(newTags) > 0 {
		_, err = r.db.NewInsert().Model(&newTags).Returning("*").Exec(ctx)
		if err != nil {
			return []domain.Tag{}, intoException(err)
		}
	}

	resultNewTags := []domain.Tag{}
	for _, tag := range newTags {
		resultNewTags = append(resultNewTags, tag.ToDomain())
	}

	// return existing and new tags
	return append(existing, resultNewTags...), nil
}

func (r *articleRepo) AssignArticleTags(ctx context.Context, arg port.AssignTagPayload) ([]domain.ArticleTag, error) {
	if len(arg.TagIDs) == 0 {
		return []domain.ArticleTag{}, exception.Validation().AddError("tags", "empty")
	}

	articleTags := make([]model.ArticleTag, len(arg.TagIDs))
	for i, tagID := range arg.TagIDs {
		articleTags[i] = model.ArticleTag{
			ArticleID: arg.ArticleID,
			TagID:     tagID,
		}
	}
	_, err := r.db.NewInsert().Model(&articleTags).Exec(ctx)
	if err != nil {
		return []domain.ArticleTag{}, intoException(err)
	}

	result := []domain.ArticleTag{}
	for _, tag := range articleTags {
		result = append(result, tag.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FilterArticleTags(ctx context.Context, arg port.FilterArticleTagPayload) ([]domain.ArticleTag, error) {
	articleTags := []model.ArticleTag{}
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

	result := []domain.ArticleTag{}
	for _, articleTag := range articleTags {
		result = append(result, articleTag.ToDomain())
	}
	return result, nil
}

func (r *articleRepo) AddFavorite(ctx context.Context, arg domain.ArticleFavorite) (domain.ArticleFavorite, error) {
	favorite := model.AsArticleFavorite(arg)
	_, err := r.db.NewInsert().Model(&favorite).Exec(ctx)
	if err != nil {
		return domain.ArticleFavorite{}, intoException(err)
	}
	return favorite.ToDomain(), nil
}

func (r *articleRepo) RemoveFavorite(ctx context.Context, arg domain.ArticleFavorite) (domain.ArticleFavorite, error) {
	favorite := model.AsArticleFavorite(arg)
	_, err := r.db.NewDelete().
		Model(&favorite).
		Where("article_id = ?", favorite.ArticleID).
		Where("user_id = ?", favorite.UserID).
		Exec(ctx)
	if err != nil {
		return domain.ArticleFavorite{}, intoException(err)
	}
	return favorite.ToDomain(), nil
}

func (r *articleRepo) FilterFavorite(ctx context.Context, arg port.FilterFavoritePayload) ([]domain.ArticleFavorite, error) {
	articleFavorites := []model.ArticleFavorite{}
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

	result := []domain.ArticleFavorite{}
	for _, articleFavorite := range articleFavorites {
		result = append(result, articleFavorite.ToDomain())
	}
	return result, nil
}

func (r *articleRepo) FilterFavoriteCount(ctx context.Context, filter port.FilterFavoritePayload) ([]domain.ArticleFavoriteCount, error) {
	counts := []model.ArticleFavoriteCount{}
	err := r.db.NewSelect().
		Model(&counts).
		Column("article_id").
		ColumnExpr("count(article_id) as favorite_count").
		Group("article_id").
		Scan(ctx)
	if err != nil {
		return []domain.ArticleFavoriteCount{}, intoException(err)
	}
	result := []domain.ArticleFavoriteCount{}
	for _, count := range counts {
		result = append(result, count.ToDomain())
	}
	return result, nil
}

func (r *articleRepo) AddComment(ctx context.Context, arg domain.Comment) (domain.Comment, error) {
	comment := model.AsComment(arg)
	_, err := r.db.NewInsert().Model(&comment).Exec(ctx)
	if err != nil {
		return domain.Comment{}, intoException(err)
	}
	return comment.ToDomain(), nil
}

func (r *articleRepo) FilterComment(ctx context.Context, arg port.FilterCommentPayload) ([]domain.Comment, error) {
	comments := []model.Comment{}
	query := r.db.NewSelect().Model(&comments)
	if len(arg.ArticleIDs) > 0 {
		query = query.Where("article_id IN (?)", bun.In(arg.ArticleIDs))
	}
	if len(arg.AuthorIDs) > 0 {
		query = query.Where("author_id IN (?)", bun.In(arg.AuthorIDs))
	}
	err := query.Scan(ctx)
	if err != nil {
		return []domain.Comment{}, intoException(err)
	}
	result := []domain.Comment{}
	for _, comment := range comments {
		result = append(result, comment.ToDomain())
	}
	return result, nil
}

func (r *articleRepo) DeleteComment(ctx context.Context, arg domain.Comment) error {
	comment := model.AsComment(arg)
	_, err := r.db.NewDelete().
		Model(&comment).
		Where("id = ?", comment.ID).
		Where("author_id = ?", comment.AuthorID).
		Where("article_id = ?", comment.ArticleID).
		Exec(ctx)
	if err != nil {
		return intoException(err)
	}
	return nil
}
