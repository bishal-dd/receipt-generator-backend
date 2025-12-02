-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD COLUMN type VARCHAR(50) NOT NULL DEFAULT 'product'
CHECK (type IN ('product', 'service'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN type;
-- +goose StatementEnd
