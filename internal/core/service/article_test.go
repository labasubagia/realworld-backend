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

func TestCreateArticleUnauthenticated(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	arg := createArticleArg(author, authorAuth)
	result, err := testService.Article().Create(context.Background(), port.CreateArticleTxParams{
		Article: arg.Article,
	})
	require.NotNil(t, err)
	require.Empty(t, result)
	fail, ok := err.(*exception.Exception)
	require.True(t, ok)
	require.Equal(t, exception.TypePermissionDenied, fail.Type)
}

func TestCreateArticleOK(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	createRandomArticle(t, author, authorAuth)
}

func TestCreateArticleOkWithTags(t *testing.T) {
	ctx := context.Background()

	author, authorAuth, _ := createRandomUser(t)
	arg1 := createArticleArg(author, authorAuth)

	// 2 article with same tags
	// new N=length(arg1.Tags) tags expected in DB
	createArticle(t, arg1)
	createArticle(t, arg1)

	// 1 article without tags
	// no additional tags expected in DB
	arg2 := createArticleArg(author, authorAuth)
	arg2.Tags = []string{}
	createArticle(t, arg2)

	// 1 article with different tags
	// new tags expected in DB
	arg3 := createArticleArg(author, authorAuth)
	createArticle(t, arg3)

	// tags to have N=len(allTags) tags in DB
	allTags := append(arg1.Tags, append(arg2.Tags, arg3.Tags...)...)
	tags, err := testRepo.Article().FilterTags(ctx, port.FilterTagPayload{Names: allTags})
	require.Nil(t, err)
	require.Len(t, tags, len(allTags))
}

func TestCreateArticleConcurrentOK(t *testing.T) {
	// make valid data
	author, authorAuth, _ := createRandomUser(t)
	arg := createArticleArg(author, authorAuth)

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

func TestUpdateArticle(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	ctx := context.Background()

	newTitle := util.RandomString(3)
	newBody := util.RandomString(5)
	newDescription := util.RandomString(12)

	t.Run("Update", func(t *testing.T) {
		article := createRandomArticle(t, author, authorAuth)
		result, err := testService.Article().Update(ctx, port.UpdateArticleParams{
			AuthArg: authorAuth,
			Slug:    article.Slug,
			Article: domain.Article{
				Title:       newTitle,
				Description: newDescription,
				Body:        newBody,
			},
		})
		require.Nil(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, article.ID, result.ID)
		require.Equal(t, newTitle, result.Title)
		require.NotEqual(t, article.Slug, result.Slug)
		require.Equal(t, newDescription, result.Description)
		require.Equal(t, newBody, result.Body)
	})

	t.Run("Partial", func(t *testing.T) {
		article := createRandomArticle(t, author, authorAuth)
		result, err := testService.Article().Update(ctx, port.UpdateArticleParams{
			AuthArg: authorAuth,
			Slug:    article.Slug,
			Article: domain.Article{
				Title: newTitle,
			},
		})
		require.Nil(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, article.ID, result.ID)
		require.Equal(t, newTitle, result.Title)
		require.NotEqual(t, article.Slug, result.Slug)
		require.Equal(t, article.Description, result.Description)
		require.Equal(t, article.Body, result.Body)
	})

	t.Run("Update other article fail", func(t *testing.T) {
		_, randomAuth, _ := createRandomUser(t)

		article := createRandomArticle(t, author, authorAuth)
		result, err := testService.Article().Update(ctx, port.UpdateArticleParams{
			AuthArg: randomAuth,
			Slug:    article.Slug,
			Article: domain.Article{
				Title:       newTitle,
				Description: newDescription,
				Body:        newBody,
			},
		})
		require.NotNil(t, err)
		require.Empty(t, result)
		fail, ok := err.(*exception.Exception)
		require.True(t, ok)
		require.Equal(t, exception.TypeNotFound, fail.Type)
	})
}

func TestDeleteArticle(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	article := createRandomArticle(t, author, authorAuth)
	ctx := context.Background()

	// other user try to delete fail
	_, otherUserAuth, _ := createRandomUser(t)
	err := testService.Article().Delete(ctx, port.DeleteArticleParams{
		AuthArg: otherUserAuth,
		Slug:    article.Slug,
	})
	require.NotNil(t, err)
	fail, ok := err.(*exception.Exception)
	require.True(t, ok)
	require.Equal(t, exception.TypeNotFound, fail.Type)

	// delete ok
	err = testService.Article().Delete(ctx, port.DeleteArticleParams{
		AuthArg: authorAuth,
		Slug:    article.Slug,
	})
	require.Nil(t, err)

	// not found
	getResult, err := testService.Article().Get(ctx, port.GetArticleParams{
		AuthArg: authorAuth,
		Slug:    article.Slug,
	})
	require.NotNil(t, err)
	require.Empty(t, getResult)
	fail, ok = err.(*exception.Exception)
	require.True(t, ok)
	require.Equal(t, exception.TypeNotFound, fail.Type)
}

func TestFavoriteArticle(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	reader, readerAuth, _ := createRandomUser(t)
	article := createArticle(t, createArticleArg(author, authorAuth))

	// add favorite
	addFav, err := testService.Article().AddFavorite(context.Background(), port.AddFavoriteParams{
		AuthArg: readerAuth,
		Slug:    article.Slug,
		UserID:  reader.ID,
	})
	require.Nil(t, err)
	require.Equal(t, article.Title, addFav.Title)
	require.Equal(t, article.ID, addFav.ID)
	require.True(t, addFav.IsFavorite)
	require.Equal(t, 1, addFav.FavoriteCount)

	// remove favorite
	removeFav, err := testService.Article().RemoveFavorite(context.Background(), port.RemoveFavoriteParams{
		AuthArg: readerAuth,
		Slug:    article.Slug,
		UserID:  reader.ID,
	})
	require.Nil(t, err)
	require.Equal(t, article.Title, removeFav.Title)
	require.Equal(t, article.ID, removeFav.ID)
	require.False(t, removeFav.IsFavorite)
	require.Equal(t, 0, removeFav.FavoriteCount)
}

func TestFeedArticle(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	_, reader1Auth, _ := createRandomUser(t)
	_, reader2Auth, _ := createRandomUser(t)
	ctx := context.Background()

	N := 5
	for i := 0; i < N; i++ {
		arg := createArticleArg(author, authorAuth)
		createArticle(t, arg)
	}

	t.Run("Feed", func(t *testing.T) {
		followResult, err := testService.User().Follow(ctx, port.ProfileParams{AuthArg: reader1Auth, Username: author.Username})
		require.Nil(t, err)
		require.NotEmpty(t, followResult)

		result, err := testService.Article().Feed(ctx, port.ListArticleParams{AuthArg: reader1Auth})
		require.Nil(t, err)
		require.NotEmpty(t, result)
		require.Len(t, result, N)
	})

	t.Run("Feed empty", func(t *testing.T) {
		result, err := testService.Article().Feed(ctx, port.ListArticleParams{AuthArg: reader2Auth})
		require.Nil(t, err)
		require.Len(t, result, 0)
	})
}

func TestListArticleOK(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	reader, readerAuth, _ := createRandomUser(t)

	tags := []string{util.RandomString(4), util.RandomString(5)}
	sort.Strings(tags)

	ctx := context.Background()

	// create other author (pollute)
	polluteN := 5
	otherAuthor, otherAuthorAuth, _ := createRandomUser(t)
	for i := 0; i < polluteN; i++ {
		arg := createArticleArg(otherAuthor, otherAuthorAuth)
		arg.Tags = []string{}
		createArticle(t, arg)
	}

	// create articles
	N := 10
	createdArticles := make([]domain.Article, N)
	for i := 0; i < N; i++ {
		arg := createArticleArg(author, authorAuth)
		arg.Tags = tags
		createResult := createArticle(t, arg)
		// assign in desc order
		createdArticles[N-1-i] = createResult
	}

	t.Run("Paginate", func(t *testing.T) {
		limit, offset := 3, 2
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			Limit:  limit,
			Offset: offset,
		})
		require.Nil(t, err)
		require.Len(t, result, limit)

		expectedArticleIDs, articleIDs := []domain.ID{}, []domain.ID{}
		for i := 0; i < limit; i++ {
			expectedArticleIDs = append(expectedArticleIDs, createdArticles[offset : offset+limit][i].ID)
			articleIDs = append(articleIDs, result[i].ID)
		}
		require.Equal(t, expectedArticleIDs, articleIDs)
	})

	t.Run("Filter by nonexistent author", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			AuthorNames: []string{"nonexistent_author"},
		})
		require.Nil(t, err)
		require.Empty(t, result)
	})

	t.Run("Filter by author", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			AuthorNames: []string{author.Username},
		})
		require.Nil(t, err)
		require.Len(t, result, N)
		for _, article := range result {
			require.Equal(t, author.ID, article.AuthorID)
			require.NotEqual(t, otherAuthor.ID, article.AuthorID)
		}
	})

	t.Run("Filter by nonexistent tag", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			Tags: []string{"nonexistent_tag"},
		})
		require.Nil(t, err)
		require.Empty(t, result)
	})

	t.Run("Filter by tags", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{Tags: tags})
		require.Nil(t, err)
		require.Len(t, result, N)
		for _, article := range result {
			sort.Strings(article.TagNames)
			require.Equal(t, tags, article.TagNames)
		}
	})

	t.Run("Filter by nonexistent favorites", func(t *testing.T) {
		result, err := testService.Article().List(ctx, port.ListArticleParams{
			FavoritedNames: []string{"nonexistent_fav_name"},
		})
		require.Nil(t, err)
		require.Empty(t, result)
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
		require.Len(t, result, favN)
		for _, article := range result {
			require.Greater(t, article.FavoriteCount, 0)
			require.True(t, article.IsFavorite)
		}
	})
}

func TestGetArticle(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	_, readerAuth, _ := createRandomUser(t)

	article := createRandomArticle(t, author, authorAuth)
	ctx := context.Background()

	result, err := testService.Article().Get(ctx, port.GetArticleParams{
		AuthArg: readerAuth,
		Slug:    article.Slug,
	})
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, article.Title, result.Title)
	require.Equal(t, article.Slug, result.Slug)
	require.Equal(t, article.Body, result.Body)
}

func TestCreateComment(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	_, user1Auth, _ := createRandomUser(t)
	_, user2Auth, _ := createRandomUser(t)

	article := createRandomArticle(t, author, authorAuth)
	ctx := context.Background()

	// user1
	arg1 := port.AddCommentParams{
		AuthArg: user1Auth,
		Slug:    article.Slug,
		Comment: domain.Comment{Body: util.RandomString(10)},
	}
	result, err := testService.Article().AddComment(ctx, arg1)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, arg1.Comment.Body, result.Body)
	require.Equal(t, article.ID, result.ArticleID)

	// user2
	arg2 := port.AddCommentParams{
		AuthArg: user2Auth,
		Slug:    article.Slug,
		Comment: domain.Comment{Body: util.RandomString(10)},
	}
	result, err = testService.Article().AddComment(ctx, arg2)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, arg2.Comment.Body, result.Body)
	require.Equal(t, article.ID, result.ArticleID)

	// get
	comments, err := testService.Article().ListComments(ctx, port.ListCommentParams{
		AuthArg: authorAuth,
		Slug:    article.Slug,
	})
	require.Nil(t, err)
	require.Len(t, comments, 2)
}

func TestDeleteComment(t *testing.T) {
	author, authorAuth, _ := createRandomUser(t)
	_, userAuth, _ := createRandomUser(t)

	article := createRandomArticle(t, author, authorAuth)
	ctx := context.Background()

	// user create comment
	arg := port.AddCommentParams{
		AuthArg: userAuth,
		Slug:    article.Slug,
		Comment: domain.Comment{Body: util.RandomString(10)},
	}
	result, err := testService.Article().AddComment(ctx, arg)
	require.Nil(t, err)
	require.NotEmpty(t, result)

	// author can see comment found
	comments, err := testService.Article().ListComments(ctx, port.ListCommentParams{
		AuthArg: authorAuth,
		Slug:    arg.Slug,
	})
	require.Nil(t, err)
	require.Len(t, comments, 1)

	// user delete his comment
	err = testService.Article().DeleteComment(ctx, port.DeleteCommentParams{
		AuthArg:   userAuth,
		Slug:      arg.Slug,
		CommentID: result.ID,
	})
	require.Nil(t, err)

	// author cannot see comment (gone)
	comments, err = testService.Article().ListComments(ctx, port.ListCommentParams{
		AuthArg: authorAuth,
		Slug:    arg.Slug,
	})
	require.Nil(t, err)
	require.Len(t, comments, 0)
}

func TestGetTags(t *testing.T) {
	ctx := context.Background()

	// create tags
	payload := port.AddTagsPayload{
		Tags: []string{
			util.RandomString(2),
			util.RandomString(3),
			util.RandomString(4),
		},
	}

	// no tags
	current, err := testService.Article().ListTags(ctx)
	require.Nil(t, err)

	// create tags
	_, err = testRepo.Article().AddTags(ctx, payload)
	require.Nil(t, err)

	// combine
	expected := append(current, payload.Tags...)
	sort.Strings(expected)

	// result
	actual, err := testService.Article().ListTags(ctx)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func createRandomArticle(t *testing.T, author domain.User, authArg port.AuthParams) domain.Article {
	return createArticle(t, createArticleArg(author, authArg))
}

func createArticle(t *testing.T, arg port.CreateArticleTxParams) domain.Article {
	result, err := testService.Article().Create(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, arg.Article.AuthorID, result.AuthorID)
	require.Equal(t, arg.Article.Title, result.Title)
	require.Equal(t, arg.Article.Slug, result.Slug)
	require.Equal(t, arg.Article.Description, result.Description)
	require.Equal(t, arg.Article.Body, result.Body)

	sort.Strings(arg.Tags)
	sort.Strings(result.TagNames)
	require.Equal(t, arg.Tags, result.TagNames)

	require.Len(t, result.TagNames, len(arg.Tags))
	return result
}

func createArticleArg(author domain.User, authorAuth port.AuthParams) port.CreateArticleTxParams {
	return port.CreateArticleTxParams{
		AuthArg: authorAuth,
		Article: domain.RandomArticle(author),
		Tags:    []string{util.RandomString(6), util.RandomString(7)},
	}
}
