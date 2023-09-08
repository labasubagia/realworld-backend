package repository

import (
	repository_sql "github.com/labasubagia/realworld-backend/internal/adapter/repository/sql"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func NewRepository(config util.Config) (port.Repository, error) {
	return repository_sql.NewSQLRepository(config)
}
