package restful

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

type Article struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	Author         Profile  `json:"author"`
}

type LisArticleResult struct {
	Articles []Article `json:"articles"`
	Count    int       `json:"articlesCount"`
}

func (server *Server) ListArticle(c *gin.Context) {
	Tag := c.Query("tag")
	Author := c.Query("author")
	FavoritedBy := c.Query("favorited")

	offset, limit := getPagination(c)
	authArg, _ := getAuthArg(c)

	arg := port.ListArticleParams{
		Tags:           []string{},
		AuthorNames:    []string{},
		FavoritedNames: []string{},
		AuthArg:        authArg,
		Offset:         offset,
		Limit:          limit,
	}
	if Tag != "" {
		arg.Tags = append(arg.Tags, Tag)
	}
	if Author != "" {
		arg.AuthorNames = append(arg.AuthorNames, Author)
	}
	if FavoritedBy != "" {
		arg.FavoritedNames = append(arg.FavoritedNames, FavoritedBy)
	}

	result, err := server.service.Article().List(context.Background(), arg)
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := LisArticleResult{
		Articles: []Article{},
		Count:    result.Count,
	}
	for _, article := range result.Articles {
		res.Articles = append(res.Articles, Article{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagNames,
			CreatedAt:      article.CreatedAt.Format(formatTime),
			UpdatedAt:      article.UpdatedAt.Format(formatTime),
			Favorited:      article.IsFavorite,
			FavoritesCount: article.FavoriteCount,
			Author: Profile{
				Username:  article.Author.Username,
				Bio:       article.Author.Bio,
				Image:     article.Author.Image,
				Following: article.Author.IsFollowed,
			},
		})
	}

	c.JSON(http.StatusOK, res)
}
