package model

import (
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            domain.ID `bun:"id,pk,autoincrement"`
	Email         string    `bun:"email,notnull"`
	Username      string    `bun:"username,notnull"`
	Password      string    `bun:"password,notnull"`
	Image         string    `bun:"image"`
	Bio           string    `bun:"bio"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (data User) ToDomain() domain.User {
	return domain.User{
		ID:        data.ID,
		Email:     data.Email,
		Username:  data.Username,
		Password:  data.Password,
		Image:     data.Image,
		Bio:       data.Bio,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func AsUser(arg domain.User) User {
	return User{
		ID:        arg.ID,
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  arg.Password,
		Image:     arg.Image,
		Bio:       arg.Bio,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
	}
}

type UserFollow struct {
	bun.BaseModel `bun:"table:user_follows,alias:uf"`
	FollowerID    domain.ID `bun:"follower_id"`
	FolloweeID    domain.ID `bun:"followee_id"`
}

func (data UserFollow) ToDomain() domain.UserFollow {
	return domain.UserFollow{
		FollowerID: data.FollowerID,
		FolloweeID: data.FolloweeID,
	}
}

func AsUserFollow(arg domain.UserFollow) UserFollow {
	return UserFollow{
		FollowerID: arg.FollowerID,
		FolloweeID: arg.FolloweeID,
	}
}
