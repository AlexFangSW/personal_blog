#/bin/bash

goose -dir=./db/migrations sqlite3 ./blog.db up && \
go run main.go
