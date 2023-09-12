package model

import (
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
)

type User struct {
	ID        domain.ID `bson:"id"`
	Email     string    `bson:"email"`
	Username  string    `bson:"username"`
	Password  string    `bson:"password"`
	Image     string    `bson:"image"`
	Bio       string    `bson:"bio"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
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
	FollowerID domain.ID `bson:"follower_id"`
	FolloweeID domain.ID `bson:"followee_id"`
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
