.PHONY: migrate-up migrate-down migrate-status build run

migrate-up:
	goose -dir migration postgres "postgresql://postgres:postgres@localhost:5432/office_db?sslmode=disable" up

migrate-down:
	goose -dir migration postgres "postgresql://postgres:postgres@localhost:5432/office_db?sslmode=disable" down

migrate-status:
	goose -dir migration postgres "postgresql://postgres:postgres@localhost:5432/office_db?sslmode=disable" status

build:
	go build -o bin/office-service ./cmd/app

run:
	go run ./cmd/app/main.go

docker-build:
	docker build -t office-service .

docker-run:
	docker run -p 8082:8082 office-service
