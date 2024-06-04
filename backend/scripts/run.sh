#/bin/bash

goose -dir=./db/migrations/sqlite sqlite3 ./blog.db up && \
go run main.go
