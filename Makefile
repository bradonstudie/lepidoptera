.PHONY: dev build generate migrate db

dev:
	air

build:
	templ generate && go build -o bin/server ./cmd/server

generate:
	sqlc generate && templ generate

migrate:
	goose -dir internal/db/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/lepidoptera?sslmode=disable" up

migrate-down:
	goose -dir internal/db/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/lepidoptera?sslmode=disable" down

db:
	docker-compose up -d db

db-stop:
	docker-compose down
