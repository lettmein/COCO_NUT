# API Контракты Backend-2

## Общая информация

**Base URLs:**
- Request Service: `http://localhost:8081`
- Office Service: `http://localhost:8082`

**Content-Type:** `application/json`

---

## Request Service API (Сервис заявок)

### 1. Создать заявку

**Endpoint:** `POST /api/v1/requests`

**Request Body:**
```json
{
  "logistic_point_id": 1,
  "customer": {
    "company_name": "ООО Компания",
    "inn": "1234567890",
    "contact_name": "Иванов Иван Иванович",
    "phone": "+79001234567",
    "email": "test@example.com"
  },
  "cargo": {
    "name": "Название груза",
    "quantity": 10,
    "weight": 1500.50,
    "volume": 25.5,
    "special_requirements": "Особые требования"
  },
  "recipient": {
    "company_name": "ООО Получатель",
    "address": "г. Москва, ул. Примерная, д. 1",
    "contact_name": "Петров Петр Петрович",
    "phone": "+79007654321"
  }
}
```

**Response:** `201 Created`
```json
{
  "id": "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
  "logistic_point_id": 1,
  "customer": {
    "company_name": "ООО Компания",
    "inn": "1234567890",
    "contact_name": "Иванов Иван Иванович",
    "phone": "+79001234567",
    "email": "test@example.com"
  },
  "cargo": {
    "name": "Название груза",
    "quantity": 10,
    "weight": 1500.5,
    "volume": 25.5,
    "special_requirements": "Особые требования"
  },
  "recipient": {
    "company_name": "ООО Получатель",
    "address": "г. Москва, ул. Примерная, д. 1",
    "contact_name": "Петров Петр Петрович",
    "phone": "+79007654321"
  },
  "status": "pending",
  "created_at": "2025-10-22T10:58:31.260932Z",
  "updated_at": "2025-10-22T10:58:31.260932Z"
}
```

---

### 2. Получить все заявки

**Endpoint:** `GET /api/v1/requests`

**Query Parameters:**
- `status` (optional) - фильтр по статусу (например: `pending`, `in_transit`, `delivered`)

**Пример:** `GET /api/v1/requests?status=pending`

**Response:** `200 OK`
```json
[
  {
    "id": "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
    "logistic_point_id": 1,
    "customer": {
      "company_name": "ООО Компания",
      "inn": "1234567890",
      "contact_name": "Иванов Иван Иванович",
      "phone": "+79001234567",
      "email": "test@example.com"
    },
    "cargo": {
      "name": "Название груза",
      "quantity": 10,
      "weight": 1500.5,
      "volume": 25.5,
      "special_requirements": "Особые требования"
    },
    "recipient": {
      "company_name": "ООО Получатель",
      "address": "г. Москва, ул. Примерная, д. 1",
      "contact_name": "Петров Петр Петрович",
      "phone": "+79007654321"
    },
    "status": "pending",
    "created_at": "2025-10-22T10:58:31.260932Z",
    "updated_at": "2025-10-22T10:58:31.260932Z"
  }
]
```

---

### 3. Получить заявку по ID

**Endpoint:** `GET /api/v1/requests/:id`

**Пример:** `GET /api/v1/requests/447fa6fb-b611-4f5a-a0cf-87f3ed6fb026`

**Response:** `200 OK`
```json
{
  "id": "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
  "logistic_point_id": 1,
  "customer": {
    "company_name": "ООО Компания",
    "inn": "1234567890",
    "contact_name": "Иванов Иван Иванович",
    "phone": "+79001234567",
    "email": "test@example.com"
  },
  "cargo": {
    "name": "Название груза",
    "quantity": 10,
    "weight": 1500.5,
    "volume": 25.5,
    "special_requirements": "Особые требования"
  },
  "recipient": {
    "company_name": "ООО Получатель",
    "address": "г. Москва, ул. Примерная, д. 1",
    "contact_name": "Петров Петр Петрович",
    "phone": "+79007654321"
  },
  "status": "pending",
  "created_at": "2025-10-22T10:58:31.260932Z",
  "updated_at": "2025-10-22T10:58:31.260932Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "Request not found"
}
```

---

### 4. Обновить статус заявки

**Endpoint:** `PATCH /api/v1/requests/:id/status`

**Request Body:**
```json
{
  "status": "in_transit"
}
```

**Возможные статусы:**
- `pending` - ожидает обработки
- `in_transit` - в пути
- `delivered` - доставлено
- `cancelled` - отменено

**Response:** `200 OK`
```json
{
  "message": "Request status updated successfully"
}
```

---

### 5. Удалить заявку

**Endpoint:** `DELETE /api/v1/requests/:id`

**Response:** `200 OK`
```json
{
  "message": "Request deleted successfully"
}
```

---

### 6. Health Check

**Endpoint:** `GET /health`

**Response:** `200 OK`
```json
{
  "status": "ok",
  "service": "request-service"
}
```

---

## Office Service API (Сервис маршрутов)

### 1. Создать маршрут

**Endpoint:** `POST /api/v1/routes`

**Request Body:**
```json
{
  "max_volume": 87.0,
  "max_weight": 16050.0,
  "departure_date": "2025-10-24T08:00:00Z",
  "route_points": [
    {
      "latitude": 54.93869544589972,
      "longitude": 43.31275114766075,
      "address": "ул. Зернова, 34, Саров, Нижегородская обл."
    },
    {
      "latitude": 56.79469832748941,
      "longitude": 61.312923642561195,
      "address": "ул. Ленинградская, 27, Заречный, Свердловская обл."
    }
  ]
}
```

**Поля:**
- `max_volume` - максимальный объем фуры (м³)
- `max_weight` - максимальный вес фуры (кг)
- `departure_date` - дата и время отправления (ISO 8601)
- `route_points` - массив точек маршрута (минимум 2)
  - `latitude` - широта
  - `longitude` - долгота
  - `address` - адрес точки

**Response:** `201 Created`
```json
{
  "id": "4ffd15b4-9e56-40a8-8ffb-ff26f13c7b53",
  "max_volume": 87,
  "max_weight": 16050,
  "current_volume": 0,
  "current_weight": 0,
  "departure_date": "2025-10-24T08:00:00Z",
  "route_points": [
    {
      "latitude": 54.93869544589972,
      "longitude": 43.31275114766075,
      "address": "ул. Зернова, 34, Саров, Нижегородская обл.",
      "arrival_time": "2025-10-24T08:00:00Z"
    },
    {
      "latitude": 56.79469832748941,
      "longitude": 61.312923642561195,
      "address": "ул. Ленинградская, 27, Заречный, Свердловская обл.",
      "arrival_time": "2025-10-25T02:58:00Z"
    }
  ],
  "status": "pending",
  "request_ids": [],
  "created_at": "2025-10-22T11:00:07.2065Z",
  "updated_at": "2025-10-22T11:00:07.2065Z"
}
```

**Примечание:** Время прибытия (`arrival_time`) рассчитывается автоматически на основе:
- Расстояния между точками (формула Haversine)
- Скорости движения: 60 км/ч

---

### 2. Получить все маршруты

**Endpoint:** `GET /api/v1/routes`

**Query Parameters:**
- `status` (optional) - фильтр по статусу

**Пример:** `GET /api/v1/routes?status=pending`

**Response:** `200 OK`
```json
[
  {
    "id": "4ffd15b4-9e56-40a8-8ffb-ff26f13c7b53",
    "max_volume": 87,
    "max_weight": 16050,
    "current_volume": 0,
    "current_weight": 0,
    "departure_date": "2025-10-24T08:00:00Z",
    "route_points": [
      {
        "latitude": 54.93869544589972,
        "longitude": 43.31275114766075,
        "address": "ул. Зернова, 34, Саров, Нижегородская обл.",
        "arrival_time": "2025-10-24T08:00:00Z"
      },
      {
        "latitude": 56.79469832748941,
        "longitude": 61.312923642561195,
        "address": "ул. Ленинградская, 27, Заречный, Свердловская обл.",
        "arrival_time": "2025-10-25T02:58:00Z"
      }
    ],
    "status": "pending",
    "request_ids": [],
    "created_at": "2025-10-22T11:00:07.2065Z",
    "updated_at": "2025-10-22T11:00:07.2065Z"
  }
]
```

---

### 3. Получить маршрут по ID

**Endpoint:** `GET /api/v1/routes/:id`

**Пример:** `GET /api/v1/routes/4ffd15b4-9e56-40a8-8ffb-ff26f13c7b53`

**Response:** `200 OK`
```json
{
  "id": "4ffd15b4-9e56-40a8-8ffb-ff26f13c7b53",
  "max_volume": 87,
  "max_weight": 16050,
  "current_volume": 35.5,
  "current_weight": 2500,
  "departure_date": "2025-10-24T08:00:00Z",
  "route_points": [
    {
      "latitude": 54.93869544589972,
      "longitude": 43.31275114766075,
      "address": "ул. Зернова, 34, Саров, Нижегородская обл.",
      "arrival_time": "2025-10-24T08:00:00Z"
    },
    {
      "latitude": 56.79469832748941,
      "longitude": 61.312923642561195,
      "address": "ул. Ленинградская, 27, Заречный, Свердловская обл.",
      "arrival_time": "2025-10-25T02:58:00Z"
    }
  ],
  "status": "pending",
  "request_ids": ["447fa6fb-b611-4f5a-a0cf-87f3ed6fb026"],
  "created_at": "2025-10-22T11:00:07.2065Z",
  "updated_at": "2025-10-22T11:00:07.2065Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "Route not found"
}
```

---

### 4. Обновить статус маршрута

**Endpoint:** `PATCH /api/v1/routes/:id/status`

**Request Body:**
```json
{
  "status": "in_progress"
}
```

**Возможные статусы:**
- `pending` - ожидает отправки
- `in_progress` - в пути
- `completed` - завершен
- `cancelled` - отменен

**Response:** `200 OK`
```json
{
  "message": "Route status updated successfully"
}
```

---

### 5. Удалить маршрут

**Endpoint:** `DELETE /api/v1/routes/:id`

**Response:** `200 OK`
```json
{
  "message": "Route deleted successfully"
}
```

---

### 6. Health Check

**Endpoint:** `GET /health`

**Response:** `200 OK`
```json
{
  "status": "ok",
  "service": "office-service"
}
```

---

## Типы данных

### Customer (Заказчик)
```typescript
interface Customer {
  company_name: string;  // Наименование организации
  inn: string;           // ИНН (10-12 цифр)
  contact_name: string;  // ФИО контактного лица
  phone: string;         // Телефон (формат: +7XXXXXXXXXX)
  email: string;         // Email (опционально)
}
```

### Cargo (Груз)
```typescript
interface Cargo {
  name: string;                    // Наименование груза
  quantity: number;                // Количество (коробок/пакетов)
  weight: number;                  // Общий вес (кг)
  volume: number;                  // Объем (м³)
  special_requirements: string;    // Особые требования (опционально)
}
```

### Recipient (Получатель)
```typescript
interface Recipient {
  company_name: string;  // Наименование организации
  address: string;       // Полный адрес
  contact_name: string;  // ФИО ответственного
  phone: string;         // Контактный телефон
}
```

### RoutePoint (Точка маршрута)
```typescript
interface RoutePoint {
  latitude: number;      // Широта
  longitude: number;     // Долгота
  address: string;       // Адрес точки
  arrival_time?: string; // Время прибытия (ISO 8601, автоматически рассчитывается)
}
```

### Request (Заявка)
```typescript
interface Request {
  id: string;              // UUID заявки
  logistic_point_id: number; // ID логистической точки
  customer: Customer;
  cargo: Cargo;
  recipient: Recipient;
  status: string;          // Статус заявки
  created_at: string;      // Дата создания (ISO 8601)
  updated_at: string;      // Дата обновления (ISO 8601)
}
```

### Route (Маршрут)
```typescript
interface Route {
  id: string;              // UUID маршрута
  max_volume: number;      // Максимальный объем (м³)
  max_weight: number;      // Максимальный вес (кг)
  current_volume: number;  // Текущий загруженный объем (м³)
  current_weight: number;  // Текущий загруженный вес (кг)
  departure_date: string;  // Дата отправления (ISO 8601)
  route_points: RoutePoint[]; // Массив точек маршрута
  status: string;          // Статус маршрута
  request_ids: string[];   // Массив ID связанных заявок
  created_at: string;      // Дата создания (ISO 8601)
  updated_at: string;      // Дата обновления (ISO 8601)
}
```

---

## Коды ошибок

| Код | Описание |
|-----|----------|
| 200 | OK - Успешный запрос |
| 201 | Created - Ресурс успешно создан |
| 400 | Bad Request - Неверный формат запроса |
| 404 | Not Found - Ресурс не найден |
| 500 | Internal Server Error - Внутренняя ошибка сервера |

### Формат ошибки
```json
{
  "error": "Описание ошибки"
}
```

---

## Примеры использования

### JavaScript/TypeScript (fetch)

```typescript
// Создание заявки
const createRequest = async () => {
  const response = await fetch('http://localhost:8081/api/v1/requests', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      logistic_point_id: 1,
      customer: {
        company_name: "ООО Компания",
        inn: "1234567890",
        contact_name: "Иванов Иван",
        phone: "+79001234567",
        email: "test@example.com"
      },
      cargo: {
        name: "Оборудование",
        quantity: 10,
        weight: 1500.5,
        volume: 25.5,
        special_requirements: "Хрупкое"
      },
      recipient: {
        company_name: "ООО Получатель",
        address: "г. Москва, ул. Примерная, д. 1",
        contact_name: "Петров Петр",
        phone: "+79007654321"
      }
    })
  });

  const data = await response.json();
  return data;
};

// Получение всех заявок
const getAllRequests = async () => {
  const response = await fetch('http://localhost:8081/api/v1/requests');
  const data = await response.json();
  return data;
};

// Создание маршрута
const createRoute = async () => {
  const response = await fetch('http://localhost:8082/api/v1/routes', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      max_volume: 87.0,
      max_weight: 16050.0,
      departure_date: "2025-10-24T08:00:00Z",
      route_points: [
        {
          latitude: 55.7558,
          longitude: 37.6173,
          address: "Москва, Красная площадь"
        },
        {
          latitude: 59.9343,
          longitude: 30.3351,
          address: "Санкт-Петербург, Невский проспект"
        }
      ]
    })
  });

  const data = await response.json();
  return data;
};
```

---

## Особенности

### Автоматический расчет времени прибытия
- Используется формула Haversine для расчета расстояния между координатами
- Скорость движения фуры: **60 км/ч** (константа)
- Время прибытия рассчитывается для каждой точки маршрута
- Формат времени: ISO 8601 (UTC)

### Middleware для аудита
Все запросы автоматически логируются с информацией:
- Название сервиса
- HTTP метод и URL
- Тело запроса и ответа
- Статус код
- Время обработки (мс)
- Данные отправляются в сервис аудита (be-3)

### CORS
Оба сервиса настроены на прием запросов с любых источников для разработки.

---

## Примечания для фронтенда

1. **Валидация на клиенте**: Рекомендуется валидировать:
   - ИНН (10-12 цифр)
   - Телефоны (формат +7XXXXXXXXXX)
   - Email (стандартный формат)
   - Числовые значения (вес, объем > 0)
   - Минимум 2 точки для маршрута

2. **Обработка дат**: Все даты в формате ISO 8601 (UTC). Для отображения конвертируйте в локальную временную зону.

3. **Загрузка данных**: При создании маршрута время прибытия рассчитывается автоматически - не нужно передавать `arrival_time`.

4. **Фильтрация**: Используйте query параметр `status` для фильтрации заявок и маршрутов.

5. **Отображение на карте**: Координаты в формате `latitude/longitude` для использования с картами (Yandex Maps, Google Maps, etc.)
