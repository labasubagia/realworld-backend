package api

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/grpc/pb"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.UserResponse, error) {
	user, err := server.service.User().Register(ctx, port.RegisterParams{
		User: domain.User{
			Email:    req.GetUser().GetEmail(),
			Username: req.GetUser().GetUsername(),
			Password: req.GetUser().GetPassword(),
		},
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.UserResponse{
		User: serializeUser(user),
	}
	return res, nil
}

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error) {
	user, err := server.service.User().Login(ctx, port.LoginParams{
		User: domain.User{
			Email:    req.GetUser().GetEmail(),
			Password: req.GetUser().GetPassword(),
		},
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.UserResponse{
		User: serializeUser(user),
	}
	return res, nil
}

func (server *Server) CurrentUser(ctx context.Context, _ *emptypb.Empty) (*pb.UserResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	user, err := server.service.User().Current(ctx, auth)
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.UserResponse{
		User: serializeUser(user),
	}
	return res, nil
}

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	user, err := server.service.User().Update(ctx, port.UpdateUserParams{
		AuthArg: auth,
		User: domain.User{
			ID:       auth.Payload.UserID,
			Email:    req.GetUser().GetEmail(),
			Username: req.GetUser().GetUsername(),
			Password: req.GetUser().GetPassword(),
			Image:    req.GetUser().GetImage(),
			Bio:      req.GetUser().GetBio(),
		},
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.UserResponse{
		User: serializeUser(user),
	}
	return res, nil
}

func (server *Server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileResponse, error) {
	auth, _ := server.authorizeUser(ctx)
	user, err := server.service.User().Profile(ctx, port.ProfileParams{
		Username: req.GetUsername(),
		AuthArg:  auth,
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ProfileResponse{
		Profile: serializeProfile(user),
	}
	return res, nil
}

func (server *Server) FollowUser(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	user, err := server.service.User().Follow(ctx, port.ProfileParams{
		Username: req.GetUsername(),
		AuthArg:  auth,
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ProfileResponse{
		Profile: serializeProfile(user),
	}
	return res, nil
}

func (server *Server) UnFollowUser(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	user, err := server.service.User().UnFollow(ctx, port.ProfileParams{
		Username: req.GetUsername(),
		AuthArg:  auth,
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ProfileResponse{
		Profile: serializeProfile(user),
	}
	return res, nil
}
