-- +goose Up
-- +goose StatementBegin
ALTER TABLE receipts
ADD COLUMN recipient_email VARCHAR(50),
ADD COLUMN recipient_address TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE receipts
DROP COLUMN recipient_email,
DROP COLUMN recipient_address;
-- +goose StatementEnd
