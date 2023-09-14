package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/labasubagia/realworld-backend/internal/adapter/logger"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/service"
	"github.com/labasubagia/realworld-backend/internal/core/util"
)

var testRepo port.Repository
var testService port.Service

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../.env.test")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err)
		os.Exit(1)
	}
	logger := logger.NewLogger(config)

	var code int

	if config.IsTestAllRepo() {
		// test every store repo against service
		// slower test
		// NOTE: add env TEST_REPO=all

		repos, err := repository.ListRepository(config, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load repository")
		}
		for _, repo := range repos {
			testRepo = repo
			testService, err = service.NewService(config, testRepo, logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to load service")
			}
			code = m.Run()
		}

	} else {
		// test main repo (currently used)
		// faster test

		testRepo, err = repository.NewRepository(config, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load repository")
		}

		testService, err = service.NewService(config, testRepo, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load service")
		}
		code = m.Run()
	}

	os.Exit(code)
}
