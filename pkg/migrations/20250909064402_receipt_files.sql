-- +goose Up
-- +goose StatementBegin
CREATE TABLE receipt_files (
    id UUID PRIMARY KEY,
    receipt_no VARCHAR(100) NOT NULL UNIQUE,
    encrypted_receipt_id UUID NOT NULL,
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (encrypted_receipt_id) REFERENCES encrypted_receipts(id) ON DELETE CASCADE

);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS receipt_files;
-- +goose StatementEnd
