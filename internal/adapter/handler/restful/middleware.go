package restful

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeToken = "token"
	authorizationArgKey    = "authorization_arg"
)

func (s *Server) AuthMiddleware(autoDenied bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authArg, err := s.parseToken(c)
		if err != nil {
			if autoDenied {
				errorHandler(c, err)
				return
			}
		}
		c.Set(authorizationArgKey, authArg)
		c.Next()
	}
}

func (s *Server) parseToken(c *gin.Context) (port.AuthParams, error) {
	authorizationHeader := c.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		msg := "authorization header not provided"
		err := exception.New(exception.TypePermissionDenied, msg, nil)
		return port.AuthParams{}, err
	}

	fields := strings.Fields(authorizationHeader)
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

func getAuthArg(c *gin.Context) (port.AuthParams, error) {
	arg, ok := c.Get(authorizationArgKey)
	if !ok {
		return port.AuthParams{}, exception.New(exception.TypePermissionDenied, "no authorization arguments provided", nil)
	}
	authArg, ok := arg.(port.AuthParams)
	if !ok {
		return port.AuthParams{}, exception.New(exception.TypePermissionDenied, "invalid authorization arguments", nil)
	}
	return authArg, nil
}
