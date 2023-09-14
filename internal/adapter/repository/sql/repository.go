package sql

import (
	"context"
	"database/sql"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql/db"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/uptrace/bun"
)

const TypePostgres = "postgres"

type sqlRepo struct {
	db          bun.IDB
	logger      port.Logger
	userRepo    port.UserRepository
	articleRepo port.ArticleRepository
}

func NewSQLRepository(config util.Config, logger port.Logger) (port.Repository, error) {
	db, err := db.New(config, logger)
	if err != nil {
		return nil, err
	}
	return create(db.DB(), logger), nil
}

func (r *sqlRepo) Atomic(ctx context.Context, fn port.RepositoryAtomicCallback) error {
	err := r.db.RunInTx(
		ctx,
		&sql.TxOptions{Isolation: sql.LevelSerializable},
		func(ctx context.Context, tx bun.Tx) error {
			return fn(create(tx, r.logger))
		},
	)
	if err != nil {
		return intoException(err)
	}
	return nil
}

func create(db bun.IDB, logger port.Logger) port.Repository {
	return &sqlRepo{
		db:          db,
		logger:      logger,
		userRepo:    NewUserRepository(db),
		articleRepo: NewArticleRepository(db),
	}
}

func (r *sqlRepo) User() port.UserRepository {
	return r.userRepo
}

func (r *sqlRepo) Article() port.ArticleRepository {
	return r.articleRepo
}
