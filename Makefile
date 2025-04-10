# Makefile commands for Golang container operations (using docker-compose)

# Service name (fixed)
.PHONY: api.install api.test api.test.cover api.help api.install-air api.dev api.ent.init api.ent.generate api.ent.run api.db.connect api.ent.install api.ent.setup api.create-infrastructure api.create-domain api.create-repository api.create-usecase api.create-handler api.setup-all api.db.setup api.test.health api.test.hello api.test.users.create api.test.users.list api.test.users.get api.test.users.update api.test.users.delete api.test.all api.install-lint api.bash api.lint api.ent.gen api.run

init: api.install

up:
	docker-compose up -d

down:
	docker-compose down

# Install packages
api.install:
	docker-compose exec api go mod download
	docker-compose exec api go mod tidy

# Run tests (excluding ent directory)
api.test:
	docker-compose exec api go test $(shell docker-compose exec -T api go list ./... | grep -v "/ent")

# Run tests with coverage (excluding ent directory)
api.test.cover:
	docker-compose exec api go test -cover $(shell docker-compose exec -T api go list ./... | grep -v "/ent")

api.bash:
	docker-compose exec api bash

# Run with Air hot reload
api.dev:
	docker-compose exec api /go/bin/air

api.lint:
	docker-compose exec api go fmt ./...
	docker-compose exec api go vet ./...
	docker-compose exec api golangci-lint run --fix
	docker-compose exec api modernize -fix ./...

# Complete ent setup - install, init, and generate
api.ent.setup: api.ent.init api.ent.generate api.install

# Initialize Ent schema
api.ent.init:
	docker-compose exec api go run entgo.io/ent/cmd/ent init User

# Generate Ent code
api.ent.generate:
	docker-compose exec api go run entgo.io/ent/cmd/ent generate ./ent/schema

# Run the application
api.run:
	docker-compose exec api go run cmd/api/main.go

api.ent.gen:
	docker-compose exec api go run cmd/entgen/entgen.go

# Install golangci-lint
api.install-lint:
	docker-compose exec api go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "golangci-lint installed successfully! Now you can use 'make api.lint'"

# Help
help:
	@echo "Available commands:"
	@echo "  make api.install      - Install dependencies"
	@echo "  make api.test         - Run all tests"
	@echo "  make api.test.cover   - Run tests with coverage"
	@echo "  make api.lint         - Run code quality checks"
	@echo "  make api.install-lint - Install golangci-lint"
	@echo "  make api.install-air  - Install Air for hot reload"
	@echo "  make api.dev          - Run with Air hot reload"
	@echo "  make api.ent.setup    - Complete ent setup (install, init, generate)"
	@echo "  make api.ent.init     - Initialize Ent schema"
	@echo "  make api.ent.generate - Generate Ent code"
	@echo "  make api.run          - Run the application"

