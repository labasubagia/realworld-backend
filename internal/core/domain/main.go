package domain

import "github.com/labasubagia/realworld-backend/internal/core/util"

type ID int64

func RandomID() ID {
	return ID(util.RandomInt(1, 100))
}
