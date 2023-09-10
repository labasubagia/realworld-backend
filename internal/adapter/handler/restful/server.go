package restful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

const formatTime string = "2006-01-02T15:04:05.999Z"

type Server struct {
	config  util.Config
	router  *gin.Engine
	service port.Service
}

func NewServer(config util.Config, service port.Service) port.Server {
	server := &Server{
		config:  config,
		service: service,
	}
	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {

	router := gin.Default()
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
	articleRouter.GET("/feed", server.Feed)

	server.router = router
}

func (server *Server) Start() error {
	srv := &http.Server{
		Addr:    server.config.HTTPServerAddress,
		Handler: server.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")
	}
	log.Println("Server existing")

	return nil
}
