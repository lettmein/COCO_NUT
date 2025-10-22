-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    logistic_point_id INTEGER NOT NULL,

    customer_company_name VARCHAR(255) NOT NULL,
    customer_inn VARCHAR(12) NOT NULL,
    customer_contact_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50) NOT NULL,
    customer_email VARCHAR(255),

    cargo_name VARCHAR(255) NOT NULL,
    cargo_quantity INTEGER NOT NULL,
    cargo_weight DECIMAL(10, 2) NOT NULL,
    cargo_volume DECIMAL(10, 2) NOT NULL,
    cargo_special_requirements TEXT,

    recipient_company_name VARCHAR(255) NOT NULL,
    recipient_address TEXT NOT NULL,
    recipient_contact_name VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(50) NOT NULL,

    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_requests_logistic_point ON requests(logistic_point_id);
CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_created_at ON requests(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests;
-- +goose StatementEnd
