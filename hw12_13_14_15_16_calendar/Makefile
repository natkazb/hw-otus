MIGRATION_DIR=migrations
DB_STRING="postgres://user_cld:pass_cld@localhost:5432/calendar?sslmode=disable"

BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migration-create:
	goose -dir $(MIGRATION_DIR) create $(name) sql

migration-up:
	goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up

migration-down:
	goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) down

migration-status:
	goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) status

migration-reset:
	goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset

generate:
	rm -rf ./internal/server/grpc/pb
	mkdir -p ./internal/server/grpc/pb

	protoc \
		--go_out=./internal/server/grpc \
		--go-grpc_out=./internal/server/grpc \
		./api/*.proto