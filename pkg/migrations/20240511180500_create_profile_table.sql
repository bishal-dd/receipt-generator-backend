-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(255) NOT NULL,
    logo_image TEXT,
    phone_no INT NOT NULL,
    address VARCHAR(255),
    email VARCHAR(255),
    city VARCHAR(255),
    title VARCHAR(255),
    signature_image TEXT,
    manual_signature_image VARCHAR(255),
    user_id VARCHAR(36) NOT NULL,
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
