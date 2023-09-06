package domain

import (
	"fmt"
	"time"

	"github.com/labasubagia/realworld-backend/util"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bu:"table:users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Email         string    `bun:"email,notnull"`
	Username      string    `bun:"username,notnull"`
	Password      string    `bun:"password,notnull"`
	Image         string    `bun:"image"`
	Bio           string    `bun:"bio"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func NewUser(arg User) (User, error) {
	user := User{}
	if err := user.SetEmail(arg.Email); err != nil {
		return user, fmt.Errorf("email %w", err)
	}
	if err := user.SetUsername(arg.Username); err != nil {
		return user, fmt.Errorf("username %w", err)
	}
	if err := user.SetImageURL(arg.Image); err != nil {
		return user, fmt.Errorf("image %w", err)
	}
	if err := user.SetPassword(arg.Password); err != nil {
		return user, fmt.Errorf("password %w", err)
	}
	user.Bio = arg.Bio
	return user, nil
}

func RandomUser() User {
	return User{
		Email:    util.RandomEmail(),
		Username: util.RandomUsername(),
		Password: util.RandomString(8),
		Image:    util.RandomURL(),
		Bio:      util.RandomString(10),
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
