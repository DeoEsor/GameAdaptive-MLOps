LOCAL_BIN:=$(CURDIR)/bin
DB_NAME:=ml_ops
DB_DSN:=postgres://postgres:postgres@db:5432/$(DB_NAME)

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/go-jet/jet/v2/cmd/jet@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up:
	@PATH=$(LOCAL_BIN):$(PATH) goose -dir migrations -allow-missing postgres "$(DB_DSN)" up

db-create:
	psql -U postgres -h db -c "CREATE DATABASE IF NOT EXISTS $(DB_NAME)"
db-down:
	psql -U postgres -h db -c "DROP DATABASE IF EXISTS $(DB_NAME) WITH (FORCE)"

db-reset:db-down db-create migrate-up

run: migrate-up
	go run ./cmd/.
	
	
	

