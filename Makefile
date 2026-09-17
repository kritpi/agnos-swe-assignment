SWAG_VERSION ?= v1.16.6
MOCKERY_VERSION ?= v2.53.3

.PHONY: swagger mocks build run test tidy docker-up docker-down docker-logs

## swagger: regenerate the OpenAPI docs into docs/
swagger:
	go run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) init \
		-g cmd/server/main.go --parseInternal --parseDependency -o docs

## mocks: regenerate the Repository port mock into internal/core/port/mocks
mocks:
	go run github.com/vektra/mockery/v2@$(MOCKERY_VERSION) \
		--dir internal/core/port --name Repository \
		--output internal/core/port/mocks --outpkg mocks --filename repository.go \
		--with-expecter --disable-version-string --issue-845-fix

## build: compile the API binary
build:
	go build -o bin/api ./cmd/server

## run: start the HTTP API locally
run:
	go run ./cmd/server

## test: run unit tests
test:
	go test ./...

tidy:
	go mod tidy

## docker-up / docker-down / docker-logs: manage the API + PostgreSQL stack
docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api
