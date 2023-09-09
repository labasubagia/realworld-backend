package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type CreateUserPayload struct {
	User domain.User
}

type UpdateUserPayload struct {
	User domain.User
}

type FilterUserPayload struct {
	IDs       []int64
	Usernames []string
	Emails    []string
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserPayload) (domain.User, error)
	UpdateUser(context.Context, UpdateUserPayload) (domain.User, error)
	FilterUser(context.Context, FilterUserPayload) ([]domain.User, error)
}
