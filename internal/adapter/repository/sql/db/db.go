package db

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql/db/migration"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type DB struct {
	config util.Config
	db     *bun.DB
}

func New(config util.Config) (*DB, error) {
	var err error
	database := &DB{
		config: config,
	}
	if database.db, err = database.connect(); err != nil {
		return nil, err
	}
	if err := database.migrate(); err != nil {
		return nil, err
	}
	return database, nil
}

func (db *DB) DB() *bun.DB {
	return db.db
}

func (db *DB) connect() (*bun.DB, error) {
	config, err := pgx.ParseConfig(db.config.DBSource)
	if err != nil {
		return nil, err
	}
	sqlDB := stdlib.OpenDB(*config)
	database := bun.NewDB(sqlDB, pgdialect.New())
	database.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return database, nil
}

func (db *DB) migrate() error {
	migration, err := migrate.New(db.config.DBMigrationURL, db.config.DBSource)
	if err != nil {
		log.Printf("cannot create new migration instance: %s", err)
		return err
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("failed to run migrate: %s", err)
		return err
	}
	log.Println("migration ok")
	return nil
}

func (db *DB) migrateBun() error {
	ctx := context.Background()
	migrator := bunMigrate.NewMigrator(db.db, migration.Migrations)

	if err := migrator.Init(ctx); err != nil {
		return err
	}
	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx)

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		log.Printf("there are no new migrations to run (database is up to date)\n")
		return nil
	}
	log.Printf("migrated to %s\n", group)
	return nil
}
