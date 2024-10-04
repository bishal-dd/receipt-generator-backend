-- +goose Up
-- +goose StatementBegin
CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_name VARCHAR(100) NOT NULL,
    recipient_name VARCHAR(100) NOT NULL,
    recipient_phone INT NOT NULL,
    amount FLOAT NOT NULL,
    transaction_no INT,
    user_id VARCHAR(36) NOT NULL,
    date DATE NOT NULL,
    total_amount FLOAT DEFAULT 0,
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
