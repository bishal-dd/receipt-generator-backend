-- +goose Up
-- +goose StatementBegin
CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_name VARCHAR(255),
    recipient_name VARCHAR(255),
    recipient_phone VARCHAR(255),
    recipient_email VARCHAR(255),
    recipient_address TEXT,
    receipt_no VARCHAR(36),
    user_id VARCHAR(36) NOT NULL,
    date DATE NOT NULL,
    total_amount FLOAT DEFAULT 0,
    sub_total_amount FLOAT DEFAULT 0,
    tax_amount FLOAT DEFAULT 0,
    payment_method VARCHAR(50),
    payment_note TEXT,
    is_receipt_send BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE receipts;
-- +goose StatementEnd
