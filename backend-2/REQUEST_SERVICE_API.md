# Request Service API

**Base URL:** `http://localhost:8081`

**Content-Type:** `application/json`

---

## Endpoints

### 1. Health Check

**GET** `/health`

**Response:**
```json
{
  "status": "ok",
  "service": "request-service"
}
```

---

### 2. Создать заявку

**POST** `/api/v1/requests`

**Request Body:**
```json
{
  "logistic_point_id": 1,
  "customer": {
    "company_name": "ВНИИЭФ ФГУП",
    "inn": "5254001234",
    "contact_name": "Иванов Сергей Петрович",
    "phone": "+79201234567",
    "email": "ivanov@vniief.ru"
  },
  "cargo": {
    "name": "Специальное оборудование",
    "quantity": 15,
    "weight": 2500.00,
    "volume": 35.5,
    "special_requirements": "Требуется температурный режим +15-25°C"
  },
  "recipient": {
    "company_name": "Заречный филиал",
    "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
    "contact_name": "Петров Алексей Иванович",
    "phone": "+79207654321"
  }
}
```

**Response: 201 Created**
```json
{
  "id": "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
  "logistic_point_id": 1,
  "customer": {
    "company_name": "ВНИИЭФ ФГУП",
    "inn": "5254001234",
    "contact_name": "Иванов Сергей Петрович",
    "phone": "+79201234567",
    "email": "ivanov@vniief.ru"
  },
  "cargo": {
    "name": "Специальное оборудование",
    "quantity": 15,
    "weight": 2500,
    "volume": 35.5,
    "special_requirements": "Требуется температурный режим +15-25°C"
  },
  "recipient": {
    "company_name": "Заречный филиал",
    "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
    "contact_name": "Петров Алексей Иванович",
    "phone": "+79207654321"
  },
  "status": "pending",
  "created_at": "2025-10-22T10:58:31.260932Z",
  "updated_at": "2025-10-22T10:58:31.260932Z"
}
```

---

### 3. Получить все заявки

**GET** `/api/v1/requests`

**Query Parameters:**
- `status` (optional) - фильтр по статусу

**Примеры:**
- `GET /api/v1/requests` - все заявки
- `GET /api/v1/requests?status=pending` - только со статусом pending

**Response: 200 OK**
```json
[
  {
    "id": "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
    "logistic_point_id": 1,
    "customer": {
      "company_name": "ВНИИЭФ ФГУП",
      "inn": "5254001234",
      "contact_name": "Иванов Сергей Петрович",
      "phone": "+79201234567",
      "email": "ivanov@vniief.ru"
    },
    "cargo": {
      "name": "Специальное оборудование",
      "quantity": 15,
      "weight": 2500,
      "volume": 35.5,
      "special_requirements": "Требуется температурный режим +15-25°C"
    },
    "recipient": {
      "company_name": "Заречный филиал",
      "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
      "contact_name": "Петров Алексей Иванович",
      "phone": "+79207654321"
    },
    "status": "pending",
    "created_at": "2025-10-22T10:58:31.260932Z",
    "updated_at": "2025-10-22T10:58:31.260932Z"
  }
]
```

---

### 4. Получить заявку по ID

**GET** `/api/v1/requests/:id`

**Пример:** `GET /api/v1/requests/447fa6fb-b611-4f5a-a0cf-87f3ed6fb026`

**Response: 200 OK**
```json
{
  "id": "447fa6fb-b611-4f5a-a0cf-87f3ed6fb026",
  "logistic_point_id": 1,
  "customer": {
    "company_name": "ВНИИЭФ ФГУП",
    "inn": "5254001234",
    "contact_name": "Иванов Сергей Петрович",
    "phone": "+79201234567",
    "email": "ivanov@vniief.ru"
  },
  "cargo": {
    "name": "Специальное оборудование",
    "quantity": 15,
    "weight": 2500,
    "volume": 35.5,
    "special_requirements": "Требуется температурный режим +15-25°C"
  },
  "recipient": {
    "company_name": "Заречный филиал",
    "address": "ул. Ленинградская, 27, офис 3, Заречный, Свердловская обл.",
    "contact_name": "Петров Алексей Иванович",
    "phone": "+79207654321"
  },
  "status": "pending",
  "created_at": "2025-10-22T10:58:31.260932Z",
  "updated_at": "2025-10-22T10:58:31.260932Z"
}
```

**Error: 404 Not Found**
```json
{
  "error": "Request not found"
}
```

---

### 5. Обновить статус заявки

**PATCH** `/api/v1/requests/:id/status`

**Request Body:**
```json
{
  "status": "in_transit"
}
```

**Доступные статусы:**
- `pending` - ожидает обработки
- `in_transit` - в пути
- `delivered` - доставлено
- `cancelled` - отменено

**Response: 200 OK**
```json
{
  "message": "Request status updated successfully"
}
```

---

### 6. Удалить заявку

**DELETE** `/api/v1/requests/:id`

**Пример:** `DELETE /api/v1/requests/447fa6fb-b611-4f5a-a0cf-87f3ed6fb026`

**Response: 200 OK**
```json
{
  "message": "Request deleted successfully"
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
  phone: string;         // Телефон (+7XXXXXXXXXX)
  email: string;         // Email
}
```

### Cargo (Груз)
```typescript
interface Cargo {
  name: string;                    // Наименование груза
  quantity: number;                // Количество (коробок/пакетов)
  weight: number;                  // Общий вес в кг
  volume: number;                  // Объем в м³
  special_requirements: string;    // Особые требования
}
```

### Recipient (Получатель)
```typescript
interface Recipient {
  company_name: string;  // Наименование организации
  address: string;       // Полный адрес доставки
  contact_name: string;  // ФИО ответственного
  phone: string;         // Контактный телефон
}
```

### Request (Заявка)
```typescript
interface Request {
  id: string;                    // UUID заявки
  logistic_point_id: number;     // ID логистической точки
  customer: Customer;            // Данные заказчика
  cargo: Cargo;                  // Данные груза
  recipient: Recipient;          // Данные получателя
  status: string;                // Статус заявки
  created_at: string;            // Дата создания (ISO 8601)
  updated_at: string;            // Дата обновления (ISO 8601)
}
```

---

## Коды ответов

| Код | Описание |
|-----|----------|
| 200 | OK - Успешный запрос |
| 201 | Created - Заявка создана |
| 400 | Bad Request - Неверный формат |
| 404 | Not Found - Заявка не найдена |
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
# Создать заявку
curl -X POST http://localhost:8081/api/v1/requests \
  -H "Content-Type: application/json" \
  -d '{
    "logistic_point_id": 1,
    "customer": {
      "company_name": "ООО Компания",
      "inn": "1234567890",
      "contact_name": "Иванов Иван",
      "phone": "+79001234567",
      "email": "test@example.com"
    },
    "cargo": {
      "name": "Оборудование",
      "quantity": 10,
      "weight": 1500.5,
      "volume": 25.5,
      "special_requirements": "Хрупкое"
    },
    "recipient": {
      "company_name": "ООО Получатель",
      "address": "г. Москва, ул. Примерная, д. 1",
      "contact_name": "Петров Петр",
      "phone": "+79007654321"
    }
  }'

# Получить все заявки
curl http://localhost:8081/api/v1/requests

# Получить заявки со статусом pending
curl http://localhost:8081/api/v1/requests?status=pending

# Обновить статус
curl -X PATCH http://localhost:8081/api/v1/requests/447fa6fb-b611-4f5a-a0cf-87f3ed6fb026/status \
  -H "Content-Type: application/json" \
  -d '{"status": "in_transit"}'
```

### JavaScript/TypeScript

```typescript
// Создать заявку
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

  return await response.json();
};

// Получить все заявки
const getAllRequests = async () => {
  const response = await fetch('http://localhost:8081/api/v1/requests');
  return await response.json();
};

// Получить заявку по ID
const getRequest = async (id: string) => {
  const response = await fetch(`http://localhost:8081/api/v1/requests/${id}`);
  return await response.json();
};

// Обновить статус
const updateStatus = async (id: string, status: string) => {
  const response = await fetch(`http://localhost:8081/api/v1/requests/${id}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ status })
  });

  return await response.json();
};
```
