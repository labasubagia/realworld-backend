package restful

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

type Server struct {
	config  util.Config
	router  *gin.Engine
	service port.Service
	logger  port.Logger
}

func NewServer(config util.Config, service port.Service, logger port.Logger) port.Server {
	server := &Server{
		config:  config,
		service: service,
		logger:  logger,
	}
	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(server.Logger(), gin.Recovery(), cors.Default())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})
	router.POST("/users", server.Register)
	router.POST("/users/login", server.Login)

	userRouter := router.Group("/user")
	userRouter.Use(server.AuthMiddleware(true))
	userRouter.GET("/", server.CurrentUser)
	userRouter.PUT("/", server.UpdateUser)

	profileRouter := router.Group("/profiles/:username")
	profileRouter.Use(server.AuthMiddleware(false))
	profileRouter.GET("/", server.Profile)
	profileRouter.POST("/follow", server.FollowUser)
	profileRouter.DELETE("/follow", server.UnFollowUser)

	articleRouter := router.Group("/articles")
	articleRouter.Use(server.AuthMiddleware(false))
	articleRouter.GET("/", server.ListArticle)
	articleRouter.GET("/feed", server.FeedArticle)
	articleRouter.GET("/:slug", server.GetArticle)
	articleRouter.POST("/", server.CreateArticle)
	articleRouter.PUT("/:slug", server.UpdateArticle)
	articleRouter.DELETE("/:slug", server.DeleteArticle)

	commentRouter := articleRouter.Group("/:slug")
	commentRouter.POST("/comments", server.AddComment)
	commentRouter.GET("/comments", server.ListComments)
	commentRouter.DELETE("/comments/:comment_id", server.DeleteComment)

	favoriteArticleRouter := articleRouter.Group("/:slug/favorite")
	favoriteArticleRouter.POST("/", server.AddFavoriteArticle)
	favoriteArticleRouter.DELETE("/", server.RemoveFavoriteArticle)

	tagRouter := router.Group("/tags")
	tagRouter.GET("/", server.ListTags)

	server.router = router
}

func (server *Server) Start() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", server.config.HTTPServerPort),
		Handler: server.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			server.logger.Fatal().Err(err).Msg("failed listen")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.logger.Info().Msg("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		server.logger.Fatal().Err(err).Msg("failed shutdown server")
	}

	select {
	case <-ctx.Done():
		server.logger.Info().Msg("timeout in 5 seconds")
	}
	server.logger.Info().Msg("server exiting")

	return nil
}
