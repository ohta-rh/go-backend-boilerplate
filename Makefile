# Makefile commands for Golang container operations (using docker-compose)

# Service name (fixed)

.PHONY: add-pkg install test test.cover test.pkg lint install-lint help api.install-air api.dev api.ent.init api.ent.generate api.ent.run api.db.connect api.ent.install api.ent.setup api.create-infrastructure api.create-domain api.create-repository api.create-usecase api.create-handler api.setup-all api.db.setup api.test.health api.test.hello api.test.users.create api.test.users.list api.test.users.get api.test.users.update api.test.users.delete api.test.all api.install-lint

init: api.install
# Add package
add-pkg:
	@read -p "Enter package name to add: " pkg; \
	docker-compose exec api go get $$pkg

# Install packages
api.install:
	docker-compose exec api go mod download
	docker-compose exec api go mod tidy

# Run tests (excluding ent directory)
test:
	docker-compose exec api go test $(shell docker-compose exec -T api go list ./... | grep -v "/ent")

# Run tests with coverage (excluding ent directory)
test.cover:
	docker-compose exec api go test -cover $(shell docker-compose exec -T api go list ./... | grep -v "/ent")

# Run tests for specific package
test.pkg:
	@read -p "Enter package path to test: " pkg; \
	docker-compose exec api go test $$pkg

# Run linting
lint:
	docker-compose exec api go fmt ./...
	docker-compose exec api go vet ./...
	docker-compose exec api golangci-lint run --fix

api.bash:
	docker-compose exec api bash

# Run with Air hot reload
api.dev:
	docker-compose exec api /go/bin/air

api.lint:
	docker-compose exec api go fmt ./...
	docker-compose exec api go vet ./...
	docker-compose exec api golangci-lint run --fix

api.lint.check:
	docker-compose exec api go fmt ./...
	docker-compose exec api go vet ./...
	docker-compose exec api golangci-lint run

api.mod.edit:
	docker-compose exec api go mod edit -replace easy-go-backend=./

# Install ent
api.ent.install:
	docker-compose exec api go get entgo.io/ent/cmd/ent
	docker-compose exec api go get github.com/lib/pq

# Complete ent setup - install, init, and generate
api.ent.setup: api.ent.install api.ent.init api.ent.generate api.install

# Initialize Ent schema
api.ent.init:
	docker-compose exec api go run entgo.io/ent/cmd/ent init User

# Generate Ent code
api.ent.generate:
	docker-compose exec api go run entgo.io/ent/cmd/ent generate ./ent/schema

# Run the application
api.run:
	docker-compose exec api go run cmd/api/main.go

# Connect to database with psql
api.db.connect:
	docker-compose exec db psql -U root -d app

# Setup all directories
api.setup-all: api.create-infrastructure api.create-domain api.create-repository api.create-usecase api.create-handler api.ent.setup

# Setup complete application with DB
api.db.setup:
	@echo "Starting containers..."
	docker-compose up -d
	@echo "Setting up application directories..."
	make api.setup-all
	@echo "Running the application..."
	make api.run

api.ent.gen:
	docker-compose exec api go run cmd/entgen/entgen.go

# Install golangci-lint
api.install-lint:
	docker-compose exec api go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "golangci-lint installed successfully! Now you can use 'make api.lint'"

# Help
help:
	@echo "Available commands:"
	@echo "  make add-pkg          - Add a new Go package"
	@echo "  make api.install      - Install dependencies"
	@echo "  make test             - Run all tests"
	@echo "  make test-cover       - Run tests with coverage"
	@echo "  make test-pkg         - Run tests for a specific package"
	@echo "  make lint             - Run code quality checks"
	@echo "  make api.install-lint - Install golangci-lint"
	@echo "  make api.install-air  - Install Air for hot reload"
	@echo "  make api.dev          - Run with Air hot reload"
	@echo "  make api.ent.install  - Install ent"
	@echo "  make api.ent.setup    - Complete ent setup (install, init, generate)"
	@echo "  make api.ent.init     - Initialize Ent schema"
	@echo "  make api.ent.generate - Generate Ent code"
	@echo "  make api.run          - Run the application"
	@echo "  make api.db.connect   - Connect to the PostgreSQL database"
	@echo "  make api.setup-all    - Set up all project directories"
	@echo "  make api.db.setup     - Setup complete application with DB"


