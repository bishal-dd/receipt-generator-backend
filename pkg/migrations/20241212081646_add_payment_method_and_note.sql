-- +goose Up
-- +goose StatementBegin
ALTER TABLE receipts
ADD COLUMN payment_method VARCHAR(50),
ADD COLUMN payment_note TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE receipts
DROP COLUMN payment_method,
DROP COLUMN payment_note;
-- +goose StatementEnd
