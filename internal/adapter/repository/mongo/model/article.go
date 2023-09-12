package model

import (
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type Article struct {
	ID          domain.ID `bson:"id"`
	AuthorID    domain.ID `bson:"author_id"`
	Title       string    `bson:"title"`
	Slug        string    `bson:"slug"`
	Description string    `bson:"description"`
	Body        string    `bson:"body"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func (data Article) ToDomain() domain.Article {
	return domain.Article{
		ID:          data.ID,
		AuthorID:    data.AuthorID,
		Title:       data.Title,
		Slug:        data.Slug,
		Description: data.Description,
		Body:        data.Body,
		CreatedAt:   data.CreatedAt.UTC(),
		UpdatedAt:   data.UpdatedAt.UTC(),
	}
}

func AsArticle(arg domain.Article) Article {
	return Article{
		ID:          arg.ID,
		AuthorID:    arg.AuthorID,
		Title:       arg.Title,
		Slug:        arg.Slug,
		Description: arg.Description,
		Body:        arg.Body,
		CreatedAt:   arg.CreatedAt.UTC(),
		UpdatedAt:   arg.UpdatedAt.UTC(),
	}
}

type Tag struct {
	ID   domain.ID `bson:"id"`
	Name string    `bson:"name"`
}

func (data Tag) ToDomain() domain.Tag {
	return domain.Tag{
		ID:   data.ID,
		Name: data.Name,
	}
}

func AsTag(arg domain.Tag) Tag {
	return Tag{
		ID:   arg.ID,
		Name: arg.Name,
	}
}

type ArticleTag struct {
	ArticleID domain.ID `bson:"article_id"`
	TagID     domain.ID `bson:"tag_id"`
}

func (data ArticleTag) ToDomain() domain.ArticleTag {
	return domain.ArticleTag{
		ArticleID: data.ArticleID,
		TagID:     data.TagID,
	}
}

func AsArticleTag(arg domain.ArticleTag) ArticleTag {
	return ArticleTag{
		ArticleID: arg.ArticleID,
		TagID:     arg.TagID,
	}
}

type Comment struct {
	ID        domain.ID `bson:"id"`
	ArticleID domain.ID `bson:"article_id"`
	AuthorID  domain.ID `bson:"author_id"`
	Body      string    `bson:"body"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (data Comment) ToDomain() domain.Comment {
	return domain.Comment{
		ID:        data.ID,
		ArticleID: data.ArticleID,
		AuthorID:  data.AuthorID,
		Body:      data.Body,
		CreatedAt: data.CreatedAt.UTC(),
		UpdatedAt: data.UpdatedAt.UTC(),
	}
}

func AsComment(arg domain.Comment) Comment {
	return Comment{
		ID:        arg.ID,
		ArticleID: arg.ArticleID,
		AuthorID:  arg.AuthorID,
		Body:      arg.Body,
		CreatedAt: arg.CreatedAt.UTC(),
		UpdatedAt: arg.UpdatedAt.UTC(),
	}
}

type ArticleFavorite struct {
	ArticleID domain.ID `bson:"article_id"`
	UserID    domain.ID `bson:"user_id"`
}

func (data ArticleFavorite) ToDomain() domain.ArticleFavorite {
	return domain.ArticleFavorite{
		ArticleID: data.ArticleID,
		UserID:    data.UserID,
	}
}

func AsArticleFavorite(arg domain.ArticleFavorite) ArticleFavorite {
	return ArticleFavorite{
		ArticleID: arg.ArticleID,
		UserID:    arg.UserID,
	}
}

type ArticleFavoriteCount struct {
	ArticleID domain.ID `bson:"article_id"`
	Count     int       `bson:"favorite_count"`
}

func (data ArticleFavoriteCount) ToDomain() domain.ArticleFavoriteCount {
	return domain.ArticleFavoriteCount{
		ArticleID: data.ArticleID,
		Count:     data.Count,
	}
}

func AsArticleFavoriteCount(arg domain.ArticleFavoriteCount) ArticleFavoriteCount {
	return ArticleFavoriteCount{
		ArticleID: arg.ArticleID,
		Count:     arg.Count,
	}
}
