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
	IDs       []domain.ID
	Usernames []string
	Emails    []string
}

type FilterUserFollowPayload struct {
	FollowerIDs []domain.ID
	FolloweeIDs []domain.ID
}

type UserRepository interface {
	CreateUser(context.Context, CreateUserPayload) (domain.User, error)
	UpdateUser(context.Context, UpdateUserPayload) (domain.User, error)
	FilterUser(context.Context, FilterUserPayload) ([]domain.User, error)
	FindOne(context.Context, FilterUserPayload) (domain.User, error)
	FilterFollow(context.Context, FilterUserFollowPayload) ([]domain.UserFollow, error)
}
