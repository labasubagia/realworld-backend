package service_test

import (
	"context"
	"testing"

	"github.com/labasubagia/realworld-backend/domain"
	"github.com/labasubagia/realworld-backend/port"
	"github.com/labasubagia/realworld-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) (user domain.User, password string) {
	password = util.RandomString(8)
	arg := port.CreateUserTxParams{
		User: domain.User{
			Email:    util.RandomEmail(),
			Username: util.RandomUsername(),
			Password: password,
			Image:    util.RandomURL(),
			Bio:      util.RandomString(10),
		},
	}
	result, err := testService.User().Create(context.Background(), arg)
	require.Nil(t, err)
	require.NotEmpty(t, result)
	user = result.User
	require.Equal(t, arg.User.Email, user.Email)
	require.Equal(t, arg.User.Username, user.Username)
	require.Equal(t, arg.User.Image, user.Image)
	require.NotEqual(t, password, user.Password)
	require.Nil(t, util.CheckPassword(password, user.Password))
	return user, password
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
