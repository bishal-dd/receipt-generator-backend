-- +goose Up
-- +goose StatementBegin
CREATE TABLE encrypted_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT NOT NULL,
    rate TEXT NOT NULL,
    quantity TEXT NOT NULL,
    amount TEXT NOT NULL,
    encrypted_receipt_id UUID NOT NULL,
    aes_key_encrypted TEXT,
    aes_iv TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (encrypted_receipt_id) REFERENCES encrypted_receipts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE encrypted_services;
-- +goose StatementEnd
