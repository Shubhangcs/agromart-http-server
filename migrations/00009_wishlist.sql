-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wishlists (
    id         UUID      PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product_id UUID      NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    UNIQUE (user_id, product_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wishlists;

-- +goose StatementEnd
