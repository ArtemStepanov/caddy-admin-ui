.PHONY: all build run dev docker clean test docker-up-build

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/ArtemStepanov/caddy-admin-ui/internal/version.Version=$(VERSION)

# Build all
all: build

# Build backend
build:
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o bin/caddy-admin-ui ./cmd/server

# Build frontend
frontend:
	cd web && npm install && npm run build

# Run locally (development)
dev:
	go run ./cmd/server

# Run frontend dev server
dev-frontend:
	cd web && npm run dev

# Build Docker image
docker:
	docker build --build-arg VERSION=$(VERSION) -t caddy-admin-ui .

# Run with Docker Compose
docker-up:
	docker compose up -d

# Run with Docker Compose (local build)
docker-up-build:
	docker build -t ghcr.io/artemstepanov/caddy-admin-ui:latest . && docker compose up -d

# Stop Docker Compose
docker-down:
	docker compose down

# View logs
logs:
	docker compose logs -f

# Run tests
test:
	CGO_ENABLED=1 go test $$(go list ./... | grep -v '/web/node_modules/')
	cd web && npm test -- run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/node_modules/

# Install dependencies
deps:
	go mod download
	cd web && npm install
