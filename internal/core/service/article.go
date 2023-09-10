package service

import (
	"context"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

type articleService struct {
	property serviceProperty
}

func NewArticleService(property serviceProperty) port.ArticleService {
	return &articleService{
		property: property,
	}
}

func (s *articleService) Create(ctx context.Context, arg port.CreateArticleTxParams) (result port.CreateArticleTxResult, err error) {
	err = s.property.repo.Atomic(ctx, func(r port.Repository) error {
		// create article
		result.Article, err = r.Article().CreateArticle(ctx, port.CreateArticlePayload{Article: arg.Article})
		if err != nil {
			return exception.Into(err)
		}

		// return when no tags
		if len(arg.Tags) == 0 {
			return nil
		}

		// add tags if not exists
		for i, tag := range arg.Tags {
			arg.Tags[i] = strings.ToLower(tag)
		}
		result.Tags, err = r.Article().AddTagsIfNotExists(ctx, port.AddTagsPayload{Tags: arg.Tags})
		if err != nil {
			return exception.Into(err)
		}

		// assign tags
		tagIDs := []domain.ID{}
		for _, tag := range result.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		_, err := r.Article().AssignArticleTags(ctx, port.AssignTagPayload{ArticleID: result.Article.ID, TagIDs: tagIDs})
		if err != nil {
			return exception.Into(err)
		}
		return nil
	})
	if err != nil {
		return port.CreateArticleTxResult{}, exception.Into(err)
	}
	return result, nil
}

func (s *articleService) AddFavorite(ctx context.Context, arg port.AddFavoriteParams) (result port.AddFavoriteResult, err error) {
	if arg.AuthArg.Payload == nil {
		return port.AddFavoriteResult{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}
	article, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return port.AddFavoriteResult{}, exception.Into(err)
	}

	addFavorite, err := s.property.repo.Article().AddFavorite(ctx, domain.ArticleFavorite{
		ArticleID: article.ID,
		UserID:    arg.AuthArg.Payload.UserID,
	})
	if err != nil {
		return port.AddFavoriteResult{}, exception.Into(err)
	}

	list, err := s.List(ctx, port.ListArticleParams{
		AuthArg: arg.AuthArg,
		IDs:     []domain.ID{addFavorite.ArticleID},
	})
	if len(list.Articles) < 1 {
		return port.AddFavoriteResult{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}
	result.Article = list.Articles[0]

	return result, nil
}

func (s *articleService) List(ctx context.Context, arg port.ListArticleParams) (result port.ListArticleResult, err error) {

	authorIDs := []domain.ID{}
	authorMap := map[domain.ID]domain.User{}
	if len(arg.AuthorNames) > 0 {
		authors, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
			Usernames: arg.AuthorNames,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		for _, author := range authors {
			authorIDs = append(authorIDs, author.ID)
			authorMap[author.ID] = author
		}
	}

	taggedArticleIDs := []domain.ID{}
	if len(arg.Tags) > 0 {

		// find tags
		tags, err := s.property.repo.Article().FilterTags(ctx, port.FilterTagPayload{Names: arg.Tags})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		tagIDs := []domain.ID{}
		for _, tag := range tags {
			tagIDs = append(tagIDs, tag.ID)
		}

		// find article tags
		articleTags, err := s.property.repo.Article().FilterArticleTags(ctx, port.FilterArticleTagPayload{
			TagIDs: tagIDs,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		for _, articleTag := range articleTags {
			taggedArticleIDs = append(taggedArticleIDs, articleTag.ArticleID)
		}
	}

	favoritedArticleIDs := []domain.ID{}
	if len(arg.FavoritedNames) > 0 {
		// find users
		users, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
			Usernames: arg.FavoritedNames,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		userIDs := []domain.ID{}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}

		// find article ids
		favorites, err := s.property.repo.Article().FilterFavorite(ctx, port.FilterFavoritePayload{
			UserIDs: userIDs,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		for _, favorite := range favorites {
			favoritedArticleIDs = append(favoritedArticleIDs, favorite.ArticleID)
		}
	}

	// get articles
	result.Articles, err = s.property.repo.Article().FilterArticle(ctx, port.FilterArticlePayload{
		IDs:       append(arg.IDs, append(taggedArticleIDs, favoritedArticleIDs...)...),
		AuthorIDs: authorIDs,
		Limit:     arg.Limit,
		Offset:    arg.Offset,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}
	result.Count = len(result.Articles)
	for i, article := range result.Articles {
		if author, ok := authorMap[article.AuthorID]; ok {
			result.Articles[i].Author = author
		}
	}

	return result, nil
}
