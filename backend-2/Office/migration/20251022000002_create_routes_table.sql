-- +goose Up
-- +goose StatementBegin
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

CREATE INDEX IF NOT EXISTS idx_routes_status ON routes(status);
CREATE INDEX IF NOT EXISTS idx_routes_depart_at ON routes(depart_at);
CREATE INDEX IF NOT EXISTS idx_routes_created_at ON routes(created_at);
CREATE INDEX IF NOT EXISTS idx_route_points_route ON route_points(route_id);
CREATE INDEX IF NOT EXISTS idx_routes_lists_route ON routes_lists(route_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes_lists;
DROP TABLE IF EXISTS route_points;
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd
