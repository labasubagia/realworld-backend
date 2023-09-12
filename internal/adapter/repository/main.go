package repository

import (
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/mongo"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

func ListRepository(config util.Config) ([]port.Repository, error) {
	type NewRepoFn func(util.Config) (port.Repository, error)
	repoFns := []NewRepoFn{
		mongo.NewMongoRepository,
		sql.NewSQLRepository,
	}
	repos := make([]port.Repository, len(repoFns))
	for i, fn := range repoFns {
		repo, err := fn(config)
		if err != nil {
			return []port.Repository{}, err
		}
		repos[i] = repo
	}
	return repos, nil
}

func NewRepository(config util.Config) (port.Repository, error) {
	return sql.NewSQLRepository(config)
	// return mongo.NewMongoRepository(config)
}
