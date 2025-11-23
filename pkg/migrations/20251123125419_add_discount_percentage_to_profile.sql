-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
ADD COLUMN discount_percentage DECIMAL(5,2) DEFAULT 0.00;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
DROP COLUMN discount_percentage;
-- +goose StatementEnd
