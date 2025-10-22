# Тут должен быть репозиторий Backend-1 на GO

Ясенко Данил Вячеславович

Для теста:
```
-- Маршрут: НИКИЭТ (id=22) -> СКЦ Росатома (id=38)
INSERT INTO routes(depart_at, max_weight_kg, max_volume_m3, status)
VALUES (now() + interval '90 minutes', 5000, 30, 'planned');

INSERT INTO route_points(route_id, seq_no, point_id) VALUES
  (1, 1, 22),
  (1, 2, 38);

INSERT INTO requests(origin_point_id, dest_point_id, weight_kg, volume_m3, ready_at, deadline_at, status) VALUES
  (22, 21, 300, 3.0, now(),           now() + interval '8 hour',  'pending'), -- Мячковский
  (22, 37, 400, 4.0, now(),           now() + interval '6 hour',  'pending'), -- Ферганская (НТЦ)
  (22, 38, 200, 2.0, now(),           now() + interval '12 hour', 'pending'), -- СКЦ (финальная точка)
  (22, 22, 100, 1.0, now(),           now() + interval '3 hour',  'pending'); -- доставить в стартовую точку (detour=0)
```

```
curl -X POST http://localhost:8001/routes/1/match

curl http://localhost:8001/routes/1/assignments
```
