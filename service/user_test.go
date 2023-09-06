package service_test

import (
	"context"
	"testing"

	"github.com/labasubagia/realworld-backend/domain"
	"github.com/labasubagia/realworld-backend/port"
	"github.com/labasubagia/realworld-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func createRandomUser(t *testing.T) (user domain.User, password string) {
	return createUser(t, createUserArg())
}

func createUser(t *testing.T, arg port.CreateUserTxParams) (user domain.User, password string) {
	result, err := testService.User().Create(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, result)
	user = result.User
	require.Equal(t, arg.User.Email, user.Email)
	require.Equal(t, arg.User.Username, user.Username)
	require.Equal(t, arg.User.Image, user.Image)
	require.NotEqual(t, arg.User.Password, user.Password)
	require.Nil(t, util.CheckPassword(arg.User.Password, user.Password))
	return user, arg.User.Password
}

func createUserArg() port.CreateUserTxParams {
	return port.CreateUserTxParams{
		User: domain.User{
			Email:    util.RandomEmail(),
			Username: util.RandomUsername(),
			Password: util.RandomString(8),
			Image:    util.RandomURL(),
			Bio:      util.RandomString(10),
		},
	}
}
