-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
ADD COLUMN currency VARCHAR(50) NOT NULL DEFAULT 'USD',
ADD COLUMN tax FLOAT NOT NULL DEFAULT 0.0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
DROP COLUMN recipient_email,
DROP COLUMN recipient_address;
-- +goose StatementEnd
