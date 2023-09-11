package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/util/token"
)

type AuthParams struct {
	Token   string
	Payload *token.Payload
}

type RegisterParams struct {
	User domain.User
}

type LoginParams struct {
	User domain.User
}
type UpdateUserParams struct {
	AuthArg AuthParams
	User    domain.User
}

type ProfileParams struct {
	AuthArg  AuthParams
	Username string
}

type UserService interface {
	Register(context.Context, RegisterParams) (domain.User, error)
	Login(context.Context, LoginParams) (domain.User, error)
	Update(context.Context, UpdateUserParams) (domain.User, error)
	Current(context.Context, AuthParams) (domain.User, error)

	Profile(context.Context, ProfileParams) (domain.User, error)
	Follow(context.Context, ProfileParams) (domain.User, error)
	UnFollow(context.Context, ProfileParams) (domain.User, error)
}
