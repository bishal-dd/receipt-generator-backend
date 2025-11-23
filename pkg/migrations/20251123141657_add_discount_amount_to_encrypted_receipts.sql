-- +goose Up
-- +goose StatementBegin
ALTER TABLE encrypted_receipts
ADD COLUMN discount_amount TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE encrypted_receipts
DROP COLUMN discount_amount;
-- +goose StatementEnd
