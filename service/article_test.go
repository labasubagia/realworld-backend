package service_test

import (
	"context"
	"sort"
	"sync"
	"testing"

	"github.com/labasubagia/realworld-backend/domain"
	"github.com/labasubagia/realworld-backend/port"
	"github.com/labasubagia/realworld-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateArticleOK(t *testing.T) {
	author, _ := createRandomUser(t)
	createRandomArticle(t, author)
}

func TestCreateArticleOkWithTags(t *testing.T) {
	ctx := context.Background()

	author, _ := createRandomUser(t)
	arg1 := createArticleArg(author)

	// 2 article with same tags
	// new N=length(arg1.Tags) tags expected in DB
	createArticle(t, arg1)
	createArticle(t, arg1)

	// 1 article without tags
	// no additional tags expected in DB
	arg2 := createArticleArg(author)
	arg2.Tags = []string{}
	createArticle(t, arg2)

	// 1 article with different tags
	// new tags expected in DB
	arg3 := createArticleArg(author)
	createArticle(t, arg3)

	// tags to have N=len(allTags) tags in DB
	allTags := append(arg1.Tags, append(arg2.Tags, arg3.Tags...)...)
	tags, err := testRepo.Article().FilterTags(ctx, port.FilterTagParams{Names: allTags})
	require.Nil(t, err)
	require.Len(t, tags, len(allTags))
}

func TestCreateArticleConcurrentOK(t *testing.T) {

	// make valid data
	author, _ := createRandomUser(t)
	arg := createArticleArg(author)

	// concurrent process
	N := 5
	wg := sync.WaitGroup{}
	chErrors := make(chan error, N)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := testService.Article().Create(context.Background(), arg)
			chErrors <- err
		}()
	}
	wg.Wait()
	close(chErrors)

	// if there is an error
	// it must be an isolation or unique error
	for range chErrors {
		err := <-chErrors
		if err != nil {
			require.Contains(t, []error{port.ErrIsolation, port.ErrUniqueKey}, err)
		}
	}

	// tags expected N=len(arg.Tags) in DB
	tags, err := testRepo.Article().FilterTags(context.Background(), port.FilterTagParams{Names: arg.Tags})
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
