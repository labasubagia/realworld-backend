package api

import (
	"github.com/labasubagia/realworld-backend/internal/adapter/handler/grpc/pb"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func serializeUser(arg domain.User) *pb.User {
	return &pb.User{
		Email:    arg.Email,
		Username: arg.Username,
		Bio:      arg.Bio,
		Image:    arg.Image,
		Token:    arg.Token,
	}
}

func serializeProfile(arg domain.User) *pb.Profile {
	return &pb.Profile{
		Username:  arg.Username,
		Image:     arg.Image,
		Bio:       arg.Bio,
		Following: arg.IsFollowed,
	}
}

func serializeArticle(arg domain.Article) *pb.Article {
	tags := []string{}
	if len(arg.TagNames) > 0 {
		tags = arg.TagNames
	}
	return &pb.Article{
		Slug:          arg.Slug,
		Title:         arg.Title,
		Description:   arg.Description,
		Body:          arg.Body,
		TagList:       tags,
		Favorited:     arg.IsFavorite,
		FavoriteCount: int64(arg.FavoriteCount),
		Author:        serializeProfile(arg.Author),
		CreatedAt:     timestamppb.New(arg.CreatedAt),
		UpdatedAt:     timestamppb.New(arg.UpdatedAt),
	}
}

func serializeComment(arg domain.Comment) *pb.Comment {
	return &pb.Comment{
		Id:        arg.ID.String(),
		Body:      arg.Body,
		Author:    serializeProfile(arg.Author),
		CreatedAt: timestamppb.New(arg.CreatedAt),
		UpdatedAt: timestamppb.New(arg.UpdatedAt),
	}
}
