-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN mode SET DATA TYPE VARCHAR(15);
ALTER TABLE users DROP CONSTRAINT users_mode_check;
ALTER TABLE users ADD CONSTRAINT users_mode_check CHECK (mode IN ('trial', 'paid', 'starter', 'growth', 'business'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN mode SET DATA TYPE VARCHAR(10);
ALTER TABLE users DROP CONSTRAINT users_mode_check;
ALTER TABLE users ADD CONSTRAINT users_mode_check CHECK (mode IN ('trial', 'paid'));
-- +goose StatementEnd
