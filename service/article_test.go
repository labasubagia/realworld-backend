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

func createRandomArticle(t *testing.T) domain.Article {
	user, _ := createRandomUser(t)
	arg := port.CreateArticleTxParams{
		Article: domain.Article{
			AuthorID:    user.ID,
			Title:       util.RandomString(10),
			Description: util.RandomString(15),
			Slug:        util.RandomString(5),
			Body:        util.RandomString(20),
		},
		Tags: []string{util.RandomString(6), util.RandomString(7)},
	}
	result, err := testService.Article().Create(context.Background(), arg)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, user.ID, result.Article.AuthorID)
	require.Equal(t, arg.Article.Title, result.Article.Title)
	require.Equal(t, arg.Article.Slug, result.Article.Slug)
	require.Equal(t, arg.Article.Description, result.Article.Description)
	require.Equal(t, arg.Article.Body, result.Article.Body)

	resultTags := []string{}
	for _, tag := range result.Tags {
		resultTags = append(resultTags, tag.Name)
	}
	sort.Strings(arg.Tags)
	sort.Strings(resultTags)
	require.Equal(t, arg.Tags, resultTags)
	return result.Article
}

func TestCreateArticleOK(t *testing.T) {
	createRandomArticle(t)
}
