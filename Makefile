DB_URL=postgresql://postgres:postgres@0.0.0.0:5432/realworld?sslmode=disable
DB_MIGRATION_PATH=internal/adapter/repository/sql/db/migration

migrate_up:
	migrate -path "$(DB_MIGRATION_PATH)" -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path "$(DB_MIGRATION_PATH)" -database "$(DB_URL)" -verbose down

migrate_drop:
	migrate -path "$(DB_MIGRATION_PATH)" -database "$(DB_URL)" -verbose drop

test:
	go test -cover ./...

e2e:
	APIURL=http://0.0.0.0:5000 ./tests/run-api-tests.sh