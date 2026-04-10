.PHONY: dev dev-backend dev-frontend build clean test test-coverage lint docker

# Development
dev-frontend:
	cd web && npm run dev

dev-backend:
	DEV=1 go run ./cmd/parcel-tracker

dev:
	$(MAKE) -j2 dev-backend dev-frontend

# Build
build: build-frontend build-backend

build-frontend:
	cd web && npm run build

build-backend:
	go build -ldflags="-s -w" -o bin/parcel-tracker ./cmd/parcel-tracker

# Test
test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Lint
lint:
	go vet ./...

# Docker
docker:
	docker build -t parcel-tracker .

# Clean
clean:
	rm -rf bin/ web/dist/ data/
