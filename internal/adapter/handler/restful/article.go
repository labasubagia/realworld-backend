package restful

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
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

type UpdateArticleParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

type UpdateArticleRequest struct {
	Article UpdateArticleParams `json:"article"`
}

type UpdateArticleResponse struct {
	Article Article `json:"article"`
}

func (server *Server) UpdateArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	var req UpdateArticleRequest
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.Article().Update(context.Background(), port.UpdateArticleParams{
		AuthArg: authArg,
		Slug:    slug,
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

	c.JSON(http.StatusOK, res)
}

func (server *Server) DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	err = server.service.Article().Delete(context.Background(), port.DeleteArticleParams{
		AuthArg: authArg,
		Slug:    slug,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

type Comment struct {
	ID        int64   `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Body      string  `json:"body"`
	Author    Profile `json:"author"`
}

type AddCommentRequest struct {
	Comment Comment `json:"comment"`
}

type AddCommentResponse struct {
	Comment Comment `json:"comment"`
}

func (server *Server) AddComment(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	var req AddCommentRequest
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.Article().AddComment(context.Background(), port.AddCommentParams{
		AuthArg: authArg,
		Slug:    slug,
		Comment: domain.Comment{
			Body: req.Comment.Body,
		},
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := AddCommentResponse{
		Comment: Comment{
			ID:        int64(result.Comment.ID),
			CreatedAt: result.Comment.CreatedAt.UTC().Format(formatTime),
			UpdatedAt: result.Comment.UpdatedAt.UTC().Format(formatTime),
			Body:      result.Comment.Body,
			Author: Profile{
				Username:  result.Comment.Author.Username,
				Bio:       result.Comment.Author.Bio,
				Image:     result.Comment.Author.Image,
				Following: result.Comment.Author.IsFollowed,
			},
		},
	}
	c.JSON(http.StatusOK, res)
}

type ListCommentResponse struct {
	Comments []Comment `json:"comments"`
}

func (server *Server) ListComments(c *gin.Context) {
	slug := c.Param("slug")
	authArg, _ := getAuthArg(c)

	result, err := server.service.Article().ListComments(context.Background(), port.ListCommentParams{
		AuthArg: authArg,
		Slug:    slug,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := ListCommentResponse{
		Comments: []Comment{},
	}
	for _, comment := range result.Comments {
		res.Comments = append(res.Comments, Comment{
			ID:        int64(comment.ID),
			CreatedAt: comment.CreatedAt.UTC().Format(formatTime),
			UpdatedAt: comment.UpdatedAt.UTC().Format(formatTime),
			Body:      comment.Body,
			Author: Profile{
				Username:  comment.Author.Username,
				Bio:       comment.Author.Bio,
				Image:     comment.Author.Image,
				Following: comment.Author.IsFollowed,
			},
		})
	}
	c.JSON(http.StatusOK, res)
}

func (server *Server) DeleteComment(c *gin.Context) {
	slug := c.Param("slug")
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		err = exception.Validation().AddError("comment_id", "should valid id")
		errorHandler(c, err)
		return
	}

	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	err = server.service.Article().DeleteComment(context.Background(), port.DeleteCommentParams{
		AuthArg:   authArg,
		Slug:      slug,
		CommentID: domain.ID(commentID),
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

type FavoriteArticleResponse struct {
	Article Article `json:"article"`
}

func (server *Server) AddFavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.Article().AddFavorite(context.Background(), port.AddFavoriteParams{
		AuthArg: authArg,
		Slug:    slug,
		UserID:  authArg.Payload.UserID,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	tags := []string{}
	if len(result.Article.TagNames) > 0 {
		tags = result.Article.TagNames
	}
	res := FavoriteArticleResponse{
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

	c.JSON(http.StatusOK, res)
}

func (server *Server) RemoveFavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.Article().RemoveFavorite(context.Background(), port.RemoveFavoriteParams{
		AuthArg: authArg,
		Slug:    slug,
		UserID:  authArg.Payload.UserID,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	tags := []string{}
	if len(result.Article.TagNames) > 0 {
		tags = result.Article.TagNames
	}
	res := FavoriteArticleResponse{
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

	c.JSON(http.StatusOK, res)
}
