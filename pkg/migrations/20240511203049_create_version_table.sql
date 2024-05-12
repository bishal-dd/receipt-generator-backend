-- +goose Up
-- +goose StatementBegin
CREATE TABLE versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mode VARCHAR(10) CHECK (mode IN ('trial', 'paid')) DEFAULT 'trial',
    user_id VARCHAR(36) UNIQUE NOT NULL,
    use_count INT NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE versions;
-- +goose StatementEnd
