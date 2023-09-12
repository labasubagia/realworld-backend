package domain

import (
	"strings"
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/util"
)

type Article struct {
	ID            ID
	AuthorID      ID
	Title         string
	Slug          string
	Description   string
	Body          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TagNames      []string
	Author        User
	IsFavorite    bool
	FavoriteCount int
}

func (article *Article) SetTitle(value string) {
	article.Title = value
	article.Slug = strings.ToLower(strings.ReplaceAll(value, " ", "-"))
}

func NewArticle(arg Article) Article {
	now := time.Now()
	article := Article{
		ID:          NewID(),
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
	ID   ID
	Name string
}

func NewTag(arg Tag) Tag {
	return Tag{
		ID:   NewID(),
		Name: arg.Name,
	}
}

type ArticleTag struct {
	ArticleID ID
	TagID     ID
}

type Comment struct {
	ID        ID
	ArticleID ID
	AuthorID  ID
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    User
}

func NewComment(arg Comment) Comment {
	now := time.Now()
	return Comment{
		ID:        NewID(),
		ArticleID: arg.ArticleID,
		AuthorID:  arg.AuthorID,
		Body:      arg.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type ArticleFavorite struct {
	ArticleID ID
	UserID    ID
}

type ArticleFavoriteCount struct {
	ArticleID ID
	Count     int
}
