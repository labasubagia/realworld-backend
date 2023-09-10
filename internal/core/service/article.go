package service

import (
	"context"
	"fmt"
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
	if arg.AuthArg.Payload == nil {
		return port.CreateArticleTxResult{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	err = s.property.repo.Atomic(ctx, func(r port.Repository) error {

		arg.Article.AuthorID = arg.AuthArg.Payload.UserID
		newArticle := domain.NewArticle(arg.Article)

		// create article
		result.Article, err = r.Article().CreateArticle(ctx, port.CreateArticlePayload{Article: newArticle})
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

	// Get decorated article
	articleInfos, err := s.infoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.AuthArg,
		articles: []domain.Article{result.Article},
	})
	if err != nil {
		return port.CreateArticleTxResult{}, exception.Into(err)
	}
	if len(articleInfos.Articles) == 0 {
		return port.CreateArticleTxResult{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}
	result.Article = articleInfos.Articles[0]

	return result, nil
}

func (s *articleService) Update(ctx context.Context, arg port.UpdateArticleParams) (port.UpdateArticleResult, error) {
	if arg.AuthArg.Payload == nil {
		return port.UpdateArticleResult{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	current, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs:     []string{arg.Slug},
		AuthorIDs: []domain.ID{arg.AuthArg.Payload.UserID},
	})
	if err != nil {
		return port.UpdateArticleResult{}, exception.Into(err)
	}

	current.Title = arg.Article.Title
	current.Description = arg.Article.Description
	current.Body = arg.Article.Body

	updated, err := s.property.repo.Article().UpdateArticle(ctx, port.UpdateArticlePayload{Article: current})
	if err != nil {
		return port.UpdateArticleResult{}, exception.Into(err)
	}

	// Get decorated article
	articleInfos, err := s.infoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.AuthArg,
		articles: []domain.Article{updated},
	})
	if err != nil {
		return port.UpdateArticleResult{}, exception.Into(err)
	}
	if len(articleInfos.Articles) == 0 {
		return port.UpdateArticleResult{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}

	result := port.UpdateArticleResult{
		Article: articleInfos.Articles[0],
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

func (s *articleService) Get(ctx context.Context, arg port.GetArticleParams) (port.GetArticleResult, error) {
	articles, err := s.property.repo.Article().FilterArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return port.GetArticleResult{}, exception.Into(err)
	}
	articleInfos, err := s.infoArticles(ctx, GetListArticleInfoParams{authArg: arg.AuthArg, articles: articles})
	if err != nil {
		return port.GetArticleResult{}, exception.Into(err)
	}
	if len(articleInfos.Articles) < 1 {
		return port.GetArticleResult{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}

	result := port.GetArticleResult{
		Article: articleInfos.Articles[0],
	}
	return result, nil
}

func (s *articleService) List(ctx context.Context, arg port.ListArticleParams) (result port.ListArticleResult, err error) {

	// filter authors
	authorIDs := []domain.ID{}
	if len(arg.AuthorNames) > 0 {
		authors, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
			Usernames: arg.AuthorNames,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		if len(authors) == 0 {
			authorNames := strings.Join(arg.AuthorNames, ", ")
			msg := fmt.Sprintf("author %s does not exists", authorNames)
			return port.ListArticleResult{}, exception.Validation().AddError("author", msg)
		}
		for _, author := range authors {
			authorIDs = append(authorIDs, author.ID)
		}
	}

	// filter tags
	taggedArticleIDs := []domain.ID{}
	if len(arg.Tags) > 0 {

		// find tags
		tags, err := s.property.repo.Article().FilterTags(ctx, port.FilterTagPayload{Names: arg.Tags})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		if len(tags) == 0 {
			tagNames := strings.Join(arg.Tags, ", ")
			msg := fmt.Sprintf("tag %s does not exists", tagNames)
			return port.ListArticleResult{}, exception.Validation().AddError("tag", msg)
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

	// filter favorites by users
	favoritedArticleIDs := []domain.ID{}
	if len(arg.FavoritedNames) > 0 {
		// find users
		users, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
			Usernames: arg.FavoritedNames,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		if len(users) == 0 {
			userNames := strings.Join(arg.FavoritedNames, ", ")
			msg := fmt.Sprintf("user for favorite filter %s does not exists", userNames)
			return port.ListArticleResult{}, exception.Validation().AddError("favorited", msg)
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

	// Get articles
	articles, err := s.property.repo.Article().FilterArticle(ctx, port.FilterArticlePayload{
		IDs:       append(arg.IDs, append(taggedArticleIDs, favoritedArticleIDs...)...),
		AuthorIDs: authorIDs,
		Limit:     arg.Limit,
		Offset:    arg.Offset,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}

	result, err = s.infoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.AuthArg,
		articles: articles,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}

	return result, nil
}

func (s *articleService) Feed(ctx context.Context, arg port.ListArticleParams) (result port.ListArticleResult, err error) {

	if arg.AuthArg.Payload == nil {
		return result, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	followingAuthors, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}
	if len(followingAuthors) == 0 {
		withInfos := port.ListArticleResult{
			Articles: []domain.Article{},
			Count:    0,
		}
		return withInfos, nil
	}

	authorIDs := []domain.ID{}
	for _, author := range followingAuthors {
		authorIDs = append(authorIDs, author.FolloweeID)
	}

	articles, err := s.property.repo.Article().FilterArticle(ctx, port.FilterArticlePayload{
		AuthorIDs: authorIDs,
		Limit:     arg.Limit,
		Offset:    arg.Offset,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}

	result, err = s.infoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.AuthArg,
		articles: articles,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}

	return result, nil
}

type GetListArticleInfoParams struct {
	authArg  port.AuthParams
	articles []domain.Article
}

// infoArticles add decorator to articles
func (s *articleService) infoArticles(ctx context.Context, arg GetListArticleInfoParams) (port.ListArticleResult, error) {

	// Get article id and author id
	authorIDs := []domain.ID{}
	articleIDs := []domain.ID{}
	for _, article := range arg.articles {
		authorIDs = append(authorIDs, article.AuthorID)
		articleIDs = append(articleIDs, article.ID)
	}

	// Get article authors
	authors, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
		IDs: authorIDs,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}
	authorMap := map[domain.ID]domain.User{}
	for _, author := range authors {
		authorMap[author.ID] = author
	}

	// Logged user check if followed author
	loggedUserFollowedAuthorMap := map[domain.ID]bool{}
	if arg.authArg.Payload != nil {
		followed, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
			FollowerIDs: []domain.ID{arg.authArg.Payload.UserID},
			FolloweeIDs: authorIDs,
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		for _, follow := range followed {
			loggedUserFollowedAuthorMap[follow.FolloweeID] = true
		}
	}

	// Get article tags
	articleTags, err := s.property.repo.Article().FilterArticleTags(ctx, port.FilterArticleTagPayload{
		ArticleIDs: articleIDs,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}
	tagIDs := []domain.ID{}
	for _, articleTag := range articleTags {
		tagIDs = append(tagIDs, articleTag.TagID)
	}
	tags, err := s.property.repo.Article().FilterTags(ctx, port.FilterTagPayload{
		IDs: tagIDs,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}
	tagMap := map[domain.ID]string{}
	for _, tag := range tags {
		tagMap[tag.ID] = tag.Name
	}
	// compose article tags
	articleTagMap := map[domain.ID][]string{} // article_id:[]tag_name
	for _, articleTag := range articleTags {
		articleTagMap[articleTag.ArticleID] = append(articleTagMap[articleTag.ArticleID], tagMap[articleTag.TagID])
	}

	// Get favorite counts
	favoriteCounts, err := s.property.repo.Article().FilterFavoriteCount(ctx, port.FilterFavoritePayload{
		ArticleIDs: articleIDs,
	})
	if err != nil {
		return port.ListArticleResult{}, exception.Into(err)
	}
	favoriteCountMap := map[domain.ID]int{}
	for _, favorite := range favoriteCounts {
		favoriteCountMap[favorite.ArticleID] = favorite.Count
	}

	// Get favorites by logged user
	loggedUserFavoritedArticleMap := map[domain.ID]bool{}
	if arg.authArg.Payload != nil {
		favorites, err := s.property.repo.Article().FilterFavorite(ctx, port.FilterFavoritePayload{
			ArticleIDs: articleIDs,
			UserIDs:    []domain.ID{arg.authArg.Payload.UserID},
		})
		if err != nil {
			return port.ListArticleResult{}, exception.Into(err)
		}
		for _, favorited := range favorites {
			loggedUserFavoritedArticleMap[favorited.ArticleID] = true
		}
	}

	// Compose result
	for index, article := range arg.articles {
		article.TagNames = articleTagMap[article.ID]
		article.FavoriteCount = favoriteCountMap[article.ID]
		article.IsFavorite = loggedUserFavoritedArticleMap[article.ID]
		if author, ok := authorMap[article.AuthorID]; ok {
			author.IsFollowed = loggedUserFollowedAuthorMap[author.ID]
			article.Author = author
		}
		arg.articles[index] = article
	}

	result := port.ListArticleResult{
		Articles: arg.articles,
		Count:    len(arg.articles),
	}
	return result, nil
}
