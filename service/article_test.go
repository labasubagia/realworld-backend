package service_test

import (
	"context"
	"sort"
	"testing"

	"github.com/labasubagia/realworld-backend/domain"
	"github.com/labasubagia/realworld-backend/port"
	"github.com/labasubagia/realworld-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateArticleOK(t *testing.T) {
	user, _ := createRandomUser(t)
	createRandomArticle(t, user)
}

func TestCreateArticleOkWithSameTags(t *testing.T) {
	ctx := context.Background()

	author, _ := createRandomUser(t)
	arg := createArticleArg(author)

	// 2 article with tags
	// only len(arg.Tags) expected in DB
	createArticle(t, arg)
	createArticle(t, arg)

	// 1 article without tags
	// no additional tags expected in DB
	arg2 := createArticleArg(author)
	arg2.Tags = []string{}
	createArticle(t, arg2)

	// tags to have len = len(arg.Tags)
	tags, err := testRepo.Article().FilterTags(ctx, port.FilterTagParams{Names: arg.Tags})
	require.Nil(t, err)
	require.Len(t, tags, len(arg.Tags))
}

func createRandomArticle(t *testing.T, author domain.User) port.CreateArticleTxResult {
	return createArticle(t, createArticleArg(author))
}

func createArticle(t *testing.T, arg port.CreateArticleTxParams) port.CreateArticleTxResult {
	result, err := testService.Article().Create(context.Background(), arg)
	resultTags := []string{}
	for _, tag := range result.Tags {
		resultTags = append(resultTags, tag.Name)
	}
	sort.Strings(arg.Tags)
	sort.Strings(resultTags)

	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, arg.Article.AuthorID, result.Article.AuthorID)
	require.Equal(t, arg.Article.Title, result.Article.Title)
	require.Equal(t, arg.Article.Slug, result.Article.Slug)
	require.Equal(t, arg.Article.Description, result.Article.Description)
	require.Equal(t, arg.Article.Body, result.Article.Body)
	require.Equal(t, arg.Tags, resultTags)
	require.Len(t, resultTags, len(arg.Tags))
	return result
}

func createArticleArg(author domain.User) port.CreateArticleTxParams {
	return port.CreateArticleTxParams{
		Article: domain.RandomArticle(author),
		Tags:    []string{util.RandomString(6), util.RandomString(7)},
	}
}
