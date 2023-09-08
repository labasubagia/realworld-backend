package service_test

import (
	"log"
	"os"
	"testing"

	repository "github.com/labasubagia/realworld-backend/internal/adapter/repository/sql"
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
	testRepo, err = repository.NewSQLRepository(config)
	if err != nil {
		log.Fatal("failed to init repository", err)
	}
	testService = service.NewService(testRepo)
	code := m.Run()
	os.Exit(code)
}
