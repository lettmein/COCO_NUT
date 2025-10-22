# Dev образ с хот-релоадом
FROM golang:1.24.6-alpine AS base
WORKDIR /app
RUN apk add --no-cache git curl build-base
COPY backend-1/go.mod backend-1/go.sum ./
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY backend-1/ .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# Prod сборка
FROM golang:1.24.6-alpine AS builder
WORKDIR /src
COPY backend-1/go.mod backend-1/go.sum ./
RUN go mod download
COPY backend-1/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/server

# Минимальный рантайм
FROM gcr.io/distroless/static-debian12 AS prod
COPY --from=builder /bin/app /app
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app"]