package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
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

type UserService interface {
	Register(context.Context, RegisterUserParams) (RegisterUserResult, error)
	Login(context.Context, LoginUserParams) (LoginUserResult, error)
}
