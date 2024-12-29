-- +goose Up
-- +goose StatementBegin
ALTER TABLE receipts
ADD COLUMN is_receipt_send BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE receipts
DROP COLUMN is_receipt_send;
-- +goose StatementEnd
