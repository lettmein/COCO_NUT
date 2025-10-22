# Request Service

Микросервис для управления заявками на перевозку грузов.

Разуваев Денис Русланович

## API Endpoints

### Health Check
```
GET /health
```

### Заявки
```
POST   /api/v1/requests          - Создать заявку
GET    /api/v1/requests          - Получить все заявки
GET    /api/v1/requests?status=pending - Получить заявки по статусу
GET    /api/v1/requests/:id      - Получить заявку по ID
PATCH  /api/v1/requests/:id/status - Обновить статус заявки
DELETE /api/v1/requests/:id      - Удалить заявку
```

## Миграции

Запуск миграций:
```bash
cd migration
goose postgres "postgresql://postgres:postgres@localhost:5432/request_db?sslmode=disable" up
```

## Запуск

Локально:
```bash
cp .env.example .env
go run cmd/app/main.go
```

Docker:
```bash
docker build -t request-service .
docker run -p 8081:8081 request-service
```
