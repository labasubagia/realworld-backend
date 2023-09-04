package main

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labasubagia/go-backend-realworld/repository"
	"github.com/labasubagia/go-backend-realworld/service"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func main() {
	config, err := pgx.ParseConfig("postgres://postgres:@localhost:5432/test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	sqlDB := stdlib.OpenDB(*config)
	db := bun.NewDB(sqlDB, pgdialect.New())

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	fmt.Println(svc)
}
