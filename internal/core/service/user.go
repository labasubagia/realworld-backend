package service

import (
	"context"
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

type userService struct {
	property serviceProperty
}

func NewUserService(property serviceProperty) port.UserService {
	return &userService{
		property: property,
	}
}

func (s *userService) Register(ctx context.Context, req port.RegisterParams) (user domain.User, err error) {
	reqUser, err := domain.NewUser(req.User)
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	user, err = s.property.repo.User().CreateUser(ctx, port.CreateUserPayload{User: reqUser})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	user.Token, _, err = s.property.tokenMaker.CreateToken(user.ID, 2*time.Hour)
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, req port.LoginParams) (user domain.User, err error) {
	existing, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{Emails: []string{req.User.Email}})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	if len(existing) < 1 {
		return domain.User{}, exception.Validation().AddError("exception", "email or password invalid")
	}

	user = existing[0]
	if err := util.CheckPassword(req.User.Password, user.Password); err != nil {
		return domain.User{}, exception.Into(err)
	}

	user.Token, _, err = s.property.tokenMaker.CreateToken(user.ID, 2*time.Hour)
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	return user, nil
}

func (s *userService) Current(ctx context.Context, arg port.AuthParams) (user domain.User, err error) {
	if arg.Payload == nil {
		return domain.User{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	existing, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{IDs: []domain.ID{arg.Payload.UserID}})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	if len(existing) == 0 {
		return domain.User{}, exception.New(exception.TypePermissionDenied, "no user found", nil)
	}

	user = existing[0]
	user.Token = arg.Token

	return user, nil
}

func (s *userService) Update(ctx context.Context, arg port.UpdateUserParams) (user domain.User, err error) {
	if arg.AuthArg.Payload == nil {
		return domain.User{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	payload := port.UpdateUserPayload{
		User: domain.User{
			ID:        arg.User.ID,
			Username:  arg.User.Username,
			Email:     arg.User.Email,
			Password:  arg.User.Password,
			Image:     arg.User.Image,
			Bio:       arg.User.Bio,
			UpdatedAt: time.Now(),
		},
	}
	if arg.User.Password != "" {
		if err := payload.User.SetPassword(arg.User.Password); err != nil {
			return domain.User{}, exception.Into(err)
		}
	}

	user, err = s.property.repo.User().UpdateUser(ctx, payload)
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	user.Token = arg.AuthArg.Token
	return user, nil
}

func (s *userService) Profile(ctx context.Context, arg port.ProfileParams) (user domain.User, err error) {
	user, err = s.property.repo.User().FindOne(ctx, port.FilterUserPayload{Usernames: []string{arg.Username}})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	if arg.AuthArg.Payload == nil {
		return user, nil
	}

	follows, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
		FolloweeIDs: []domain.ID{user.ID},
	})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	if len(follows) > 0 {
		user.IsFollowed = true
	}

	return user, nil
}

func (s *userService) Follow(ctx context.Context, arg port.ProfileParams) (user domain.User, err error) {
	if arg.AuthArg.Payload == nil {
		return user, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	user, err = s.property.repo.User().FindOne(ctx, port.FilterUserPayload{Usernames: []string{arg.Username}})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	follows, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
		FolloweeIDs: []domain.ID{user.ID},
	})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	if len(follows) > 0 {
		user.IsFollowed = true
		return user, nil
	}

	newFollow, err := domain.NewUserFollow(domain.UserFollow{
		FollowerID: arg.AuthArg.Payload.UserID,
		FolloweeID: user.ID,
	})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	_, err = s.property.repo.User().Follow(ctx, port.FollowPayload{Follow: newFollow})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	user.IsFollowed = true
	return user, nil
}

func (s *userService) UnFollow(ctx context.Context, arg port.ProfileParams) (user domain.User, err error) {
	if arg.AuthArg.Payload == nil {
		return user, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	user, err = s.property.repo.User().FindOne(ctx, port.FilterUserPayload{Usernames: []string{arg.Username}})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	follows, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
		FolloweeIDs: []domain.ID{user.ID},
	})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}
	if len(follows) == 0 {
		user.IsFollowed = false
		return user, nil
	}

	_, err = s.property.repo.User().UnFollow(ctx, port.UnFollowPayload{Follow: domain.UserFollow{
		FollowerID: arg.AuthArg.Payload.UserID,
		FolloweeID: user.ID,
	}})
	if err != nil {
		return domain.User{}, exception.Into(err)
	}

	user.IsFollowed = false
	return user, nil
}
