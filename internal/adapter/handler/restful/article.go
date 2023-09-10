package restful

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
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
		tags := []string{}
		if len(article.TagNames) > 0 {
			tags = article.TagNames
		}
		res.Articles = append(res.Articles, Article{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tags,
			CreatedAt:      article.CreatedAt.UTC().Format(formatTime),
			UpdatedAt:      article.UpdatedAt.UTC().Format(formatTime),
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

func (server *Server) FeedArticle(c *gin.Context) {
	offset, limit := getPagination(c)
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	arg := port.ListArticleParams{
		AuthArg: authArg,
		Offset:  offset,
		Limit:   limit,
	}
	result, err := server.service.Article().Feed(context.Background(), arg)
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := LisArticleResult{
		Articles: []Article{},
		Count:    result.Count,
	}
	for _, article := range result.Articles {
		tags := []string{}
		if len(article.TagNames) > 0 {
			tags = article.TagNames
		}
		res.Articles = append(res.Articles, Article{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tags,
			CreatedAt:      article.CreatedAt.UTC().Format(formatTime),
			UpdatedAt:      article.UpdatedAt.UTC().Format(formatTime),
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

type GetArticleResult struct {
	Article Article `json:"article"`
}

func (server *Server) GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, _ := getAuthArg(c)

	result, err := server.service.Article().Get(context.Background(), port.GetArticleParams{
		AuthArg: authArg,
		Slug:    slug,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	tags := []string{}
	if len(result.Article.TagNames) > 0 {
		tags = result.Article.TagNames
	}

	res := GetArticleResult{
		Article{
			Slug:           result.Article.Slug,
			Title:          result.Article.Title,
			Description:    result.Article.Description,
			Body:           result.Article.Body,
			TagList:        tags,
			CreatedAt:      result.Article.CreatedAt.UTC().Format(formatTime),
			UpdatedAt:      result.Article.UpdatedAt.UTC().Format(formatTime),
			Favorited:      result.Article.IsFavorite,
			FavoritesCount: result.Article.FavoriteCount,
			Author: Profile{
				Username:  result.Article.Author.Username,
				Bio:       result.Article.Author.Bio,
				Image:     result.Article.Author.Image,
				Following: result.Article.Author.IsFollowed,
			},
		},
	}
	c.JSON(http.StatusOK, res)
}

type CreateArticleParams struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type CreateArticleRequest struct {
	Article CreateArticleParams `json:"article"`
}

type CreateArticleResponse struct {
	Article Article `json:"article"`
}

func (server *Server) CreateArticle(c *gin.Context) {
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	var req CreateArticleRequest
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.Article().Create(context.Background(), port.CreateArticleTxParams{
		AuthArg: authArg,
		Tags:    req.Article.TagList,
		Article: domain.Article{
			Title:       req.Article.Title,
			Description: req.Article.Description,
			Body:        req.Article.Body,
		},
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	tags := []string{}
	if len(result.Article.TagNames) > 0 {
		tags = result.Article.TagNames
	}
	res := CreateArticleResponse{
		Article: Article{
			Slug:           result.Article.Slug,
			Title:          result.Article.Title,
			Description:    result.Article.Description,
			Body:           result.Article.Body,
			TagList:        tags,
			CreatedAt:      result.Article.CreatedAt.UTC().Format(formatTime),
			UpdatedAt:      result.Article.UpdatedAt.UTC().Format(formatTime),
			Favorited:      result.Article.IsFavorite,
			FavoritesCount: result.Article.FavoriteCount,
			Author: Profile{
				Username:  result.Article.Author.Username,
				Bio:       result.Article.Author.Bio,
				Image:     result.Article.Author.Image,
				Following: result.Article.Author.IsFollowed,
			},
		},
	}

	c.JSON(http.StatusCreated, res)
}
