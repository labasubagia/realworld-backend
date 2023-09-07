package repository

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labasubagia/realworld-backend/port"
)

type Err pgconn.PgError

var mapServiceError = map[string]error{
	"40001": port.ErrIsolation,
	"23503": port.ErrForeignKey,
	"23505": port.ErrUniqueKey,
}

func errCode(err error) string {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		return pgErr.Code
	}
	return ""
}

func AsServiceError(err error) error {
	if serviceErr, isConfigured := mapServiceError[errCode(err)]; isConfigured {
		return serviceErr
	}
	return err
}
