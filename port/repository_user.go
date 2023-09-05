package port

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
)

type CreateUserParams struct {
	User domain.User
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserParams) (domain.User, error)
}
