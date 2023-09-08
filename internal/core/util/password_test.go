package util_test

import (
	"testing"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	err = util.CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := util.RandomString(6)
	err = util.CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	differentHashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, differentHashedPassword)
	require.NotEqual(t, hashedPassword, differentHashedPassword)
}
