# Тут должен быть репозиторий Backend-1 на GO

Ясенко Данил Вячеславович

Для теста:

```
-- если таблиц ещё нет (их должны делать другие сервисы), создадим минимально:
CREATE TABLE IF NOT EXISTS routes (
    id BIGSERIAL PRIMARY KEY,
    depart_at TIMESTAMPTZ NOT NULL,
    max_weight_kg NUMERIC(12,3) NOT NULL CHECK (max_weight_kg > 0),
    max_volume_m3 NUMERIC(12,3) NOT NULL CHECK (max_volume_m3 > 0),
    status TEXT NOT NULL DEFAULT 'planned',
    current_volume NUMERIC(12,3) DEFAULT 0,
    current_weight NUMERIC(12,3) DEFAULT 0,
    route_points JSONB,
    request_ids TEXT[],
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS route_points (
    id BIGSERIAL PRIMARY KEY,
    route_id BIGINT NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    seq_no INT NOT NULL,
    point_id BIGINT NOT NULL REFERENCES logistic_points(id),
    UNIQUE(route_id, seq_no)
);

CREATE TABLE IF NOT EXISTS requests (
    id BIGSERIAL PRIMARY KEY,
    origin_point_id BIGINT REFERENCES logistic_points(id),
    dest_point_id BIGINT REFERENCES logistic_points(id),
    weight_kg NUMERIC(12,3) NOT NULL CHECK (weight_kg > 0),
    volume_m3 NUMERIC(12,3) NOT NULL CHECK (volume_m3 > 0),
    ready_at TIMESTAMPTZ NOT NULL,
    deadline_at TIMESTAMPTZ,
    customer_company_name VARCHAR(255) NOT NULL,
    customer_inn VARCHAR(12) NOT NULL,
    customer_contact_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    customer_email VARCHAR(255),
    cargo_name VARCHAR(255) NOT NULL,
    cargo_quantity INTEGER NOT NULL,
    cargo_special_requirements TEXT,
    recipient_company_name VARCHAR(255) NOT NULL,
    recipient_address TEXT NOT NULL,
    recipient_contact_name VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(50) NOT NULL,
    recipient_lat DOUBLE PRECISION,
    recipient_lon DOUBLE PRECISION,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

```
-- Маршрут: НИКИЭТ (id=22) -> СКЦ Росатома (id=38)
INSERT INTO routes(depart_at, max_weight_kg, max_volume_m3, status)
VALUES (now() + interval '90 minutes', 5000, 30, 'planned');

INSERT INTO route_points(route_id, seq_no, point_id) VALUES
  (1, 1, 22),
  (1, 2, 38);

INSERT INTO requests (
  origin_point_id, dest_point_id,
  weight_kg, volume_m3, ready_at, deadline_at,
  customer_company_name, customer_inn, customer_contact_name, customer_phone, customer_email,
  cargo_name, cargo_quantity, cargo_special_requirements,
  recipient_company_name, recipient_address, recipient_contact_name, recipient_phone, recipient_lat, recipient_lon,
  status
) VALUES
  -- 22 (НИКИЭТ) -> 21 (Мячковский)
  (22, 21,
   300, 3.0, now(), now() + interval '8 hour',
   'ООО Альфа', '7701234567', 'Иванов Иван', '+7-999-000-0001', 'client1@example.com',
   'Электроника', 10, NULL,
   'Получатель: Мячковский',
   (SELECT address FROM logistic_points WHERE id=21),
   'Петров Петр', '+7-999-100-0001',
   (SELECT lat FROM logistic_points WHERE id=21),
   (SELECT lon FROM logistic_points WHERE id=21),
   'pending'
  ),

  -- 22 (НИКИЭТ) -> 37 (НТЦ, Ферганская)
  (22, 37,
   400, 4.0, now(), now() + interval '6 hour',
   'ООО Альфа', '7701234567', 'Иванов Иван', '+7-999-000-0001', 'client1@example.com',
   'Комплектующие', 8, NULL,
   'Получатель: НТЦ АЭС',
   (SELECT address FROM logistic_points WHERE id=37),
   'Сидоров Сидор', '+7-999-100-0002',
   (SELECT lat FROM logistic_points WHERE id=37),
   (SELECT lon FROM logistic_points WHERE id=37),
   'pending'
  ),

  -- 22 (НИКИЭТ) -> 38 (СКЦ Росатома, финиш)
  (22, 38,
   200, 2.0, now(), now() + interval '12 hour',
   'ООО Альфа', '7701234567', 'Иванов Иван', '+7-999-000-0001', NULL,
   'Документы', 1, NULL,
   'Получатель: СКЦ Росатома',
   (SELECT address FROM logistic_points WHERE id=38),
   'Романов Роман', '+7-999-100-0003',
   (SELECT lat FROM logistic_points WHERE id=38),
   (SELECT lon FROM logistic_points WHERE id=38),
   'pending'
  ),

  -- 22 (НИКИЭТ) -> 22 (доставка в стартовую точку; detour=0)
  (22, 22,
   100, 1.0, now(), now() + interval '3 hour',
   'ООО Альфа', '7701234567', 'Иванов Иван', '+7-999-000-0001', 'client1@example.com',
   'Образцы', 5, 'Не кантовать',
   'Получатель: НИКИЭТ',
   (SELECT address FROM logistic_points WHERE id=22),
   'Дежурный', '+7-999-100-0004',
   (SELECT lat FROM logistic_points WHERE id=22),
   (SELECT lon FROM logistic_points WHERE id=22),
   'pending'
  );
```

```
curl -X POST http://localhost:8001/routes/1/match

curl http://localhost:8001/routes/1/assignments
```

```
BEGIN;
TRUNCATE TABLE routes_lists RESTART IDENTITY CASCADE;
TRUNCATE TABLE route_points RESTART IDENTITY CASCADE;
TRUNCATE TABLE routes RESTART IDENTITY CASCADE;
TRUNCATE TABLE requests RESTART IDENTITY CASCADE;
COMMIT;
```

```
update requests set volume_m3=10000.0 where id = 4;
```