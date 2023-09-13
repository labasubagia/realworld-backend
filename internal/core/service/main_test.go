package service_test

import (
	"log"
	"os"
	"testing"

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
		log.Fatal("failed load config", err)
	}

	var code int

	if config.IsTestAllRepo() {
		// test every store repo against service
		// slower test
		// NOTE: add env TEST_REPO=all

		repos, err := repository.ListRepository(config)
		if err != nil {
			log.Fatal("failed load repositories", err)
		}
		for _, repo := range repos {
			testRepo = repo
			testService, err = service.NewService(config, testRepo)
			if err != nil {
				log.Fatal("failed to init service", err)
			}
			code = m.Run()
		}

	} else {
		// test main repo (currently used)
		// faster test

		testRepo, err = repository.NewRepository(config)
		if err != nil {
			log.Fatal("failed load main repository", err)
		}

		testService, err = service.NewService(config, testRepo)
		if err != nil {
			log.Fatal("failed to init service", err)
		}
		code = m.Run()
	}

	os.Exit(code)
}
