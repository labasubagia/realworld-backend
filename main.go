package main

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labasubagia/realworld-backend/api"
	"github.com/labasubagia/realworld-backend/repository"
	"github.com/labasubagia/realworld-backend/service"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func main() {
	config, err := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/realworld?sslmode=disable")
	if err != nil {
		panic(err)
	}
	sqlDB := stdlib.OpenDB(*config)
	db := bun.NewDB(sqlDB, pgdialect.New())

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	server := api.NewServer(svc)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
