package service

import (
	"context"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/core/port"
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
		// create article
		result.Article, err = r.Article().CreateArticle(ctx, port.CreateArticleParams{Article: arg.Article})
		if err != nil {
			return err
		}

		// return when no tags
		if len(arg.Tags) == 0 {
			return nil
		}

		// add tags if not exists
		for i, tag := range arg.Tags {
			arg.Tags[i] = strings.ToLower(tag)
		}
		result.Tags, err = r.Article().AddTagsIfNotExists(ctx, port.AddTagsParams{Tags: arg.Tags})
		if err != nil {
			return err
		}

		// assign tags
		tagIDs := []int64{}
		for _, tag := range result.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		_, err := r.Article().AssignTags(ctx, port.AssignTags{ArticleID: result.Article.ID, TagIDs: tagIDs})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return port.CreateArticleTxResult{}, err
	}
	return result, nil
}
