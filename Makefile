.DEFAULT_GOAL := run
.PHONY: run 
	migrate-up 
	migrate-down 
	migrate-create 
	test 
	migrate-schema
	migrate-dev 
	migrate-prod 
	swagger 
	docs 
	load-test

GOOSE_DRIVER := postgres
GOOSE_DBSTRING := 'user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=${DB_SSLMODE}'

run:
	go run cmd/api/main.go

migrate-up:
	goose -dir ./migrations/schema ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up

migrate-down:
	goose -dir ./migrations/schema ${GOOSE_DRIVER} ${GOOSE_DBSTRING} down

migrate-create:
	@if [ -z "$(dir)" ]; then \
		dir="schema"; \
	fi
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name [dir=schema|development|production]"; \
		exit 1; \
	fi
	@mkdir -p ./migrations/$(dir)
	goose -dir ./migrations/$(dir) create $(name) sql

migrate-status:
	@echo "Schema migrations:"
	goose -dir ./migrations/schema ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" status
	@echo "\nDevelopment migrations:"
	goose -dir ./migrations/development ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" status
	@echo "\nProduction migrations:"
	goose -dir ./migrations/production ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" status

migrate-schema:
	goose -dir ./migrations/schema ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" up

migrate-dev:
	goose -dir ./migrations/development ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" up

migrate-prod:
	goose -dir ./migrations/production ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" up

test:
	go test -v ./...

dev:
	ENVIRONMENT=development docker-compose up --build

db:
	docker-compose up -d db

validate: test
	go mod tidy
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
	gosec ./...

install-swag-deps:
	go get -u github.com/swaggo/swag/cmd/swag

swagger:
	swag fmt
	swag init -g cmd/api/main.go -o ./docs

docs: swagger
	redocly lint docs/swagger.json
	redocly build-docs docs/swagger.json -o docs/index.html