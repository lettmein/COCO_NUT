#!/bin/bash

set -e

BASE_URL_REQUEST="http://localhost:8081"
BASE_URL_OFFICE="http://localhost:8082"

echo "==================================="
echo "Testing COCO_NUT Backend Services"
echo "==================================="

echo ""
echo "1. Health checks..."
curl -s "$BASE_URL_REQUEST/health" | jq .
curl -s "$BASE_URL_OFFICE/health" | jq .

echo ""
echo "2. Creating test request (заявка)..."
REQUEST_RESPONSE=$(curl -s -X POST "$BASE_URL_REQUEST/api/v1/requests" \
  -H "Content-Type: application/json" \
  -d '{
    "logistic_point_id": 1,
    "customer": {
      "company_name": "ООО Тестовая компания",
      "inn": "1234567890",
      "contact_name": "Иванов Иван Иванович",
      "phone": "+79001234567",
      "email": "test@example.com"
    },
    "cargo": {
      "name": "Тестовый груз",
      "quantity": 10,
      "weight": 1500,
      "volume": 25,
      "special_requirements": "Хрупкое"
    },
    "recipient": {
      "company_name": "ООО Получатель",
      "address": "г. Москва, ул. Примерная, д. 1",
      "contact_name": "Петров Петр Петрович",
      "phone": "+79007654321"
    }
  }')

REQUEST_ID=$(echo $REQUEST_RESPONSE | jq -r '.id')
echo "Created request with ID: $REQUEST_ID"
echo $REQUEST_RESPONSE | jq .

echo ""
echo "3. Creating test route (маршрут)..."
ROUTE_RESPONSE=$(curl -s -X POST "$BASE_URL_OFFICE/api/v1/routes" \
  -H "Content-Type: application/json" \
  -d '{
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
  }')

ROUTE_ID=$(echo $ROUTE_RESPONSE | jq -r '.id')
echo "Created route with ID: $ROUTE_ID"
echo $ROUTE_RESPONSE | jq .

echo ""
echo "4. Getting all requests..."
curl -s "$BASE_URL_REQUEST/api/v1/requests" | jq .

echo ""
echo "5. Getting all routes..."
curl -s "$BASE_URL_OFFICE/api/v1/routes" | jq .

echo ""
echo "6. Getting request by ID..."
curl -s "$BASE_URL_REQUEST/api/v1/requests/$REQUEST_ID" | jq .

echo ""
echo "7. Getting route by ID..."
curl -s "$BASE_URL_OFFICE/api/v1/routes/$ROUTE_ID" | jq .

echo ""
echo "==================================="
echo "All tests passed!"
echo "==================================="
