-- справочник логистических точек
CREATE TABLE IF NOT EXISTS logistic_points (
  id          BIGINT PRIMARY KEY,
  name        TEXT NOT NULL,
  address     TEXT,
  lat         DOUBLE PRECISION NOT NULL,
  lon         DOUBLE PRECISION NOT NULL
);

-- результаты матчинга для маршрутов
CREATE TABLE IF NOT EXISTS routes_lists (
  id          BIGSERIAL PRIMARY KEY,
  route_id    BIGINT NOT NULL,
  request_id  BIGINT NOT NULL,
  seq_no      INT NOT NULL,
  eta_plan    TIMESTAMPTZ,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(route_id, request_id),
  UNIQUE(route_id, seq_no)
);

CREATE INDEX IF NOT EXISTS idx_routes_lists_route ON routes_lists(route_id);
