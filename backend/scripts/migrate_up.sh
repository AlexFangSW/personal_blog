#/bin/bash
migrate -database sqlite://blog.db -path ./db/migrations/ up
