-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS push_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT NOT NULL UNIQUE,
    platform   VARCHAR(10) NOT NULL DEFAULT 'android',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_push_tokens_user_id ON push_tokens (user_id);

CREATE TABLE IF NOT EXISTS banners (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title      VARCHAR(100),
    image      TEXT,
    target_url TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS push_tokens;
-- +goose StatementEnd
