package db

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labasubagia/realworld-backend/internal/adapter/repository/sql/db/migration"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type DB struct {
	config util.Config
	logger port.Logger
	db     *bun.DB
}

func New(config util.Config, logger port.Logger) (*DB, error) {
	var err error
	database := &DB{
		config: config,
		logger: logger,
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
	config, err := pgx.ParseConfig(db.config.PostgresSource)
	if err != nil {
		return nil, err
	}
	sqlDB := stdlib.OpenDB(*config)
	database := bun.NewDB(sqlDB, pgdialect.New())

	// log
	if !db.config.IsProduction() {
		// database.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		database.AddQueryHook(&LoggerHook{verbose: true, logger: db.logger})
	}
	return database, nil
}

func (db *DB) migrate() error {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("unable to get the current current file")
	}
	currentDir := filepath.Dir(currentFile)
	migrationURL := fmt.Sprintf("file://%s", path.Join(currentDir, "migration"))

	migration, err := migrate.New(migrationURL, db.config.PostgresSource)
	if err != nil {
		return fmt.Errorf("cannot create new migration instance: %s", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate: %s", err)
	}
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
		return nil
	}
	return nil
}
