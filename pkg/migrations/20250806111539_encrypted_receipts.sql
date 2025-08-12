-- +goose Up
-- +goose StatementBegin
CREATE TABLE encrypted_receipts (
    id UUID PRIMARY KEY,
    receipt_name TEXT,
    recipient_name TEXT,
    recipient_phone TEXT,
    recipient_email TEXT,
    recipient_address TEXT,
    receipt_no TEXT,
    user_id VARCHAR(36) NOT NULL,
    date DATE NOT NULL,
    total_amount TEXT,
    sub_total_amount TEXT,
    tax_amount TEXT,
    payment_method TEXT,
    payment_note TEXT,
    is_receipt_send BOOLEAN NOT NULL DEFAULT false,
    aes_key_encrypted TEXT,
    aes_iv TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE encrypted_receipts;
-- +goose StatementEnd
