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
	result.User, err = s.property.repo.User().CreateUser(ctx, port.CreateUserParams{User: reqUser})
	if err != nil {
		return port.RegisterUserResult{}, exception.Into(err)
	}

	result.Token, _, err = s.property.tokenMaker.CreateToken(result.User.Username, 2*time.Hour)
	if err != nil {
		return port.RegisterUserResult{}, exception.Into(err)
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, req port.LoginUserParams) (result port.LoginUserResult, err error) {
	existing, err := s.property.repo.User().FilterUser(ctx, port.FilterUserParams{Emails: []string{req.User.Email}})
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

	result.Token, _, err = s.property.tokenMaker.CreateToken(result.User.Username, 2*time.Hour)
	if err != nil {
		return port.LoginUserResult{}, exception.Into(err)
	}

	return result, nil
}
