package repository

import (
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/mongo"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

const DefaultRepoKey = "default"

type FnNewRepo func(util.Config) (port.Repository, error)

var RepoFnMap = map[string]FnNewRepo{
	DefaultRepoKey: sql.NewSQLRepository,
	"postgres":     sql.NewSQLRepository,
	"mongo":        mongo.NewMongoRepository,
}

func ListRepository(config util.Config) ([]port.Repository, error) {
	repos := []port.Repository{}
	for _, fn := range RepoFnMap {
		repo, err := fn(config)
		if err != nil {
			return []port.Repository{}, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func NewRepository(config util.Config) (port.Repository, error) {
	return RepoFnMap[DefaultRepoKey](config)
}
