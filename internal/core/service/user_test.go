package service_test

import (
	"context"
	"testing"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	createRandomUser(t)
}

func TestRegisterWithImage(t *testing.T) {
	arg := createUserArg()
	arg.User.SetImageURL(util.RandomURL())
	createUser(t, arg)
}

func TestLogin(t *testing.T) {
	createRandomLogin(t)
}

func createRandomUser(t *testing.T) (user domain.User, token, password string) {
	return createUser(t, createUserArg())
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
	require.Equal(t, result.User.Username, payload.Username)

	return user, result.Token, password
}

func createUser(t *testing.T, arg port.RegisterUserParams) (user domain.User, token, password string) {
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
	require.Equal(t, user.Username, payload.Username)

	return user, result.Token, arg.User.Password
}

func createUserArg() port.RegisterUserParams {
	return port.RegisterUserParams{
		User: domain.RandomUser(),
	}
}
