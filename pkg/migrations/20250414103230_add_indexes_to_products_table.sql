-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_deleted_at;
-- +goose StatementEnd
