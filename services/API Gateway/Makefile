APP_NAME := APIGateway
CMD_PATH := ./cmd/main
BUILD_DIR := ./bin
GO_FILES := $(shell go list ./...)

GO := go
LINT := golangci-lint

.PHONY: help build run test lint fmt vet clean tidy docker-build

help:
	@echo "Available commands:"
	@echo "  make build         - build binary"
	@echo "  make run           - run application"
	@echo "  make test          - run tests"
	@echo "  make lint          - run linters"
	@echo "  make fmt           - format code"
	@echo "  make vet           - run go vet"
	@echo "  make tidy          - go mod tidy"
	@echo "  make clean         - remove build artifacts"
	@echo "  make docker-build  - build docker image"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

run:
	$(GO) run $(CMD_PATH)

test:
	$(GO) test -race -cover ./...

lint:
	$(LINT) run

fmt:
	$(GO) fmt ./...
	$(GO) run golang.org/x/tools/cmd/goimports@latest -w .

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BUILD_DIR)

docker-build:
	docker build -t $(APP_NAME):latest .