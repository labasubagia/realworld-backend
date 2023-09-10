package domain

import (
	"strings"
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/uptrace/bun"
)

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:a"`
	ID            ID        `bun:"id,pk,autoincrement"`
	AuthorID      ID        `bun:"author_id,notnull"`
	Title         string    `bun:"title,notnull"`
	Slug          string    `bun:"slug,notnull"`
	Description   string    `bun:"description,notnull"`
	Body          string    `bun:"body,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	TagNames      []string `bun:"-"`
	Author        User     `bun:"-"`
	IsFavorite    bool     `bun:"-"`
	FavoriteCount int      `bun:"-"`
}

func (article *Article) SetTitle(value string) {
	article.Title = value
	article.Slug = strings.ToLower(strings.ReplaceAll(value, " ", "-"))
}

func NewArticle(arg Article) Article {
	now := time.Now()
	article := Article{
		AuthorID:    arg.AuthorID,
		Title:       arg.Title,
		Description: arg.Description,
		Body:        arg.Body,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	article.SetTitle(arg.Title)
	return article
}

func RandomArticle(author User) Article {
	article := Article{
		AuthorID:    author.ID,
		Description: util.RandomString(15),
		Body:        util.RandomString(20),
	}
	article.SetTitle(util.RandomString(10))
	return article
}

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`
	ID            ID     `bun:"id,pk,autoincrement"`
	Name          string `bun:"name,notnull"`
}

type ArticleTag struct {
	bun.BaseModel `bun:"table:article_tags,alias:at"`
	ArticleID     ID `bun:"article_id,notnull"`
	TagID         ID `bun:"tag_id,notnull"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comments,alias:c"`
	ID            ID        `bun:"id,pk,autoincrement"`
	ArticleID     ID        `bun:"article_id,notnull"`
	AuthorID      ID        `bun:"author_id,notnull"`
	Body          string    `bun:"body,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	Author User `bun:"-"`
}

type ArticleFavorite struct {
	bun.BaseModel `bun:"table:article_favorites,alias:af"`
	ArticleID     ID `bun:"article_id,notnull"`
	UserID        ID `bun:"user_id,notnull"`
}

type ArticleFavoriteCount struct {
	bun.BaseModel `bun:"table:article_favorites,alias:af"`
	ArticleID     ID  `bun:"article_id"`
	Count         int `bun:"favorite_count"`
}
