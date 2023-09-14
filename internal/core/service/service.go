package service

import (
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/labasubagia/realworld-backend/internal/core/util/token"
)

type serviceProperty struct {
	config     util.Config
	tokenMaker token.Maker
	repo       port.Repository
	logger     port.Logger
}

type services struct {
	property       serviceProperty
	articleService port.ArticleService
	userService    port.UserService
}

func NewService(config util.Config, repo port.Repository, logger port.Logger) (port.Service, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	property := serviceProperty{
		config:     config,
		repo:       repo,
		tokenMaker: tokenMaker,
		logger:     logger,
	}
	svc := services{
		property:       property,
		articleService: NewArticleService(property),
		userService:    NewUserService(property),
	}
	return &svc, nil
}

func (s *services) TokenMaker() token.Maker {
	return s.property.tokenMaker
}

func (s *services) Article() port.ArticleService {
	return s.articleService
}

func (s *services) User() port.UserService {
	return s.userService
}
