package sql

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
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

func (r *userRepo) CreateUser(ctx context.Context, req port.CreateUserPayload) (domain.User, error) {
	user := req.User
	_, err := r.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return domain.User{}, intoException(err)
	}
	return user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, req port.UpdateUserPayload) (domain.User, error) {

	// find current
	current, err := r.FindOne(ctx, port.FilterUserPayload{IDs: []domain.ID{req.User.ID}})
	if err != nil {
		return domain.User{}, intoException(err)
	}

	// omit same
	if current.Email == req.User.Email {
		req.User.Email = ""
	}
	if current.Username == req.User.Username {
		req.User.Username = ""
	}

	// update
	_, err = r.db.NewUpdate().Model(&req.User).OmitZero().Where("id = ?", req.User.ID).Exec(ctx)
	if err != nil {
		return domain.User{}, intoException(err)
	}

	// find updated
	updated, err := r.FindOne(ctx, port.FilterUserPayload{IDs: []domain.ID{req.User.ID}})
	if err != nil {
		return domain.User{}, intoException(err)
	}

	return updated, nil
}

func (r *userRepo) FilterUser(ctx context.Context, filter port.FilterUserPayload) ([]domain.User, error) {
	users := []domain.User{}
	query := r.db.NewSelect().Model(&users)
	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.IDs))
	}
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

func (r *userRepo) FindOne(ctx context.Context, filter port.FilterUserPayload) (domain.User, error) {
	users, err := r.FilterUser(ctx, filter)
	if err != nil {
		return domain.User{}, intoException(err)
	}
	if len(users) < 1 {
		return domain.User{}, exception.New(exception.TypeNotFound, "not found", nil)
	}
	return users[0], nil
}
