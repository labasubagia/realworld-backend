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

func (r *userRepo) CreateUser(ctx context.Context, req port.CreateUserParams) (domain.User, error) {
	user := req.User
	_, err := r.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
