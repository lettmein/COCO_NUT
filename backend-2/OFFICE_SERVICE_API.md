# Office Service API

**Base URL:** `http://localhost:8082`

**Content-Type:** `application/json`

---

## Endpoints

### 1. Health Check

**GET** `/health`

**Response:**
```json
{
  "status": "ok",
  "service": "office-service"
}
```

---

### 2. Создать маршрут

**POST** `/api/v1/routes`

**Request Body:**
```json
{
  "max_volume": 100.5,
  "max_weight": 5000.0,
  "departure_date": "2025-10-23T08:00:00Z",
  "route_points": [
    {
      "latitude": 54.93869544589972,
      "longitude": 43.31275114766075,
      "address": "ул. Зернова, 34, Саров, Нижегородская обл."
    },
    {
      "latitude": 56.79469832748941,
      "longitude": 61.312923642561195,
      "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл."
    },
    {
      "latitude": 56.13157339089896,
      "longitude": 94.6077417485157,
      "address": "Первомайская ул., 7, Зеленогорск, Красноярский край"
    }
  ]
}
```

**Response: 201 Created**
```json
{
  "id": "a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f",
  "max_volume": 100.5,
  "max_weight": 5000,
  "current_volume": 0,
  "current_weight": 0,
  "departure_date": "2025-10-23T08:00:00Z",
  "route_points": [
    {
      "latitude": 54.93869544589972,
      "longitude": 43.31275114766075,
      "address": "ул. Зернова, 34, Саров, Нижегородская обл.",
      "arrival_time": "2025-10-23T08:00:00Z"
    },
    {
      "latitude": 56.79469832748941,
      "longitude": 61.312923642561195,
      "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
      "arrival_time": "2025-10-23T32:15:00Z"
    },
    {
      "latitude": 56.13157339089896,
      "longitude": 94.6077417485157,
      "address": "Первомайская ул., 7, Зеленогорск, Красноярский край",
      "arrival_time": "2025-10-24T20:45:00Z"
    }
  ],
  "status": "pending",
  "request_ids": [],
  "created_at": "2025-10-22T11:30:00.123456Z",
  "updated_at": "2025-10-22T11:30:00.123456Z"
}
```

**Важно:**
- `arrival_time` для каждой точки рассчитывается автоматически
- Первая точка получает время из `departure_date`
- Для последующих точек время рассчитывается по формуле Haversine с учетом скорости 60 км/ч
- Расстояние между точками считается по координатам (latitude, longitude)

---

### 3. Получить все маршруты

**GET** `/api/v1/routes`

**Query Parameters:**
- `status` (optional) - фильтр по статусу

**Примеры:**
- `GET /api/v1/routes` - все маршруты
- `GET /api/v1/routes?status=pending` - только со статусом pending

**Response: 200 OK**
```json
[
  {
    "id": "a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f",
    "max_volume": 100.5,
    "max_weight": 5000,
    "current_volume": 0,
    "current_weight": 0,
    "departure_date": "2025-10-23T08:00:00Z",
    "route_points": [
      {
        "latitude": 54.93869544589972,
        "longitude": 43.31275114766075,
        "address": "ул. Зернова, 34, Саров, Нижегородская обл.",
        "arrival_time": "2025-10-23T08:00:00Z"
      },
      {
        "latitude": 56.79469832748941,
        "longitude": 61.312923642561195,
        "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
        "arrival_time": "2025-10-23T32:15:00Z"
      }
    ],
    "status": "pending",
    "request_ids": [],
    "created_at": "2025-10-22T11:30:00.123456Z",
    "updated_at": "2025-10-22T11:30:00.123456Z"
  }
]
```

---

### 4. Получить маршрут по ID

**GET** `/api/v1/routes/:id`

**Пример:** `GET /api/v1/routes/a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f`

**Response: 200 OK**
```json
{
  "id": "a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f",
  "max_volume": 100.5,
  "max_weight": 5000,
  "current_volume": 35.5,
  "current_weight": 2500,
  "departure_date": "2025-10-23T08:00:00Z",
  "route_points": [
    {
      "latitude": 54.93869544589972,
      "longitude": 43.31275114766075,
      "address": "ул. Зернова, 34, Саров, Нижегородская обл.",
      "arrival_time": "2025-10-23T08:00:00Z"
    },
    {
      "latitude": 56.79469832748941,
      "longitude": 61.312923642561195,
      "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
      "arrival_time": "2025-10-23T32:15:00Z"
    }
  ],
  "status": "in_transit",
  "request_ids": ["447fa6fb-b611-4f5a-a0cf-87f3ed6fb026"],
  "created_at": "2025-10-22T11:30:00.123456Z",
  "updated_at": "2025-10-22T12:15:00.789012Z"
}
```

**Error: 404 Not Found**
```json
{
  "error": "Route not found"
}
```

---

### 5. Обновить статус маршрута

**PATCH** `/api/v1/routes/:id/status`

**Request Body:**
```json
{
  "status": "in_transit"
}
```

**Доступные статусы:**
- `pending` - ожидает выполнения
- `in_transit` - в пути
- `completed` - завершен
- `cancelled` - отменен

**Response: 200 OK**
```json
{
  "message": "Route status updated successfully"
}
```

---

### 6. Удалить маршрут

**DELETE** `/api/v1/routes/:id`

**Пример:** `DELETE /api/v1/routes/a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f`

**Response: 200 OK**
```json
{
  "message": "Route deleted successfully"
}
```

---

## Типы данных

### RoutePoint (Точка маршрута)
```typescript
interface RoutePoint {
  latitude: number;        // Широта (градусы)
  longitude: number;       // Долгота (градусы)
  address: string;         // Адрес точки
  arrival_time?: string;   // Время прибытия (ISO 8601, рассчитывается автоматически)
}
```

### Route (Маршрут)
```typescript
interface Route {
  id: string;                     // UUID маршрута
  max_volume: number;             // Максимальный объем в м³
  max_weight: number;             // Максимальный вес в кг
  current_volume: number;         // Текущий загруженный объем в м³
  current_weight: number;         // Текущий загруженный вес в кг
  departure_date: string;         // Дата/время отправления (ISO 8601)
  route_points: RoutePoint[];     // Массив точек маршрута
  status: string;                 // Статус маршрута
  request_ids?: string[];         // Массив ID связанных заявок
  created_at: string;             // Дата создания (ISO 8601)
  updated_at: string;             // Дата обновления (ISO 8601)
}
```

### CreateRouteDTO (Создание маршрута)
```typescript
interface CreateRouteDTO {
  max_volume: number;             // Максимальный объем
  max_weight: number;             // Максимальный вес
  departure_date: string;         // Дата/время отправления
  route_points: RoutePoint[];     // Точки маршрута (без arrival_time)
}
```

### UpdateRouteDTO (Обновление статуса)
```typescript
interface UpdateRouteDTO {
  status: string;                 // Новый статус
}
```

---

## Расчет времени прибытия

Сервис автоматически рассчитывает время прибытия в каждую точку маршрута:

### Алгоритм:
1. **Первая точка** получает `arrival_time` = `departure_date`
2. **Последующие точки** рассчитываются по формуле:
   - Расстояние между точками вычисляется по **формуле Haversine** (учитывает кривизну Земли)
   - Скорость движения: **60 км/ч** (константа)
   - Время в пути = расстояние / скорость
   - `arrival_time` = `arrival_time` предыдущей точки + время в пути

### Формула Haversine:
```
a = sin²(Δφ/2) + cos(φ1) × cos(φ2) × sin²(Δλ/2)
c = 2 × atan2(√a, √(1−a))
d = R × c
```
Где:
- φ1, φ2 — широта точек (в радианах)
- Δφ — разница широт
- Δλ — разница долгот
- R — радиус Земли (6371 км)
- d — расстояние в километрах

### Пример расчета:
```
Точка 1: Саров (54.9387, 43.3128) - отправление в 08:00
Точка 2: Заречный (56.7947, 61.3129)

Расстояние: ~1456 км
Время в пути: 1456 / 60 = 24.27 часа ≈ 24 часа 16 минут
Прибытие: 08:00 + 24:16 = 08:16 (следующий день)
```

---

## Коды ответов

| Код | Описание |
|-----|----------|
| 200 | OK - Успешный запрос |
| 201 | Created - Маршрут создан |
| 400 | Bad Request - Неверный формат |
| 404 | Not Found - Маршрут не найден |
| 500 | Internal Server Error |

### Формат ошибки
```json
{
  "error": "Описание ошибки"
}
```

---

## Примеры использования

### cURL

```bash
# Создать маршрут с 3 точками
curl -X POST http://localhost:8082/api/v1/routes \
  -H "Content-Type: application/json" \
  -d '{
    "max_volume": 100.5,
    "max_weight": 5000.0,
    "departure_date": "2025-10-23T08:00:00Z",
    "route_points": [
      {
        "latitude": 54.93869544589972,
        "longitude": 43.31275114766075,
        "address": "ул. Зернова, 34, Саров, Нижегородская обл."
      },
      {
        "latitude": 56.79469832748941,
        "longitude": 61.312923642561195,
        "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл."
      },
      {
        "latitude": 56.13157339089896,
        "longitude": 94.6077417485157,
        "address": "Первомайская ул., 7, Зеленогорск, Красноярский край"
      }
    ]
  }'

# Получить все маршруты
curl http://localhost:8082/api/v1/routes

# Получить маршруты со статусом pending
curl http://localhost:8082/api/v1/routes?status=pending

# Обновить статус маршрута
curl -X PATCH http://localhost:8082/api/v1/routes/a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f/status \
  -H "Content-Type: application/json" \
  -d '{"status": "in_transit"}'

# Удалить маршрут
curl -X DELETE http://localhost:8082/api/v1/routes/a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f
```

### JavaScript/TypeScript

```typescript
// Создать маршрут
const createRoute = async () => {
  const response = await fetch('http://localhost:8082/api/v1/routes', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      max_volume: 100.5,
      max_weight: 5000.0,
      departure_date: "2025-10-23T08:00:00Z",
      route_points: [
        {
          latitude: 54.93869544589972,
          longitude: 43.31275114766075,
          address: "ул. Зернова, 34, Саров, Нижегородская обл."
        },
        {
          latitude: 56.79469832748941,
          longitude: 61.312923642561195,
          address: "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл."
        }
      ]
    })
  });

  return await response.json();
};

// Получить все маршруты
const getAllRoutes = async () => {
  const response = await fetch('http://localhost:8082/api/v1/routes');
  return await response.json();
};

// Получить маршрут по ID
const getRoute = async (id: string) => {
  const response = await fetch(`http://localhost:8082/api/v1/routes/${id}`);
  return await response.json();
};

// Обновить статус маршрута
const updateRouteStatus = async (id: string, status: string) => {
  const response = await fetch(`http://localhost:8082/api/v1/routes/${id}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ status })
  });

  return await response.json();
};

// Удалить маршрут
const deleteRoute = async (id: string) => {
  const response = await fetch(`http://localhost:8082/api/v1/routes/${id}`, {
    method: 'DELETE'
  });

  return await response.json();
};
```

---

## Интеграция с Request Service

Маршруты связаны с заявками через поле `request_ids`:

```typescript
// Пример маршрута с привязанными заявками
{
  "id": "a5f7e8d9-c3b2-4a1e-9f6d-8c5b4a3e2d1f",
  "max_volume": 100.5,
  "max_weight": 5000,
  "current_volume": 35.5,
  "current_weight": 2500,
  "request_ids": [
    "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
    "8b2c9d3e-4f5a-6b7c-8d9e-0f1a2b3c4d5e"
  ],
  "status": "in_transit"
}
```

**Примечание:** Логика привязки заявок к маршрутам реализуется внешним сервисом оптимизации маршрутов.
