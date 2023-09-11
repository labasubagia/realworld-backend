package sql

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql/model"
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

func (r *userRepo) CreateUser(ctx context.Context, arg domain.User) (domain.User, error) {
	user := model.AsUser(arg)
	_, err := r.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return domain.User{}, intoException(err)
	}
	return user.ToDomain(), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, arg domain.User) (domain.User, error) {

	// find current
	current, err := r.FindOne(ctx, port.FilterUserPayload{IDs: []domain.ID{arg.ID}})
	if err != nil {
		return domain.User{}, intoException(err)
	}

	// omit same
	if current.Email == arg.Email {
		arg.Email = ""
	}
	if current.Username == arg.Username {
		arg.Username = ""
	}

	// update
	req := model.AsUser(arg)
	_, err = r.db.NewUpdate().Model(&req).OmitZero().Where("id = ?", req.ID).Exec(ctx)
	if err != nil {
		return domain.User{}, intoException(err)
	}

	// find updated
	updated, err := r.FindOne(ctx, port.FilterUserPayload{IDs: []domain.ID{req.ID}})
	if err != nil {
		return domain.User{}, intoException(err)
	}

	return updated, nil
}

func (r *userRepo) FilterUser(ctx context.Context, filter port.FilterUserPayload) ([]domain.User, error) {
	users := []model.User{}
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
	result := []domain.User{}
	for _, user := range users {
		result = append(result, user.ToDomain())
	}
	return result, nil
}

func (r *userRepo) FindOne(ctx context.Context, filter port.FilterUserPayload) (domain.User, error) {
	users, err := r.FilterUser(ctx, filter)
	if err != nil {
		return domain.User{}, intoException(err)
	}
	if len(users) == 0 {
		return domain.User{}, exception.New(exception.TypeNotFound, "user not found", nil)
	}
	return users[0], nil
}

func (r *userRepo) FilterFollow(ctx context.Context, filter port.FilterUserFollowPayload) ([]domain.UserFollow, error) {
	follows := []model.UserFollow{}
	query := r.db.NewSelect().Model(&follows)
	if len(filter.FollowerIDs) > 0 {
		query = query.Where("follower_id IN (?)", bun.In(filter.FollowerIDs))
	}
	if len(filter.FolloweeIDs) > 0 {
		query = query.Where("followee_id IN (?)", bun.In(filter.FolloweeIDs))
	}
	err := query.Scan(ctx)
	if err != nil {
		return []domain.UserFollow{}, intoException(err)
	}
	result := []domain.UserFollow{}
	for _, follow := range follows {
		result = append(result, follow.ToDomain())
	}
	return result, nil
}

func (r *userRepo) Follow(ctx context.Context, arg domain.UserFollow) (domain.UserFollow, error) {
	req := model.AsUserFollow(arg)
	_, err := r.db.NewInsert().Model(&req).Exec(ctx)
	if err != nil {
		return domain.UserFollow{}, exception.Into(err)
	}
	return req.ToDomain(), nil
}

func (r *userRepo) UnFollow(ctx context.Context, arg domain.UserFollow) (domain.UserFollow, error) {
	req := model.AsUserFollow(arg)
	_, err := r.db.
		NewDelete().
		Model(&req).
		Where("follower_id = ?", req.FollowerID).
		Where("followee_id = ?", req.FolloweeID).
		Exec(ctx)
	if err != nil {
		return domain.UserFollow{}, exception.Into(err)
	}
	return req.ToDomain(), nil
}
