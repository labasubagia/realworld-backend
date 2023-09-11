package service

import (
	"context"
	"sort"
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

func (s *articleService) Create(ctx context.Context, arg port.CreateArticleTxParams) (article domain.Article, err error) {
	if arg.AuthArg.Payload == nil {
		return domain.Article{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	err = s.property.repo.Atomic(ctx, func(r port.Repository) error {

		arg.Article.AuthorID = arg.AuthArg.Payload.UserID
		newArticle := domain.NewArticle(arg.Article)

		// create article
		article, err = r.Article().CreateArticle(ctx, newArticle)
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
		tags, err := r.Article().AddTags(ctx, port.AddTagsPayload{Tags: arg.Tags})
		if err != nil {
			return exception.Into(err)
		}

		// assign tags
		tagIDs := []domain.ID{}
		for _, tag := range tags {
			tagIDs = append(tagIDs, tag.ID)
		}
		_, err = r.Article().AssignArticleTags(ctx, port.AssignTagPayload{ArticleID: article.ID, TagIDs: tagIDs})
		if err != nil {
			return exception.Into(err)
		}
		return nil
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	return s.infoArticle(ctx, GetArticleInfoParams{authArg: arg.AuthArg, article: article})
}

func (s *articleService) Update(ctx context.Context, arg port.UpdateArticleParams) (domain.Article, error) {
	if arg.AuthArg.Payload == nil {
		return domain.Article{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	current, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs:     []string{arg.Slug},
		AuthorIDs: []domain.ID{arg.AuthArg.Payload.UserID},
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	current.Title = arg.Article.Title
	current.Description = arg.Article.Description
	current.Body = arg.Article.Body

	updated, err := s.property.repo.Article().UpdateArticle(ctx, current)
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	return s.infoArticle(ctx, GetArticleInfoParams{authArg: arg.AuthArg, article: updated})
}

func (s *articleService) Delete(ctx context.Context, arg port.DeleteArticleParams) error {
	if arg.AuthArg.Payload == nil {
		return exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	current, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs:     []string{arg.Slug},
		AuthorIDs: []domain.ID{arg.AuthArg.Payload.UserID},
	})
	if err != nil {
		return exception.Into(err)
	}

	err = s.property.repo.Article().DeleteArticle(ctx, current)
	if err != nil {
		return exception.Into(err)
	}
	return nil
}

func (s *articleService) AddFavorite(ctx context.Context, arg port.AddFavoriteParams) (result domain.Article, err error) {
	if arg.AuthArg.Payload == nil {
		return domain.Article{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}
	article, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	favorites, err := s.property.repo.Article().FilterFavorite(ctx, port.FilterFavoritePayload{
		ArticleIDs: []domain.ID{article.ID},
		UserIDs:    []domain.ID{arg.AuthArg.Payload.UserID},
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	if len(favorites) == 0 {
		_, err := s.property.repo.Article().AddFavorite(ctx, domain.ArticleFavorite{
			ArticleID: article.ID,
			UserID:    arg.AuthArg.Payload.UserID,
		})
		if err != nil {
			return domain.Article{}, exception.Into(err)
		}
	}

	return s.infoArticle(ctx, GetArticleInfoParams{authArg: arg.AuthArg, article: article})
}

func (s *articleService) RemoveFavorite(ctx context.Context, arg port.RemoveFavoriteParams) (result domain.Article, err error) {
	if arg.AuthArg.Payload == nil {
		return domain.Article{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}
	article, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	_, err = s.property.repo.Article().RemoveFavorite(ctx, domain.ArticleFavorite{
		ArticleID: article.ID,
		UserID:    arg.AuthArg.Payload.UserID,
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}

	return s.infoArticle(ctx, GetArticleInfoParams{authArg: arg.AuthArg, article: article})
}

func (s *articleService) Get(ctx context.Context, arg port.GetArticleParams) (domain.Article, error) {
	articles, err := s.property.repo.Article().FilterArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}
	articleInfos, err := s.listInfoArticles(ctx, GetListArticleInfoParams{authArg: arg.AuthArg, articles: articles})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}
	if len(articleInfos) == 0 {
		return domain.Article{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}

	return articleInfos[0], nil
}

func (s *articleService) List(ctx context.Context, arg port.ListArticleParams) (result []domain.Article, err error) {

	// filter authors
	authorIDs := []domain.ID{}
	if len(arg.AuthorNames) > 0 {
		authors, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
			Usernames: arg.AuthorNames,
		})
		if err != nil {
			return []domain.Article{}, exception.Into(err)
		}

		if len(authors) == 0 {
			return []domain.Article{}, nil
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
			return []domain.Article{}, exception.Into(err)
		}

		if len(tags) == 0 {
			return []domain.Article{}, nil
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
			return []domain.Article{}, exception.Into(err)
		}
		if len(articleTags) == 0 {
			return []domain.Article{}, nil
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
			return []domain.Article{}, exception.Into(err)
		}
		if len(users) == 0 {
			return []domain.Article{}, nil
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
			return []domain.Article{}, exception.Into(err)
		}
		if len(favorites) == 0 {
			return []domain.Article{}, nil
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
		return []domain.Article{}, exception.Into(err)
	}

	result, err = s.listInfoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.AuthArg,
		articles: articles,
	})
	if err != nil {
		return []domain.Article{}, exception.Into(err)
	}

	return result, nil
}

func (s *articleService) Feed(ctx context.Context, arg port.ListArticleParams) (result []domain.Article, err error) {

	if arg.AuthArg.Payload == nil {
		return result, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	followingAuthors, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
	})
	if err != nil {
		return []domain.Article{}, exception.Into(err)
	}
	if len(followingAuthors) == 0 {
		return []domain.Article{}, nil
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
		return []domain.Article{}, exception.Into(err)
	}

	result, err = s.listInfoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.AuthArg,
		articles: articles,
	})
	if err != nil {
		return []domain.Article{}, exception.Into(err)
	}

	return result, nil
}

type GetArticleInfoParams struct {
	authArg port.AuthParams
	article domain.Article
}

func (s *articleService) infoArticle(ctx context.Context, arg GetArticleInfoParams) (domain.Article, error) {
	listInfos, err := s.listInfoArticles(ctx, GetListArticleInfoParams{
		authArg:  arg.authArg,
		articles: []domain.Article{arg.article},
	})
	if err != nil {
		return domain.Article{}, exception.Into(err)
	}
	if len(listInfos) == 0 {
		return domain.Article{}, exception.New(exception.TypeNotFound, "article not found", nil)
	}
	return listInfos[0], nil
}

type GetListArticleInfoParams struct {
	authArg  port.AuthParams
	articles []domain.Article
}

// listInfoArticles add decorator to articles
func (s *articleService) listInfoArticles(ctx context.Context, arg GetListArticleInfoParams) ([]domain.Article, error) {

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
		return []domain.Article{}, exception.Into(err)
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
			return []domain.Article{}, exception.Into(err)
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
		return []domain.Article{}, exception.Into(err)
	}
	tagIDs := []domain.ID{}
	for _, articleTag := range articleTags {
		tagIDs = append(tagIDs, articleTag.TagID)
	}
	tags, err := s.property.repo.Article().FilterTags(ctx, port.FilterTagPayload{
		IDs: tagIDs,
	})
	if err != nil {
		return []domain.Article{}, exception.Into(err)
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
		return []domain.Article{}, exception.Into(err)
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
			return []domain.Article{}, exception.Into(err)
		}
		for _, favorited := range favorites {
			loggedUserFavoritedArticleMap[favorited.ArticleID] = true
		}
	}

	// Compose result
	for index, article := range arg.articles {
		article.FavoriteCount = favoriteCountMap[article.ID]
		article.IsFavorite = loggedUserFavoritedArticleMap[article.ID]
		if author, ok := authorMap[article.AuthorID]; ok {
			author.IsFollowed = loggedUserFollowedAuthorMap[author.ID]
			article.Author = author
		}
		if tags, ok := articleTagMap[article.ID]; ok {
			sort.Strings(tags)
			article.TagNames = tags
		} else {
			article.TagNames = []string{}
		}
		arg.articles[index] = article
	}

	return arg.articles, nil
}

func (s *articleService) AddComment(ctx context.Context, arg port.AddCommentParams) (result domain.Comment, err error) {

	if arg.AuthArg.Payload == nil {
		return result, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	article, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return domain.Comment{}, exception.Into(err)
	}

	comment, err := s.property.repo.Article().AddComment(ctx, domain.NewComment(domain.Comment{
		ArticleID: article.ID,
		AuthorID:  arg.AuthArg.Payload.UserID,
		Body:      arg.Comment.Body,
	}))
	if err != nil {
		return domain.Comment{}, exception.Into(err)
	}

	// Get decorator info
	comments, err := s.listInfoComments(ctx, GetCommentInfo{
		authArg:  arg.AuthArg,
		comments: []domain.Comment{comment},
	})
	if err != nil {
		return domain.Comment{}, exception.Into(err)
	}
	if len(comments) == 0 {
		return domain.Comment{}, exception.New(exception.TypeNotFound, "comment not found", nil)
	}

	return comments[0], nil
}

func (s *articleService) ListComments(ctx context.Context, arg port.ListCommentParams) (result []domain.Comment, err error) {

	article, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return []domain.Comment{}, exception.Into(err)
	}

	comments, err := s.property.repo.Article().FilterComment(ctx, port.FilterCommentPayload{
		ArticleIDs: []domain.ID{article.ID},
	})
	if err != nil {
		return []domain.Comment{}, exception.Into(err)
	}

	// Get decorator info
	result, err = s.listInfoComments(ctx, GetCommentInfo{
		authArg:  arg.AuthArg,
		comments: comments,
	})
	if err != nil {
		return []domain.Comment{}, exception.Into(err)
	}

	return result, nil
}

func (s *articleService) DeleteComment(ctx context.Context, arg port.DeleteCommentParams) error {
	if arg.AuthArg.Payload == nil {
		return exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	article, err := s.property.repo.Article().FindOneArticle(ctx, port.FilterArticlePayload{
		Slugs: []string{arg.Slug},
	})
	if err != nil {
		return exception.Into(err)
	}

	err = s.property.repo.Article().DeleteComment(ctx, domain.Comment{
		ArticleID: article.ID,
		AuthorID:  arg.AuthArg.Payload.UserID,
		ID:        arg.CommentID,
	})
	if err != nil {
		return exception.Into(err)
	}

	return nil
}

type GetCommentInfo struct {
	authArg  port.AuthParams
	comments []domain.Comment
}

func (s *articleService) listInfoComments(ctx context.Context, arg GetCommentInfo) ([]domain.Comment, error) {

	// Get article id and author id
	authorIDs := []domain.ID{}
	for _, article := range arg.comments {
		authorIDs = append(authorIDs, article.AuthorID)
	}

	// Get article authors
	authors, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{
		IDs: authorIDs,
	})
	if err != nil {
		return []domain.Comment{}, exception.Into(err)
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
			return []domain.Comment{}, exception.Into(err)
		}
		for _, follow := range followed {
			loggedUserFollowedAuthorMap[follow.FolloweeID] = true
		}
	}

	for index, comment := range arg.comments {
		if author, ok := authorMap[comment.AuthorID]; ok {
			author.IsFollowed = loggedUserFollowedAuthorMap[comment.AuthorID]
			comment.Author = author
			arg.comments[index] = comment
		}
	}

	return arg.comments, nil
}

func (s *articleService) ListTags(ctx context.Context) ([]string, error) {
	tags, err := s.property.repo.Article().FilterTags(ctx, port.FilterTagPayload{})
	if err != nil {
		return []string{}, exception.Into(err)
	}
	tagNames := []string{}
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames, nil
}
