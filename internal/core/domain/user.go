package domain

import (
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

const UserDefaultImage string = "https://api.realworld.io/images/demo-avatar.png"

type User struct {
	ID         ID
	Email      string
	Username   string
	Password   string
	Image      string
	Bio        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsFollowed bool
	Token      string
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

func NewUser(arg User) (User, error) {
	validator := exception.Validation()
	now := time.Now()

	user := User{
		ID:        NewID(),
		Bio:       arg.Bio,
		CreatedAt: now,
		UpdatedAt: now,
	}

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

	return user, nil
}

func RandomUser() User {
	return User{
		Email:    util.RandomEmail(),
		Username: util.RandomUsername(),
		Password: util.RandomString(8),
	}
}

type UserFollow struct {
	FollowerID ID
	FolloweeID ID
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
