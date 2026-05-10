.PHONY: all auth lms runner gamification gateway migrate seed build clean

# Build individual services
auth:
	go build -o bin/auth-service ./cmd/auth-service

lms:
	go build -o bin/lms-service ./cmd/lms-service

runner:
	go build -o bin/code-runner ./cmd/code-runner

gamification:
	go build -o bin/gamification-service ./cmd/gamification-service

gateway:
	go build -o bin/gateway ./cmd/gateway

# Build all
all: auth lms runner gamification gateway

build: all

# Database
migrate:
	bash scripts/migrate.sh

seed:
	bash scripts/seed.sh

backup:
	bash scripts/backup.sh

# Docker
up:
	cd deploy && docker compose up --build -d

down:
	cd deploy && docker compose down

logs:
	cd deploy && docker compose logs -f

# Development
fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

test:
	go test ./... -v

clean:
	rm -rf bin/
