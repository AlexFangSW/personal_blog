#/bin/bash

# swag fmt && \
swag init --output ./swagger_docs --parseDependency -g ./cmd/server/main.go && \
goose -dir=./db/migrations/sqlite sqlite3 ./blog.db up && \
go run ./cmd/server/main.go
