-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    max_volume DECIMAL(10, 2) NOT NULL,
    max_weight DECIMAL(10, 2) NOT NULL,
    current_volume DECIMAL(10, 2) DEFAULT 0,
    current_weight DECIMAL(10, 2) DEFAULT 0,
    departure_date TIMESTAMP NOT NULL,
    route_points JSONB NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    request_ids TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_routes_status ON routes(status);
CREATE INDEX idx_routes_departure_date ON routes(departure_date);
CREATE INDEX idx_routes_created_at ON routes(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd
