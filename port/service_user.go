package port

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
)

type CreateUserTxParams struct {
	User      domain.User
	AfterFunc func(domain.User) error
}

type CreateUserTxResult struct {
	User domain.User
}

type UserService interface {
	Create(context.Context, CreateUserTxParams) (CreateUserTxResult, error)
}
