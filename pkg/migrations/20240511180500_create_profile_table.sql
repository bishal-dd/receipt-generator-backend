-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(255),
    logo_image TEXT,
    phone_no VARCHAR(255),
    address VARCHAR(255),
    email VARCHAR(255),
    city VARCHAR(255),
    title VARCHAR(255),
    signature_image TEXT,
    manual_signature_image VARCHAR(255),
    user_id VARCHAR(36) NOT NULL,
    currency VARCHAR(50) NOT NULL DEFAULT 'USD',
    tax FLOAT NOT NULL DEFAULT 0.0,
    phone_number_country_code VARCHAR(50) NOT NULL DEFAULT 'US',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE profiles;
-- +goose StatementEnd
