-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    mode VARCHAR(10) CHECK (mode IN ('trial', 'paid')) DEFAULT 'trial',
    use_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
