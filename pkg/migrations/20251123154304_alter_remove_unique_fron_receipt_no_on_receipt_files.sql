-- +goose Up
-- +goose StatementBegin
ALTER TABLE receipt_files
DROP CONSTRAINT IF EXISTS receipt_files_receipt_no_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE receipt_files
ADD CONSTRAINT receipt_files_receipt_no_key UNIQUE (receipt_no);
-- +goose StatementEnd
