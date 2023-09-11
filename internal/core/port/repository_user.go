package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

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
	CreateUser(context.Context, domain.User) (domain.User, error)
	UpdateUser(context.Context, domain.User) (domain.User, error)
	FilterUser(context.Context, FilterUserPayload) ([]domain.User, error)
	FindOne(context.Context, FilterUserPayload) (domain.User, error)

	FilterFollow(context.Context, FilterUserFollowPayload) ([]domain.UserFollow, error)
	Follow(context.Context, domain.UserFollow) (domain.UserFollow, error)
	UnFollow(context.Context, domain.UserFollow) (domain.UserFollow, error)
}
