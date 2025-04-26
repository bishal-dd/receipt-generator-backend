-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_receipts_user_id ON receipts(user_id);
CREATE INDEX idx_receipts_date ON receipts(date);
CREATE INDEX idx_receipts_receipt_no ON receipts(receipt_no);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_receipts_user_id;
DROP INDEX IF EXISTS idx_receipts_date;
DROP INDEX IF EXISTS idx_receipts_receipt_no;
-- +goose StatementEnd
