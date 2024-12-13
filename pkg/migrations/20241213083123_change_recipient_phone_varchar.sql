-- +goose Up
-- +goose StatementBegin
ALTER TABLE receipts
    ALTER COLUMN recipient_phone TYPE VARCHAR(255) USING recipient_phone::VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE receipts
    ALTER COLUMN recipient_phone TYPE INT USING recipient_phone::INTEGER;
-- +goose StatementEnd
