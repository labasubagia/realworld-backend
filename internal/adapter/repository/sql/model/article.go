package model

import (
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/uptrace/bun"
)

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:a"`
	ID            domain.ID `bun:"id,pk"`
	AuthorID      domain.ID `bun:"author_id,notnull"`
	Title         string    `bun:"title,notnull"`
	Slug          string    `bun:"slug,notnull"`
	Description   string    `bun:"description,notnull"`
	Body          string    `bun:"body,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (data Article) ToDomain() domain.Article {
	return domain.Article{
		ID:          data.ID,
		AuthorID:    data.AuthorID,
		Title:       data.Title,
		Slug:        data.Slug,
		Description: data.Description,
		Body:        data.Body,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
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
		CreatedAt:   arg.CreatedAt,
		UpdatedAt:   arg.UpdatedAt,
	}
}

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`
	ID            domain.ID `bun:"id,pk"`
	Name          string    `bun:"name,notnull"`
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
	bun.BaseModel `bun:"table:article_tags,alias:at"`
	ArticleID     domain.ID `bun:"article_id,notnull"`
	TagID         domain.ID `bun:"tag_id,notnull"`
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
	bun.BaseModel `bun:"table:comments,alias:c"`
	ID            domain.ID `bun:"id,pk"`
	ArticleID     domain.ID `bun:"article_id,notnull"`
	AuthorID      domain.ID `bun:"author_id,notnull"`
	Body          string    `bun:"body,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (data Comment) ToDomain() domain.Comment {
	return domain.Comment{
		ID:        data.ID,
		ArticleID: data.ArticleID,
		AuthorID:  data.AuthorID,
		Body:      data.Body,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func AsComment(arg domain.Comment) Comment {
	return Comment{
		ID:        arg.ID,
		ArticleID: arg.ArticleID,
		AuthorID:  arg.AuthorID,
		Body:      arg.Body,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
	}
}

type ArticleFavorite struct {
	bun.BaseModel `bun:"table:article_favorites,alias:af"`
	ArticleID     domain.ID `bun:"article_id,notnull"`
	UserID        domain.ID `bun:"user_id,notnull"`
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
	bun.BaseModel `bun:"table:article_favorites,alias:af"`
	ArticleID     domain.ID `bun:"article_id"`
	Count         int       `bun:"favorite_count"`
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
