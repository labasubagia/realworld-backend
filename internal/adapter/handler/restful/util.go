package restful

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

const formatTime string = "2006-01-02T15:04:05.999Z"

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

func hasToken(c *gin.Context) bool {
	authorizationHeader := c.GetHeader(authorizationHeaderKey)
	return len(authorizationHeader) > 0
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

func getPagination(c *gin.Context) (offset, limit int) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}
	return offset, limit
}

func timeString(t time.Time) string {
	return t.UTC().Format(formatTime)
}
