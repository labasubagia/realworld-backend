package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type CreateUserParams struct {
	User domain.User
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserParams) (domain.User, error)
}