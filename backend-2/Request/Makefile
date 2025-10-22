.PHONY: migrate-up migrate-down migrate-status build run

migrate-up:
	goose -dir migration postgres "postgresql://postgres:postgres@localhost:5432/request_db?sslmode=disable" up

migrate-down:
	goose -dir migration postgres "postgresql://postgres:postgres@localhost:5432/request_db?sslmode=disable" down

migrate-status:
	goose -dir migration postgres "postgresql://postgres:postgres@localhost:5432/request_db?sslmode=disable" status

build:
	go build -o bin/request-service ./cmd/app

run:
	go run ./cmd/app/main.go

docker-build:
	docker build -t request-service .

docker-run:
	docker run -p 8081:8081 request-service
