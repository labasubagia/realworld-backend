# database

POSTGRES_URL=postgresql://postgres:postgres@0.0.0.0:5432/realworld?sslmode=disable
POSTGRES_MIGRATION_PATH=internal/adapter/repository/sql/db/migration

postgres_migrate_up:
	migrate -path "$(POSTGRES_MIGRATION_PATH)" -database "$(POSTGRES_URL)" -verbose up

postgres_migrate_down:
	migrate -path "$(POSTGRES_MIGRATION_PATH)" -database "$(POSTGRES_URL)" -verbose down

postgres_migrate_drop:
	migrate -path "$(POSTGRES_MIGRATION_PATH)" -database "$(POSTGRES_URL)" -verbose drop

# testing

test:
	go test -cover ./...

test_all:
	export TEST_REPO=all
	make test

test_cover:
	go test -coverprofile=coverage.profile -cover ./...
	go tool cover -html coverage.profile -o coverage.html

e2e:
	APIURL=http://0.0.0.0:5000 ./tests/run-api-tests.sh