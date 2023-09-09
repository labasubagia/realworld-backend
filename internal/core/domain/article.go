package domain

import (
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
}

func RandomArticle(author User) Article {
	return Article{
		AuthorID:    author.ID,
		Title:       util.RandomString(10),
		Description: util.RandomString(15),
		Slug:        util.RandomString(5),
		Body:        util.RandomString(20),
	}
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
