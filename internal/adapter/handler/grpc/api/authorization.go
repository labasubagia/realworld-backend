package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader    = "authorization"
	authorizationTypeToken = "token"
	authorizationArgKey    = "authorization_arg"
)

func (s *Server) authorizeUser(ctx context.Context) (port.AuthParams, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return port.AuthParams{}, exception.New(exception.TypePermissionDenied, "missing metadata", nil)
	}
	values := metaData.Get(authorizationHeader)
	if len(values) == 0 {
		return port.AuthParams{}, exception.New(exception.TypePermissionDenied, "missing authorization header", nil)
	}

	fields := strings.Fields(values[0])
	if len(fields) < 2 {
		msg := "invalid authorization format"
		err := exception.New(exception.TypePermissionDenied, msg, nil)
		return port.AuthParams{}, err
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeToken {
		msg := fmt.Sprintf("authorization type %s not supported", authorizationType)
		err := exception.New(exception.TypePermissionDenied, msg, nil)
		return port.AuthParams{}, err
	}

	token := fields[1]
	payload, err := s.service.TokenMaker().VerifyToken(token)
	if err != nil {
		return port.AuthParams{}, err
	}
	return port.AuthParams{Token: token, Payload: payload}, nil
}
