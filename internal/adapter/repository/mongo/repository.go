package mongo

import (
	"context"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type mongoRepo struct {
	db          DB
	userRepo    port.UserRepository
	articleRepo port.ArticleRepository
}

func NewMongoRepository(config util.Config) (port.Repository, error) {
	db, err := NewDB(config)
	if err != nil {
		return nil, err
	}
	return create(db), nil
}

func create(db DB) port.Repository {
	return &mongoRepo{
		db:          db,
		userRepo:    NewUserRepository(db),
		articleRepo: NewArticleRepository(db),
	}
}

func (r *mongoRepo) Atomic(ctx context.Context, fn port.RepositoryAtomicCallback) error {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := r.db.Client().StartSession()
	if err != nil {
		return intoException(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (any, error) {
		if err := fn(create(r.db)); err != nil {
			return nil, intoException(err)
		}
		return nil, nil
	}, txnOptions)

	if err != nil {
		return intoException(err)
	}

	return nil
}

func (r *mongoRepo) User() port.UserRepository {
	return r.userRepo
}

func (r *mongoRepo) Article() port.ArticleRepository {
	return r.articleRepo
}
