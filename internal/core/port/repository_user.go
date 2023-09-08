package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type CreateUserParams struct {
	User domain.User
}

type FilterUserParams struct {
	Usernames []string
	Emails    []string
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserParams) (domain.User, error)
	FilterUser(context.Context, FilterUserParams) ([]domain.User, error)
}
