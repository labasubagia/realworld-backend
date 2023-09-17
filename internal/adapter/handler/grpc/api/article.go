package api

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/adapter/handler/grpc/pb"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) ListArticle(ctx context.Context, req *pb.FilterArticleRequest) (*pb.ArticlesResponse, error) {

	auth, _ := server.authorizeUser(ctx)

	offset := 0
	if req.Offset != nil {
		offset = int(req.GetOffset())
	}

	limit := DefaultPaginationSize
	if req.Limit != nil {
		limit = int(req.GetLimit())
	}

	arg := port.ListArticleParams{
		Tags:           []string{},
		AuthorNames:    []string{},
		FavoritedNames: []string{},
		AuthArg:        auth,
		Offset:         offset,
		Limit:          limit,
	}
	if req.GetTag() != "" {
		arg.Tags = append(arg.Tags, req.GetTag())
	}
	if req.GetAuthor() != "" {
		arg.AuthorNames = append(arg.AuthorNames, req.GetAuthor())
	}
	if req.GetFavorited() != "" {
		arg.FavoritedNames = append(arg.FavoritedNames, req.GetFavorited())
	}

	articles, err := server.service.Article().List(ctx, arg)
	if err != nil {
		return nil, handleError(err)
	}

	res := &pb.ArticlesResponse{
		Articles: []*pb.Article{},
		Count:    int64(len(articles)),
	}
	for _, article := range articles {
		res.Articles = append(res.Articles, serializeArticle(article))
	}

	return res, nil
}

func (server *Server) FeedArticle(ctx context.Context, req *pb.FilterArticleRequest) (*pb.ArticlesResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	offset := 0
	if req.Offset != nil {
		offset = int(req.GetOffset())
	}

	limit := DefaultPaginationSize
	if req.Limit != nil {
		limit = int(req.GetLimit())
	}

	arg := port.ListArticleParams{
		AuthArg: auth,
		Offset:  offset,
		Limit:   limit,
	}
	articles, err := server.service.Article().Feed(ctx, arg)
	if err != nil {
		return nil, handleError(err)
	}

	res := &pb.ArticlesResponse{
		Articles: []*pb.Article{},
		Count:    int64(len(articles)),
	}
	for _, article := range articles {
		res.Articles = append(res.Articles, serializeArticle(article))
	}

	return res, nil
}

func (server *Server) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.ArticleResponse, error) {

	authArg, _ := server.authorizeUser(ctx)

	article, err := server.service.Article().Get(ctx, port.GetArticleParams{
		AuthArg: authArg,
		Slug:    req.GetSlug(),
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ArticleResponse{
		Article: serializeArticle(article),
	}
	return res, nil
}

func (server *Server) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.ArticleResponse, error) {
	authArg, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	article, err := server.service.Article().Create(ctx, port.CreateArticleTxParams{
		AuthArg: authArg,
		Tags:    req.GetArticle().GetTagList(),
		Article: domain.Article{
			Title:       req.GetArticle().GetTitle(),
			Description: req.GetArticle().GetDescription(),
			Body:        req.GetArticle().GetBody(),
		},
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ArticleResponse{
		Article: serializeArticle(article),
	}
	return res, nil
}

func (server *Server) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.ArticleResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	article, err := server.service.Article().Update(ctx, port.UpdateArticleParams{
		AuthArg: auth,
		Slug:    req.GetSlug(),
		Article: domain.Article{
			Title:       req.GetArticle().GetTitle(),
			Description: req.GetArticle().GetDescription(),
			Body:        req.GetArticle().GetBody(),
		},
	})
	if err != nil {
		return nil, handleError(err)
	}

	res := &pb.ArticleResponse{
		Article: serializeArticle(article),
	}
	return res, nil
}

func (server *Server) DeleteArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.Response, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	err = server.service.Article().Delete(ctx, port.DeleteArticleParams{
		AuthArg: auth,
		Slug:    req.GetSlug(),
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.Response{Status: "OK"}
	return res, nil
}

func (server *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CommentResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	result, err := server.service.Article().AddComment(ctx, port.AddCommentParams{
		AuthArg: auth,
		Slug:    req.GetSlug(),
		Comment: domain.Comment{
			Body: req.GetComment().GetBody(),
		},
	})
	if err != nil {
		return nil, handleError(err)
	}

	res := &pb.CommentResponse{
		Comment: serializeComment(result),
	}
	return res, nil
}

func (server *Server) ListComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.CommentsResponse, error) {
	auth, _ := server.authorizeUser(ctx)
	comments, err := server.service.Article().ListComments(ctx, port.ListCommentParams{
		AuthArg: auth,
		Slug:    req.GetSlug(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	res := &pb.CommentsResponse{
		Comments: []*pb.Comment{},
	}
	for _, comment := range comments {
		res.Comments = append(res.Comments, serializeComment(comment))
	}
	return res, nil
}

func (server *Server) DeleteComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.Response, error) {
	commentID, err := domain.ParseID(req.GetCommentId())
	if err != nil {
		return nil, handleError(err)
	}
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	err = server.service.Article().DeleteComment(ctx, port.DeleteCommentParams{
		AuthArg:   auth,
		Slug:      req.Slug,
		CommentID: domain.ID(commentID),
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.Response{Status: "OK"}
	return res, nil
}

func (server *Server) FavoriteArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.ArticleResponse, error) {
	auth, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	article, err := server.service.Article().AddFavorite(ctx, port.AddFavoriteParams{
		AuthArg: auth,
		Slug:    req.GetSlug(),
		UserID:  auth.Payload.UserID,
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ArticleResponse{
		Article: serializeArticle(article),
	}
	return res, nil
}

func (server *Server) UnFavoriteArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.ArticleResponse, error) {
	authArg, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	article, err := server.service.Article().RemoveFavorite(ctx, port.RemoveFavoriteParams{
		AuthArg: authArg,
		Slug:    req.GetSlug(),
		UserID:  authArg.Payload.UserID,
	})
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ArticleResponse{
		Article: serializeArticle(article),
	}
	return res, nil
}

func (server *Server) ListTag(ctx context.Context, _ *emptypb.Empty) (*pb.ListTagResponse, error) {
	tags, err := server.service.Article().ListTags(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	res := &pb.ListTagResponse{
		Tags: tags,
	}
	return res, nil
}
