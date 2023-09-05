package service

import "github.com/labasubagia/realworld-backend/port"

type service struct {
	repo           port.Repository
	articleService port.ArticleService
	userService    port.UserService
}

func NewService(repo port.Repository) port.Service {
	return &service{
		repo:           repo,
		articleService: NewArticleService(repo),
		userService:    NewUserService(repo),
	}
}

func (s *service) Article() port.ArticleService {
	return s.articleService
}

func (s *service) User() port.UserService {
	return s.userService
}
