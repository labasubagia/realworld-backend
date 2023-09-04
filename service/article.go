package service

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/port"
)

type articleService struct {
	repo port.Repository
}

func NewArticleService(repo port.Repository) port.ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (s *articleService) Create(ctx context.Context, arg port.CreateArticleTxParams) (result port.CreateArticleTxResult, err error) {
	err = s.repo.Atomic(ctx, func(r port.Repository) error {
		result.Article, err = r.Article().CreateArticle(ctx, port.CreateArticleParams{Article: arg.Article})
		if err != nil {
			return err
		}
		for _, item := range arg.Tags {
			arg := port.CreateTagParams{ArticleID: result.Article.ID, Tag: item}
			tag, err := r.Article().CreateTag(ctx, arg)
			if err != nil {
				return err
			}
			result.Tags = append(result.Tags, tag)
		}
		return nil
	})
	return result, err
}
