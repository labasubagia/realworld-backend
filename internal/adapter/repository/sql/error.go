package sql

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

var mapException = map[string]string{
	"40001": exception.TypeValidation,
	"23503": exception.TypeValidation,
	"23505": exception.TypeValidation,
}

func postgresErrCode(err error) string {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		return pgErr.Code
	}
	return ""
}

func intoException(err error) *exception.Exception {
	kind, ok := mapException[postgresErrCode(err)]
	if ok {
		return exception.New(kind, err.Error(), err)
	}
	return exception.Into(err)
}
