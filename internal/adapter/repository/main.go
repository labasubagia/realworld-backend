package repository

import (
	"sort"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/mongo"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

const DefaultRepoKey = "default"

type FnNew func(util.Config, port.Logger) (port.Repository, error)

var FnNewMap = map[string]FnNew{
	DefaultRepoKey: sql.NewSQLRepository,
	"postgres":     sql.NewSQLRepository,
	"mongo":        mongo.NewMongoRepository,
}

func Keys() (keys []string) {
	for key := range FnNewMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func ListRepository(config util.Config, logger port.Logger) ([]port.Repository, error) {
	repos := []port.Repository{}
	for _, fn := range FnNewMap {
		repo, err := fn(config, logger)
		if err != nil {
			return []port.Repository{}, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func NewRepository(config util.Config, logger port.Logger) (port.Repository, error) {
	return FnNewMap[DefaultRepoKey](config, logger)
}
