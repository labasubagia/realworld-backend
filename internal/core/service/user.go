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

func (s *userService) Register(ctx context.Context, req port.RegisterUserParams) (result port.RegisterUserResult, err error) {
	reqUser, err := domain.NewUser(req.User)
	if err != nil {
		return port.RegisterUserResult{}, exception.Into(err)
	}
	result.User, err = s.property.repo.User().CreateUser(ctx, port.CreateUserPayload{User: reqUser})
	if err != nil {
		return port.RegisterUserResult{}, exception.Into(err)
	}

	result.Token, _, err = s.property.tokenMaker.CreateToken(result.User.ID, 2*time.Hour)
	if err != nil {
		return port.RegisterUserResult{}, exception.Into(err)
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, req port.LoginUserParams) (result port.LoginUserResult, err error) {
	existing, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{Emails: []string{req.User.Email}})
	if err != nil {
		return port.LoginUserResult{}, exception.Into(err)
	}
	if len(existing) < 1 {
		return port.LoginUserResult{}, exception.Validation().AddError("exception", "email or password invalid")
	}

	result.User = existing[0]
	if err := util.CheckPassword(req.User.Password, result.User.Password); err != nil {
		return port.LoginUserResult{}, exception.Into(err)
	}

	result.Token, _, err = s.property.tokenMaker.CreateToken(result.User.ID, 2*time.Hour)
	if err != nil {
		return port.LoginUserResult{}, exception.Into(err)
	}

	return result, nil
}

func (s *userService) Current(ctx context.Context, arg port.AuthParams) (result port.CurrentUserResult, err error) {
	if arg.Payload == nil {
		return port.CurrentUserResult{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
	}

	existing, err := s.property.repo.User().FilterUser(ctx, port.FilterUserPayload{IDs: []domain.ID{arg.Payload.UserID}})
	if err != nil {
		return port.CurrentUserResult{}, exception.Into(err)
	}
	if len(existing) < 1 {
		return port.CurrentUserResult{}, exception.New(exception.TypePermissionDenied, "no user found", nil)
	}

	result.User = existing[0]
	result.Token = arg.Token

	return result, nil
}

func (s *userService) Update(ctx context.Context, arg port.UpdateUserParams) (result port.UpdateUserResult, err error) {
	if arg.AuthArg.Payload == nil {
		return port.UpdateUserResult{}, exception.New(exception.TypePermissionDenied, "token payload not provided", nil)
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
			return port.UpdateUserResult{}, exception.Into(err)
		}
	}

	result.User, err = s.property.repo.User().UpdateUser(ctx, payload)
	if err != nil {
		return port.UpdateUserResult{}, exception.Into(err)
	}

	result.Token = arg.AuthArg.Token
	return result, nil
}

func (s *userService) Profile(ctx context.Context, arg port.ProfileParams) (result port.ProfileResult, err error) {
	result.User, err = s.property.repo.User().FindOne(ctx, port.FilterUserPayload{Usernames: []string{arg.Username}})
	if err != nil {
		return port.ProfileResult{}, exception.Into(err)
	}

	if arg.AuthArg.Payload == nil {
		return result, nil
	}

	follows, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
		FolloweeIDs: []domain.ID{result.User.ID},
	})
	if err != nil {
		return port.ProfileResult{}, exception.Into(err)
	}
	if len(follows) > 0 {
		result.IsFollow = true
	}

	return result, nil
}

func (s *userService) Follow(ctx context.Context, arg port.ProfileParams) (result port.ProfileResult, err error) {
	if arg.AuthArg.Payload == nil {
		return result, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	result.User, err = s.property.repo.User().FindOne(ctx, port.FilterUserPayload{Usernames: []string{arg.Username}})
	if err != nil {
		return port.ProfileResult{}, exception.Into(err)
	}

	follows, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
		FolloweeIDs: []domain.ID{result.User.ID},
	})
	if len(follows) > 0 {
		result.IsFollow = true
		return result, nil
	}

	_, err = s.property.repo.User().Follow(ctx, port.FollowPayload{Follow: domain.UserFollow{
		FollowerID: arg.AuthArg.Payload.UserID,
		FolloweeID: result.User.ID,
	}})
	if err != nil {
		return port.ProfileResult{}, exception.Into(err)
	}

	result.IsFollow = true
	return result, nil
}

func (s *userService) UnFollow(ctx context.Context, arg port.ProfileParams) (result port.ProfileResult, err error) {
	if arg.AuthArg.Payload == nil {
		return result, exception.New(exception.TypePermissionDenied, "authentication required", nil)
	}

	result.User, err = s.property.repo.User().FindOne(ctx, port.FilterUserPayload{Usernames: []string{arg.Username}})
	if err != nil {
		return port.ProfileResult{}, exception.Into(err)
	}

	follows, err := s.property.repo.User().FilterFollow(ctx, port.FilterUserFollowPayload{
		FollowerIDs: []domain.ID{arg.AuthArg.Payload.UserID},
		FolloweeIDs: []domain.ID{result.User.ID},
	})
	if len(follows) < 1 {
		result.IsFollow = false
		return result, nil
	}

	_, err = s.property.repo.User().UnFollow(ctx, port.UnFollowPayload{Follow: domain.UserFollow{
		FollowerID: arg.AuthArg.Payload.UserID,
		FolloweeID: result.User.ID,
	}})
	if err != nil {
		return port.ProfileResult{}, exception.Into(err)
	}

	result.IsFollow = false
	return result, nil
}
