package service

import "github.com/labasubagia/go-backend-realworld/port"

type service struct {
	repo           port.Repository
	articleService port.ArticleService
	userService    port.UserService
}

func NewService(repo port.Repository) port.Service {
	return &service{
		repo:           repo,
		articleService: NewArticleService(repo),
	}
}

func (s *service) Article(repo port.Repository) port.ArticleService {
	return s.articleService
}

// User implements port.Service.
func (*service) User(repo port.Repository) port.UserService {
	panic("unimplemented")
}
