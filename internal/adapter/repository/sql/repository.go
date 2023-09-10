package sql

import (
	"context"
	"database/sql"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql/db"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/uptrace/bun"
)

type sqlRepo struct {
	db          bun.IDB
	userRepo    port.UserRepository
	articleRepo port.ArticleRepository
}

func NewSQLRepository(config util.Config) (port.Repository, error) {
	db, err := db.New(config)
	if err != nil {
		return nil, err
	}
	return create(db.DB()), nil
}

func (r *sqlRepo) Atomic(ctx context.Context, fn port.RepositoryAtomicCallback) error {
	err := r.db.RunInTx(
		ctx,
		&sql.TxOptions{Isolation: sql.LevelSerializable},
		func(ctx context.Context, tx bun.Tx) error {
			return fn(create(tx))
		},
	)
	if err != nil {
		return intoException(err)
	}
	return nil
}

func create(db bun.IDB) port.Repository {
	return &sqlRepo{
		db:          db,
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
