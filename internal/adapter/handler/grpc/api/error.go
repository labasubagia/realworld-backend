package api

import (
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}
	fail, ok := err.(*exception.Exception)
	if !ok {
		return status.Error(codes.Internal, err.Error())
	}
	if !fail.HasError() {
		fail.AddError("exception", fail.Message)
	}
	var code codes.Code
	switch fail.Type {
	case exception.TypeNotFound:
		code = codes.NotFound
	case exception.TypeTokenExpired, exception.TypeTokenInvalid, exception.TypePermissionDenied:
		code = codes.Unauthenticated
	case exception.TypeValidation:
		code = codes.InvalidArgument
	default:
		code = codes.Internal
	}
	return status.Error(code, err.Error())
}
