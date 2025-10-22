# Сборка бинарника
FROM golang:1.24.6-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY backend-1/go.mod backend-1/go.sum ./
RUN go mod download
COPY backend-1/ .
# статическая сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app ./cmd/server

# Рантайм без Air
FROM alpine:3.20 AS prod
RUN apk add --no-cache ca-certificates curl
WORKDIR /app
COPY --from=builder /app/app /app/app
# Копируем миграции и данные, чтобы applyMigrations и сидер работали внутри контейнера
COPY --from=builder /src/migrations /app/migrations
COPY --from=builder /src/data /app/data
EXPOSE 8080
ENTRYPOINT ["/app/app"]