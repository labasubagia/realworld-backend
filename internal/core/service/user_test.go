package service_test

import (
	"context"
	"testing"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestCreateUserWithImage(t *testing.T) {
	arg := createUserArg()
	arg.User.SetImageURL(util.RandomURL())
	createUser(t, arg)
}

func createRandomUser(t *testing.T) (user domain.User, password string) {
	return createUser(t, createUserArg())
}

func createUser(t *testing.T, arg port.CreateUserTxParams) (user domain.User, password string) {
	image := arg.User.Image
	if image == "" {
		image = domain.UserDefaultImage
	}

	result, err := testService.User().Create(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, result)
	user = result.User
	require.Equal(t, arg.User.Email, user.Email)
	require.Equal(t, arg.User.Username, user.Username)
	require.Equal(t, image, user.Image)
	require.NotEqual(t, arg.User.Password, user.Password)
	require.Nil(t, util.CheckPassword(arg.User.Password, user.Password))
	return user, arg.User.Password
}

func createUserArg() port.CreateUserTxParams {
	return port.CreateUserTxParams{
		User: domain.RandomUser(),
	}
}
