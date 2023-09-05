package port

import (
	"context"

	"github.com/labasubagia/realworld-backend/domain"
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
