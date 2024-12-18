-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
ADD COLUMN phone_number_country_code VARCHAR(50) NOT NULL DEFAULT 'US';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profiles
DROP COLUMN phone_number_country_code;
-- +goose StatementEnd
