DB_URL=postgresql://postgres:postgres@localhost:5432/realworld?sslmode=disable
DB_MIGRATION_PATH=internal/adapter/repository/sql/db/migration

migrate_up:
	migrate -path "$(DB_MIGRATION_PATH)" -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path "$(DB_MIGRATION_PATH)" -database "$(DB_URL)" -verbose down

migrate_drop:
	migrate -path "$(DB_MIGRATION_PATH)" -database "$(DB_URL)" -verbose drop

# make new_migration -name=add_new_table
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

test:
	go test -cover ./...