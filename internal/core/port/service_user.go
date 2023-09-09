package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/util/token"
)

type RegisterUserParams struct {
	User domain.User
}

type RegisterUserResult struct {
	User  domain.User
	Token string
}

type LoginUserParams struct {
	User domain.User
}

type LoginUserResult struct {
	User  domain.User
	Token string
}

type CurrentUserResult struct {
	User  domain.User
	Token string
}

type AuthParams struct {
	Token   string
	Payload *token.Payload
}

type UpdateUserParams struct {
	AuthArg AuthParams
	User    domain.User
}

type UpdateUserResult struct {
	User  domain.User
	Token string
}

type UserService interface {
	Register(context.Context, RegisterUserParams) (RegisterUserResult, error)
	Login(context.Context, LoginUserParams) (LoginUserResult, error)
	Update(context.Context, UpdateUserParams) (UpdateUserResult, error)
	Current(context.Context, AuthParams) (CurrentUserResult, error)
}
