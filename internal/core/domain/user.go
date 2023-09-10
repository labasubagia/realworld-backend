package domain

import (
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"github.com/uptrace/bun"
)

const UserDefaultImage string = "https://api.realworld.io/images/demo-avatar.png"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            ID        `bun:"id,pk,autoincrement"`
	Email         string    `bun:"email,notnull"`
	Username      string    `bun:"username,notnull"`
	Password      string    `bun:"password,notnull"`
	Image         string    `bun:"image"`
	Bio           string    `bun:"bio"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	IsFollowed    bool      `bun:"-"`
}

func NewUser(arg User) (User, error) {
	validator := exception.Validation()

	user := User{}
	if err := user.SetEmail(arg.Email); err != nil {
		validator.AddError("email", err.Error())
	}
	if err := user.SetUsername(arg.Username); err != nil {
		validator.AddError("username", err.Error())
	}
	image := UserDefaultImage
	if arg.Image != "" {
		image = arg.Image
	}
	if err := user.SetImageURL(image); err != nil {
		validator.AddError("image", err.Error())
	}
	if err := user.SetPassword(arg.Password); err != nil {
		validator.AddError("password", err.Error())
	}

	if validator.HasError() {
		return user, validator
	}

	user.Bio = arg.Bio
	return user, nil
}

func RandomUser() User {
	return User{
		Email:    util.RandomEmail(),
		Username: util.RandomUsername(),
		Password: util.RandomString(8),
	}
}

func (user *User) SetEmail(email string) error {
	if err := util.ValidateEmail(email); err != nil {
		return err
	}
	user.Email = email
	return nil
}

func (user *User) SetPassword(password string) error {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}

func (user *User) SetUsername(username string) error {
	if err := util.ValidateUsername(username); err != nil {
		return err
	}
	user.Username = username
	return nil
}

func (user *User) SetImageURL(url string) error {
	if err := util.ValidateURL(url); err != nil {
		return err
	}
	user.Image = url
	return nil
}

type UserFollow struct {
	bun.BaseModel `bun:"table:user_follows,alias:uf"`
	FollowerID    ID `bun:"follower_id"`
	FolloweeID    ID `bun:"followee_id"`
}

func NewUserFollow(arg UserFollow) (UserFollow, error) {
	if arg.FollowerID == arg.FolloweeID {
		return UserFollow{}, exception.Validation().AddError("exception", "cannot follow yourself")
	}
	follow := UserFollow{
		FollowerID: arg.FollowerID,
		FolloweeID: arg.FolloweeID,
	}
	return follow, nil
}
