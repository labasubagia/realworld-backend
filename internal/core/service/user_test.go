package service_test

import (
	"context"
	"testing"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/stretchr/testify/require"
)

func TestRegisterOK(t *testing.T) {
	createRandomUser(t)
}

func TestRegisterWithImageOK(t *testing.T) {
	arg := createUserArg()
	arg.User.SetImageURL(util.RandomURL())
	createUser(t, arg)
}

func TestLoginOK(t *testing.T) {
	createRandomLogin(t)
}

func TestLoginInvalid(t *testing.T) {
	result, err := testService.User().Login(context.Background(), port.LoginUserParams{
		User: domain.RandomUser(),
	})
	require.NotNil(t, err)
	require.Empty(t, result)
}

func TestCurrentUserOK(t *testing.T) {
	user, authArg, _ := createRandomUser(t)

	result, err := testService.User().Current(context.Background(), authArg)
	require.Nil(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user, result.User)
	require.Equal(t, authArg.Token, result.Token)
}

func TestUpdateUserOK(t *testing.T) {
	user, authArg, _ := createRandomUser(t)

	newEmail := util.RandomEmail()
	newUsername := util.RandomUsername()
	newPassword := util.RandomString(8)
	newImage := util.RandomURL()
	newBio := util.RandomString(5)

	result, err := testService.User().Update(context.Background(), port.UpdateUserParams{
		AuthArg: authArg,
		User: domain.User{
			ID:       user.ID,
			Email:    newEmail,
			Username: newUsername,
			Password: newPassword,
			Image:    newImage,
			Bio:      newBio,
		},
	})
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, newEmail, result.User.Email)
	require.Equal(t, newUsername, result.User.Username)
	require.Equal(t, newImage, result.User.Image)
	require.Equal(t, newBio, result.User.Bio)
	require.Nil(t, util.CheckPassword(newPassword, result.User.Password))
}

func TestUpdateUserSameDataOK(t *testing.T) {
	user, authArg, password := createRandomUser(t)

	result, err := testService.User().Update(context.Background(), port.UpdateUserParams{
		AuthArg: authArg,
		User: domain.User{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Password: password,
			Image:    user.Image,
			Bio:      user.Bio,
		},
	})
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, user.Email, result.User.Email)
	require.Equal(t, user.Username, result.User.Username)
	require.Equal(t, user.Image, result.User.Image)
	require.Equal(t, user.Bio, result.User.Bio)
	require.Nil(t, util.CheckPassword(password, result.User.Password))
}

func TestUpdateUserEmptyOK(t *testing.T) {
	user, authArg, password := createRandomUser(t)

	result, err := testService.User().Update(context.Background(), port.UpdateUserParams{
		AuthArg: authArg,
		User: domain.User{
			ID:       user.ID,
			Email:    "",
			Username: "",
			Password: "",
			Image:    "",
			Bio:      "",
		},
	})
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, user.Email, result.User.Email)
	require.Equal(t, user.Username, result.User.Username)
	require.Equal(t, user.Image, result.User.Image)
	require.Equal(t, user.Bio, result.User.Bio)
	require.Nil(t, util.CheckPassword(password, result.User.Password))
}

func TestProfile(t *testing.T) {
	user, authArg, _ := createRandomUser(t)
	result, err := testService.User().Profile(context.Background(), port.ProfileParams{
		Username: user.Username,
		AuthArg:  authArg,
	})
	require.Nil(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, user.Email, result.User.Email)
	require.Equal(t, user.Username, result.User.Username)
	require.Equal(t, user.Image, result.User.Image)
	require.Equal(t, user.Bio, result.User.Bio)
	require.False(t, result.IsFollow)
}

func TestFollowUnFollow(t *testing.T) {
	followee, _, _ := createRandomUser(t)
	_, followerAuthArg, _ := createRandomUser(t)

	arg := port.ProfileParams{
		Username: followee.Username,
		AuthArg:  followerAuthArg,
	}

	// follow
	followResult, err := testService.User().Follow(context.Background(), arg)
	require.Nil(t, err)
	require.True(t, followResult.IsFollow)

	profileResult, err := testService.User().Profile(context.Background(), arg)
	require.Nil(t, err)
	require.True(t, profileResult.IsFollow)

	// un follow
	unFollowResult, err := testService.User().UnFollow(context.Background(), arg)
	require.Nil(t, err)
	require.False(t, unFollowResult.IsFollow)

	profileResult, err = testService.User().Profile(context.Background(), arg)
	require.Nil(t, err)
	require.False(t, profileResult.IsFollow)

}

func createRandomLogin(t *testing.T) {
	user, _, password := createRandomUser(t)
	createLogin(t, port.LoginUserParams{User: domain.User{
		Email:    user.Email,
		Password: password,
	}})
}

func createLogin(t *testing.T, arg port.LoginUserParams) (user domain.User, token, password string) {
	result, err := testService.User().Login(context.Background(), port.LoginUserParams{
		User: domain.User{
			Email:    arg.User.Email,
			Password: arg.User.Password,
		},
	})
	require.Nil(t, err)
	require.NotEmpty(t, result)

	payload, err := testService.TokenMaker().VerifyToken(result.Token)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, result.User.ID, payload.UserID)

	return user, result.Token, password
}

func createRandomUser(t *testing.T) (user domain.User, authArg port.AuthParams, password string) {
	return createUser(t, createUserArg())
}

func createUser(t *testing.T, arg port.RegisterUserParams) (user domain.User, authArg port.AuthParams, password string) {
	image := arg.User.Image
	if image == "" {
		image = domain.UserDefaultImage
	}

	result, err := testService.User().Register(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, result)

	user = result.User
	require.Equal(t, arg.User.Email, user.Email)
	require.Equal(t, arg.User.Username, user.Username)
	require.Equal(t, image, user.Image)
	require.NotEqual(t, arg.User.Password, user.Password)
	require.Nil(t, util.CheckPassword(arg.User.Password, user.Password))

	payload, err := testService.TokenMaker().VerifyToken(result.Token)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, user.ID, payload.UserID)

	authArg = port.AuthParams{
		Token:   result.Token,
		Payload: payload,
	}
	return user, authArg, arg.User.Password
}

func createUserArg() port.RegisterUserParams {
	return port.RegisterUserParams{
		User: domain.RandomUser(),
	}
}
