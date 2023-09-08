package sql

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
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
		return domain.User{}, intoException(err)
	}
	return user, nil
}

func (r *userRepo) FilterUser(ctx context.Context, filter port.FilterUserParams) ([]domain.User, error) {
	users := []domain.User{}
	query := r.db.NewSelect().Model(&users)
	if len(filter.Emails) > 0 {
		query = query.Where("email IN (?)", bun.In(filter.Emails))
	}
	if len(filter.Usernames) > 0 {
		query = query.Where("username IN (?)", bun.In(filter.Usernames))
	}
	err := query.Scan(ctx)
	if err != nil {
		return []domain.User{}, intoException(err)
	}
	return users, nil
}
