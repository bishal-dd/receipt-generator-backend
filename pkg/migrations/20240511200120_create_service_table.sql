-- +goose Up
-- +goose StatementBegin
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description VARCHAR(5000) NOT NULL,
    rate FLOAT NOT NULL,
    quantity INT NOT NULL,
    amount INT NOT NULL,
    receipt_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE services;
-- +goose StatementEnd
