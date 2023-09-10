package service_test

import (
	"context"
	"sort"
	"sync"
	"testing"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"github.com/stretchr/testify/require"
)

func TestCreateArticleOK(t *testing.T) {
	author, _, _ := createRandomUser(t)
	createRandomArticle(t, author)
}

func TestCreateArticleOkWithTags(t *testing.T) {
	ctx := context.Background()

	author, _, _ := createRandomUser(t)
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
	tags, err := testRepo.Article().FilterTags(ctx, port.FilterTagPayload{Names: allTags})
	require.Nil(t, err)
	require.Len(t, tags, len(allTags))
}

func TestCreateArticleConcurrentOK(t *testing.T) {

	// make valid data
	author, _, _ := createRandomUser(t)
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
	// it must be an validation (isolation or unique error)
	for range chErrors {
		err := <-chErrors
		if err != nil {
			fail, ok := err.(*exception.Exception)
			require.True(t, ok)
			require.Equal(t, exception.TypeValidation, fail.Type)
		}
	}

	// tags expected N=len(arg.Tags) in DB
	tags, err := testRepo.Article().FilterTags(context.Background(), port.FilterTagPayload{Names: arg.Tags})
	require.Nil(t, err)
	require.Len(t, tags, len(arg.Tags))
}

func TestAddFavoriteArticleOK(t *testing.T) {
	author, _, _ := createRandomUser(t)
	reader, readerAuth, _ := createRandomUser(t)

	resultCreateArticle := createArticle(t, createArticleArg(author))
	result, err := testService.Article().AddFavorite(context.Background(), port.AddFavoriteParams{
		AuthArg: readerAuth,
		Slug:    resultCreateArticle.Article.Slug,
		UserID:  reader.ID,
	})
	require.Nil(t, err)
	require.Equal(t, resultCreateArticle.Article.Title, result.Article.Title)
}

func TestFeedArticle(t *testing.T) {
	author, _, _ := createRandomUser(t)
	_, reader1Auth, _ := createRandomUser(t)
	_, reader2Auth, _ := createRandomUser(t)
	ctx := context.Background()

	N := 5
	for i := 0; i < N; i++ {
		arg := createArticleArg(author)
		createArticle(t, arg)
	}

	// follow author
	followResult, err := testService.User().Follow(ctx, port.ProfileParams{
		AuthArg:  reader1Auth,
		Username: author.Username,
	})
	require.Nil(t, err)
	require.NotEmpty(t, followResult)

	t.Run("Feed", func(t *testing.T) {
		result, err := testService.Article().Feed(ctx, port.ListArticleParams{AuthArg: reader1Auth})
		require.Nil(t, err)
		require.NotEmpty(t, result)
		require.Len(t, result.Articles, N)
		require.Equal(t, N, result.Count)
	})

	t.Run("Feed empty", func(t *testing.T) {
		result, err := testService.Article().Feed(ctx, port.ListArticleParams{AuthArg: reader2Auth})
		require.Nil(t, err)
		require.Len(t, result.Articles, 0)
		require.Equal(t, 0, result.Count)
	})
}

func TestListArticleOK(t *testing.T) {
	author, _, _ := createRandomUser(t)
	reader, readerAuth, _ := createRandomUser(t)

	tags := []string{util.RandomString(4), util.RandomString(5)}
	sort.Strings(tags)

	ctx := context.Background()

	// create other author (pollute)
	polluteN := 5
	otherAuthor, _, _ := createRandomUser(t)
	for i := 0; i < polluteN; i++ {
		arg := createArticleArg(otherAuthor)
		arg.Tags = []string{}
		createArticle(t, arg)
	}

	// create articles
	N := 10
	createdArticles := make([]domain.Article, N)
	for i := 0; i < N; i++ {
		arg := createArticleArg(author)
		arg.Tags = tags
		createResult := createArticle(t, arg)
		// assign in desc order
		createdArticles[N-1-i] = createResult.Article
	}

	t.Run("Paginate", func(t *testing.T) {
		limit, offset := 3, 2
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			Limit:  limit,
			Offset: offset,
		})
		require.Nil(t, err)
		require.Equal(t, limit, result.Count)
		require.Len(t, result.Articles, limit)

		expectedArticleIDs, articleIDs := []domain.ID{}, []domain.ID{}
		for i := 0; i < limit; i++ {
			expectedArticleIDs = append(expectedArticleIDs, createdArticles[offset : offset+limit][i].ID)
			articleIDs = append(articleIDs, result.Articles[i].ID)
		}
		require.Equal(t, expectedArticleIDs, articleIDs)
	})

	t.Run("Filter by nonexistent author", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			AuthorNames: []string{"nonexistent_author"},
		})
		require.NotNil(t, err)
		require.Empty(t, result)
		fail, ok := err.(*exception.Exception)
		require.True(t, ok)
		require.Equal(t, exception.TypeValidation, fail.Type)
	})

	t.Run("Filter by author", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			AuthorNames: []string{author.Username},
		})
		require.Nil(t, err)
		require.Equal(t, N, result.Count)
		require.Len(t, result.Articles, N)
		for _, article := range result.Articles {
			require.Equal(t, author.ID, article.AuthorID)
			require.NotEqual(t, otherAuthor.ID, article.AuthorID)
		}
	})

	t.Run("Filter by nonexistent tag", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			Tags: []string{"nonexistent_tag"},
		})
		require.NotNil(t, err)
		require.Empty(t, result)
		fail, ok := err.(*exception.Exception)
		require.True(t, ok)
		require.Equal(t, exception.TypeValidation, fail.Type)
	})

	t.Run("Filter by tags", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{Tags: tags})
		require.Nil(t, err)
		require.Equal(t, N, result.Count)
		require.Len(t, result.Articles, N)
		for _, article := range result.Articles {
			sort.Strings(article.TagNames)
			require.Equal(t, tags, article.TagNames)
		}
	})

	t.Run("Filter by nonexistent favorites", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			FavoritedNames: []string{"nonexistent_fav_name"},
		})
		require.NotNil(t, err)
		require.Empty(t, result)
		fail, ok := err.(*exception.Exception)
		require.True(t, ok)
		require.Equal(t, exception.TypeValidation, fail.Type)
	})

	t.Run("Filter by favorites", func(t *testing.T) {
		favN := int(util.RandomInt(1, int64(N)))
		for i := 0; i < favN; i++ {
			article := createdArticles[i]
			favResult, err := testService.Article().AddFavorite(ctx, port.AddFavoriteParams{
				AuthArg: readerAuth,
				Slug:    article.Slug,
				UserID:  reader.ID,
			})
			require.Nil(t, err)
			require.NotEmpty(t, favResult)
		}

		result, err := testService.Article().List(ctx, port.ListArticleParams{
			AuthArg:        readerAuth,
			FavoritedNames: []string{reader.Username},
		})
		require.Nil(t, err)
		require.Equal(t, favN, result.Count)
		require.Len(t, result.Articles, favN)
		for _, article := range result.Articles {
			require.Greater(t, article.FavoriteCount, 0)
			require.True(t, article.IsFavorite)
		}
	})
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
