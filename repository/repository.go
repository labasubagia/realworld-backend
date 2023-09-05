package repository

import (
	"context"
	"database/sql"

	"github.com/labasubagia/realworld-backend/port"
	"github.com/uptrace/bun"
)

type repo struct {
	db          bun.IDB
	userRepo    port.UserRepository
	articleRepo port.ArticleRepository
}

func NewRepository(db bun.IDB) port.Repository {
	return &repo{
		db:          db,
		userRepo:    NewUserRepository(db),
		articleRepo: NewArticleRepository(db),
	}
}

func (r *repo) Atomic(ctx context.Context, fn port.RepositoryAtomicCallback) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		return fn(NewRepository(tx))
	})
}

func (r *repo) User() port.UserRepository {
	return r.userRepo
}

func (r *repo) Article() port.ArticleRepository {
	return r.articleRepo
}
