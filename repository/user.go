package repository

import (
	"context"

	"github.com/labasubagia/go-backend-realworld/domain"
	"github.com/labasubagia/go-backend-realworld/port"
	"github.com/uptrace/bun"
)

type userRepo struct {
	db bun.IDB
}

func NewUserRepository(db bun.IDB) port.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(context.Context, port.CreateUserParams) (domain.User, error) {
	return domain.User{}, nil
}
