package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labasubagia/realworld-backend/domain"
	"github.com/labasubagia/realworld-backend/port"
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

	// server := api.NewServer(svc)
	// if err := server.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// just try here
	ctx := context.Background()
	userArg := port.CreateUserTxParams{
		User: domain.User{
			Email:    "user@mail.com",
			Username: "user",
			Password: "12345678",
		},
	}
	userTx, err := svc.User().Create(ctx, userArg)
	if err != nil {
		log.Fatal("failed", err)
	}
	articleArg := port.CreateArticleTxParams{
		Article: domain.Article{
			AuthorID:    userTx.User.ID,
			Title:       "title1",
			Slug:        "title1",
			Description: "desc of title1",
			Body:        "body of title1",
		},
		Tags: []string{"first", "second"},
	}
	articleTx, err := svc.Article().Create(ctx, articleArg)
	if err != nil {
		log.Fatal("fa", err)
	}
	log.Println(articleTx.Article.Title, articleTx.Tags)
}
