package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

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

	articles, err := server.service.Article().List(c, arg)
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := ArticlesResponse{
		Articles: []Article{},
		Count:    len(articles),
	}
	for _, article := range articles {
		res.Articles = append(res.Articles, serializeArticle(article))
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
	articles, err := server.service.Article().Feed(c, arg)
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := ArticlesResponse{
		Articles: []Article{},
		Count:    len(articles),
	}
	for _, article := range articles {
		res.Articles = append(res.Articles, serializeArticle(article))
	}

	c.JSON(http.StatusOK, res)
}

func (server *Server) GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, _ := getAuthArg(c)

	article, err := server.service.Article().Get(c, port.GetArticleParams{
		AuthArg: authArg,
		Slug:    slug,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := ArticleResponse{serializeArticle(article)}
	c.JSON(http.StatusOK, res)
}

type CreateArticle struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type CreateArticleRequest struct {
	Article CreateArticle `json:"article"`
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

	article, err := server.service.Article().Create(c, port.CreateArticleTxParams{
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

	res := ArticleResponse{serializeArticle(article)}
	c.JSON(http.StatusCreated, res)
}

type UpdateArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

type UpdateArticleRequest struct {
	Article UpdateArticle `json:"article"`
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

	article, err := server.service.Article().Update(c, port.UpdateArticleParams{
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

	res := ArticleResponse{serializeArticle(article)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	err = server.service.Article().Delete(c, port.DeleteArticleParams{
		AuthArg: authArg,
		Slug:    slug,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

type AddCommentRequest struct {
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

	result, err := server.service.Article().AddComment(c, port.AddCommentParams{
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

	res := CommentResponse{serializeComment(result)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) ListComments(c *gin.Context) {
	slug := c.Param("slug")
	authArg, _ := getAuthArg(c)

	comments, err := server.service.Article().ListComments(c, port.ListCommentParams{
		AuthArg: authArg,
		Slug:    slug,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := CommentsResponse{
		Comments: []Comment{},
	}
	for _, comment := range comments {
		res.Comments = append(res.Comments, serializeComment(comment))
	}
	c.JSON(http.StatusOK, res)
}

func (server *Server) DeleteComment(c *gin.Context) {
	slug := c.Param("slug")
	commentID, err := domain.ParseID(c.Param("comment_id"))
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

	err = server.service.Article().DeleteComment(c, port.DeleteCommentParams{
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

func (server *Server) AddFavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}

	article, err := server.service.Article().AddFavorite(c, port.AddFavoriteParams{
		AuthArg: authArg,
		Slug:    slug,
		UserID:  authArg.Payload.UserID,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := ArticleResponse{
		Article: serializeArticle(article),
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

	article, err := server.service.Article().RemoveFavorite(c, port.RemoveFavoriteParams{
		AuthArg: authArg,
		Slug:    slug,
		UserID:  authArg.Payload.UserID,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := ArticleResponse{serializeArticle(article)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) ListTags(c *gin.Context) {
	tags, err := server.service.Article().ListTags(c)
	if err != nil {
		errorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
