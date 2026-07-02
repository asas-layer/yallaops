export DATABASE_URL := env_var_or_default("DATABASE_URL", "postgres://postgres:yalla@localhost:5432/yallaops?sslmode=disable")
export REDIS_URL    := env_var_or_default("REDIS_URL", "redis://localhost:6379")

# Regenerate all protobuf stubs
proto:
    buf generate

# Regenerate sqlc Go code
sqlc:
    cd core && sqlc generate

# Run all pending migrations
migrate:
    cd core && go run github.com/golang-migrate/migrate/v4/cmd/migrate \
        -path db/migrations \
        -database "$DATABASE_URL" \
        up

# Roll back the last migration
migrate-down:
    cd core && go run github.com/golang-migrate/migrate/v4/cmd/migrate \
        -path db/migrations \
        -database "$DATABASE_URL" \
        down 1

# Start the gRPC server in dev mode
dev-core:
    cd core && go run ./cmd/server

# Build the CLI binary
build-cli:
    cd cli && go build -o yallaops ./cmd/yallaops

# Run all tests
test:
    cd core && go test ./...
    cd cli && go test ./...

# Run linters
lint:
    cd core && golangci-lint run ./...
    cd cli && golangci-lint run ./...
    buf lint
