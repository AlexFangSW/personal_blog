#/bin/bash

swag init --parseDependency && \
goose -dir=./db/migrations/sqlite sqlite3 ./blog.db up && \
go run main.go
