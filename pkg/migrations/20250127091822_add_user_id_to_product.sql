-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN user_id VARCHAR(36) NOT NULL;

ALTER TABLE products
    ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
    DROP CONSTRAINT fk_user;

ALTER TABLE products
    DROP COLUMN user_id;
-- +goose StatementEnd
