package service_test

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labasubagia/realworld-backend/port"
	"github.com/labasubagia/realworld-backend/repository"
	"github.com/labasubagia/realworld-backend/service"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var testService port.Service

func TestMain(m *testing.M) {
	config, err := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/realworld?sslmode=disable")
	if err != nil {
		panic(err)
	}
	sqlDB := stdlib.OpenDB(*config)
	db := bun.NewDB(sqlDB, pgdialect.New())

	repo := repository.NewRepository(db)
	testService = service.NewService(repo)
	os.Exit(m.Run())
}
