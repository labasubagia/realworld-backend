package mongo

import (
	"context"
	"time"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/mongo/model"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type articleRepo struct {
	db DB
}

func NewArticleRepository(db DB) port.ArticleRepository {
	return &articleRepo{
		db: db,
	}
}

func (r *articleRepo) AddComment(ctx context.Context, arg domain.Comment) (domain.Comment, error) {
	comment := model.AsComment(arg)
	_, err := r.db.Collection(CollectionComment).InsertOne(ctx, comment)
	if err != nil {
		return domain.Comment{}, intoException(err)
	}
	return comment.ToDomain(), nil
}

func (r *articleRepo) AddFavorite(ctx context.Context, arg domain.ArticleFavorite) (domain.ArticleFavorite, error) {
	favorite := model.AsArticleFavorite(arg)
	_, err := r.db.Collection(CollectionArticleFavorite).InsertOne(ctx, favorite)
	if err != nil {
		return domain.ArticleFavorite{}, intoException(err)
	}
	return favorite.ToDomain(), nil
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

	newTags := []any{}
	for _, tag := range arg.Tags {
		if _, exist := existMap[tag]; exist {
			continue
		}
		newTag := domain.NewTag(domain.Tag{Name: tag})
		newTags = append(newTags, model.AsTag(newTag))
	}
	if len(newTags) > 0 {
		_, err = r.db.Collection(CollectionTag).InsertMany(ctx, newTags)
		if err != nil {
			return []domain.Tag{}, intoException(err)
		}
	}

	resultNewTags := []domain.Tag{}
	for _, tag := range newTags {
		tag, ok := tag.(model.Tag)
		if !ok {
			continue
		}
		resultNewTags = append(resultNewTags, tag.ToDomain())
	}

	// return existing and new tags
	return append(existing, resultNewTags...), nil

}

func (r *articleRepo) AssignArticleTags(ctx context.Context, arg port.AssignTagPayload) ([]domain.ArticleTag, error) {
	if len(arg.TagIDs) == 0 {
		return []domain.ArticleTag{}, exception.Validation().AddError("tags", "empty")
	}

	articleTags := make([]any, len(arg.TagIDs))
	for i, tagID := range arg.TagIDs {
		articleTags[i] = model.ArticleTag{
			ArticleID: arg.ArticleID,
			TagID:     tagID,
		}
	}
	_, err := r.db.Collection(CollectionArticleTag).InsertMany(ctx, articleTags)
	if err != nil {
		return []domain.ArticleTag{}, intoException(err)
	}

	result := []domain.ArticleTag{}
	for _, tag := range articleTags {
		tag, ok := tag.(model.ArticleTag)
		if !ok {
			continue
		}
		result = append(result, tag.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) CreateArticle(ctx context.Context, arg domain.Article) (domain.Article, error) {
	article := model.AsArticle(arg)
	_, err := r.db.Collection(CollectionArticle).InsertOne(ctx, article)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	return article.ToDomain(), nil
}

func (r *articleRepo) DeleteArticle(ctx context.Context, arg domain.Article) error {
	_, err := r.db.Collection(CollectionArticle).DeleteOne(ctx, bson.M{
		"id":   arg.ID,
		"slug": arg.Slug,
	})
	if err != nil {
		return intoException(err)
	}
	return nil
}

func (r *articleRepo) DeleteComment(ctx context.Context, arg domain.Comment) error {
	_, err := r.db.Collection(CollectionComment).DeleteOne(ctx, bson.M{
		"id":         arg.ID,
		"author_id":  arg.AuthorID,
		"article_id": arg.ArticleID,
	})
	if err != nil {
		return intoException(err)
	}
	return nil
}

func (r *articleRepo) FilterArticle(ctx context.Context, arg port.FilterArticlePayload) ([]domain.Article, error) {

	query := []bson.M{}
	if len(arg.IDs) > 0 {
		query = append(query, bson.M{"id": bson.M{"$in": arg.IDs}})
	}
	if len(arg.AuthorIDs) > 0 {
		query = append(query, bson.M{"author_id": bson.M{"$in": arg.AuthorIDs}})
	}
	if len(arg.Slugs) > 0 {
		query = append(query, bson.M{"slug": bson.M{"$in": arg.Slugs}})
	}
	if len(arg.Slugs) > 0 {
		query = append(query, bson.M{"slug": bson.M{"$in": arg.Slugs}})
	}
	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}

	limit := int64(arg.Limit)
	offset := int64(arg.Offset)
	option := options.FindOptions{Limit: &limit, Skip: &offset, Sort: bson.M{"id": -1}}

	cursor, err := r.db.Collection(CollectionArticle).Find(ctx, filter, &option)
	if err != nil {
		return []domain.Article{}, intoException(err)
	}

	result := []domain.Article{}
	for cursor.Next(ctx) {
		data := model.Article{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.Article{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FilterArticleTags(ctx context.Context, arg port.FilterArticleTagPayload) ([]domain.ArticleTag, error) {

	query := []bson.M{}
	if len(arg.ArticleIDs) > 0 {
		query = append(query, bson.M{"article_id": bson.M{"$in": arg.ArticleIDs}})
	}
	if len(arg.TagIDs) > 0 {
		query = append(query, bson.M{"tag_id": bson.M{"$in": arg.TagIDs}})
	}
	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}

	cursor, err := r.db.Collection(CollectionArticleTag).Find(ctx, filter)
	if err != nil {
		return []domain.ArticleTag{}, intoException(err)
	}

	result := []domain.ArticleTag{}
	for cursor.Next(ctx) {
		data := model.ArticleTag{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.ArticleTag{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FilterComment(ctx context.Context, arg port.FilterCommentPayload) ([]domain.Comment, error) {
	query := []bson.M{}
	if len(arg.ArticleIDs) > 0 {
		query = append(query, bson.M{"article_id": bson.M{"$in": arg.ArticleIDs}})
	}
	if len(arg.AuthorIDs) > 0 {
		query = append(query, bson.M{"author_id": bson.M{"$in": arg.AuthorIDs}})
	}
	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}

	cursor, err := r.db.Collection(CollectionComment).Find(ctx, filter)
	if err != nil {
		return []domain.Comment{}, intoException(err)
	}

	result := []domain.Comment{}
	for cursor.Next(ctx) {
		data := model.Comment{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.Comment{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FilterFavorite(ctx context.Context, arg port.FilterFavoritePayload) ([]domain.ArticleFavorite, error) {
	query := []bson.M{}
	if len(arg.ArticleIDs) > 0 {
		query = append(query, bson.M{"article_id": bson.M{"$in": arg.ArticleIDs}})
	}
	if len(arg.UserIDs) > 0 {
		query = append(query, bson.M{"user_id": bson.M{"$in": arg.UserIDs}})
	}
	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}

	cursor, err := r.db.Collection(CollectionArticleFavorite).Find(ctx, filter)
	if err != nil {
		return []domain.ArticleFavorite{}, intoException(err)
	}

	result := []domain.ArticleFavorite{}
	for cursor.Next(ctx) {
		data := model.ArticleFavorite{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.ArticleFavorite{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FilterFavoriteCount(ctx context.Context, arg port.FilterFavoritePayload) ([]domain.ArticleFavoriteCount, error) {

	groupStage := bson.A{}

	// match query stage
	matchQueryItems := []bson.M{}
	if len(arg.ArticleIDs) > 0 {
		matchQueryItems = append(matchQueryItems, bson.M{"article_id": bson.M{"$in": arg.ArticleIDs}})
	}
	if len(arg.UserIDs) > 0 {
		matchQueryItems = append(matchQueryItems, bson.M{"user_id": bson.M{"$in": arg.UserIDs}})
	}
	if len(matchQueryItems) > 0 {
		groupStage = append(groupStage, bson.M{
			"$match": bson.M{"$and": matchQueryItems},
		})
	}

	groupStage = append(
		groupStage,
		// grouping stage
		bson.M{
			"$group": bson.M{
				"_id": "$article_id",
				"favorite_count": bson.M{
					"$count": bson.M{},
				},
			},
		},
		// custom field stage
		bson.M{
			"$addFields": bson.M{
				"article_id": "$_id",
				"_id":        "$$REMOVE",
			},
		},
	)

	cursor, err := r.db.Collection(CollectionArticleFavorite).Aggregate(ctx, groupStage)
	if err != nil {
		return []domain.ArticleFavoriteCount{}, intoException(err)
	}

	result := []domain.ArticleFavoriteCount{}
	for cursor.Next(ctx) {
		data := model.ArticleFavoriteCount{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.ArticleFavoriteCount{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FilterTags(ctx context.Context, arg port.FilterTagPayload) ([]domain.Tag, error) {

	query := []bson.M{}
	if len(arg.IDs) > 0 {
		query = append(query, bson.M{"id": bson.M{"$in": arg.IDs}})
	}
	if len(arg.Names) > 0 {
		query = append(query, bson.M{"name": bson.M{"$in": arg.Names}})
	}
	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}
	option := options.FindOptions{Sort: bson.M{"name": 1}}

	cursor, err := r.db.Collection(CollectionTag).Find(ctx, filter, &option)
	if err != nil {
		return []domain.Tag{}, intoException(err)
	}

	result := []domain.Tag{}
	for cursor.Next(ctx) {
		tag := model.Tag{}
		if err := cursor.Decode(&tag); err != nil {
			return []domain.Tag{}, intoException(err)
		}
		result = append(result, tag.ToDomain())
	}

	return result, nil
}

func (r *articleRepo) FindOneArticle(ctx context.Context, arg port.FilterArticlePayload) (domain.Article, error) {
	articles, err := r.FilterArticle(ctx, arg)
	if err != nil {
		return domain.Article{}, intoException(err)
	}
	if len(articles) == 0 {
		return domain.Article{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}
	return articles[0], nil
}

func (r *articleRepo) RemoveFavorite(ctx context.Context, arg domain.ArticleFavorite) (domain.ArticleFavorite, error) {
	favorite := model.AsArticleFavorite(arg)
	_, err := r.db.Collection(CollectionArticleFavorite).DeleteOne(ctx, bson.M{
		"user_id":    arg.UserID,
		"article_id": arg.ArticleID,
	})
	if err != nil {
		return domain.ArticleFavorite{}, intoException(err)
	}
	return favorite.ToDomain(), nil
}

func (r *articleRepo) UpdateArticle(ctx context.Context, arg domain.Article) (domain.Article, error) {
	if arg.Title != "" {
		arg.SetTitle(arg.Title)
	}
	article := model.AsArticle(arg)

	filter := bson.M{"_id": arg.ID}

	fields := bson.M{}
	if arg.Title != "" {
		fields["title"] = arg.Title
	}
	if arg.Slug != "" {
		fields["slug"] = arg.Slug
	}
	if arg.Body != "" {
		fields["body"] = arg.Body
	}
	if arg.Description != "" {
		fields["description"] = arg.Description
	}
	if len(fields) > 0 {
		fields["updated_at"] = time.Now()
	}

	_, err := r.db.Collection(CollectionArticle).UpdateOne(ctx, filter, bson.M{"$set": fields})
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
