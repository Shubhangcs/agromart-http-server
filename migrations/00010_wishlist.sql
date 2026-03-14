-- +goose Up
CREATE TABLE IF NOT EXISTS wishlists (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_wishlists_user_id    ON wishlists (user_id);
CREATE INDEX IF NOT EXISTS idx_wishlists_product_id ON wishlists (product_id);

-- +goose Down
DROP TABLE IF EXISTS wishlists;
