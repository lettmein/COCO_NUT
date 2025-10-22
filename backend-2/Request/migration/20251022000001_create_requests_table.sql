-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logistic_points (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT,
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL
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

CREATE INDEX IF NOT EXISTS idx_requests_origin ON requests(origin_point_id);
CREATE INDEX IF NOT EXISTS idx_requests_dest ON requests(dest_point_id);
CREATE INDEX IF NOT EXISTS idx_requests_status ON requests(status);
CREATE INDEX IF NOT EXISTS idx_requests_status_ready ON requests(status, ready_at);
CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS logistic_points;
-- +goose StatementEnd
