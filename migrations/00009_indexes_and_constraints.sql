-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_businesses_user_id ON businesses (user_id);

CREATE INDEX IF NOT EXISTS idx_sub_categories_category_id ON sub_categories (category_id);

CREATE INDEX IF NOT EXISTS idx_products_business_id ON products (business_id);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products (category_id);
CREATE INDEX IF NOT EXISTS idx_products_sub_category_id ON products (sub_category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products (is_product_active);

CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images (product_id);

CREATE INDEX IF NOT EXISTS idx_product_ratings_product_id ON product_ratings (product_id);
CREATE INDEX IF NOT EXISTS idx_product_ratings_user_id ON product_ratings (user_id);

CREATE INDEX IF NOT EXISTS idx_product_reviews_product_id ON product_reviews (product_id);
CREATE INDEX IF NOT EXISTS idx_product_reviews_user_id ON product_reviews (user_id);

CREATE INDEX IF NOT EXISTS idx_rfqs_business_id ON rfqs (business_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_is_active ON rfqs (is_rfq_active);

CREATE INDEX IF NOT EXISTS idx_followers_user_id ON followers (user_id);
CREATE INDEX IF NOT EXISTS idx_followers_business_id ON followers (business_id);

CREATE INDEX IF NOT EXISTS idx_business_ratings_business_id ON business_ratings (business_id);
CREATE INDEX IF NOT EXISTS idx_business_ratings_user_id ON business_ratings (user_id);

CREATE INDEX IF NOT EXISTS idx_business_reviews_business_id ON business_reviews (business_id);
CREATE INDEX IF NOT EXISTS idx_business_reviews_user_id ON business_reviews (user_id);

CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages (sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages (receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages (
    LEAST(sender_id::text, receiver_id::text),
    GREATEST(sender_id::text, receiver_id::text),
    created_at ASC
);

CREATE INDEX IF NOT EXISTS idx_wishlists_user_id ON wishlists (user_id);
CREATE INDEX IF NOT EXISTS idx_wishlists_product_id ON wishlists (product_id);

DO $$
BEGIN
    -- Drop old single-column unique constraint if it exists (from old schema)
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'business_ratings_user_id_key'
          AND conrelid = 'business_ratings'::regclass
    ) THEN
        ALTER TABLE business_ratings DROP CONSTRAINT business_ratings_user_id_key;
    END IF;

    -- Fix rating column type from NUMERIC(1,1) to NUMERIC(2,1) for existing DBs
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'business_ratings'
          AND column_name = 'rating'
          AND numeric_precision = 1
    ) THEN
        ALTER TABLE business_ratings ALTER COLUMN rating TYPE NUMERIC(2, 1);
        ALTER TABLE business_ratings DROP CONSTRAINT IF EXISTS business_ratings_rating_check;
        ALTER TABLE business_ratings ADD CONSTRAINT business_ratings_rating_check
            CHECK (rating >= 0.5 AND rating <= 5.0);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_ratings DROP CONSTRAINT IF EXISTS unique_user_business_rating;

DROP INDEX IF EXISTS idx_wishlists_product_id;
DROP INDEX IF EXISTS idx_wishlists_user_id;

DROP INDEX IF EXISTS idx_messages_conversation;
DROP INDEX IF EXISTS idx_messages_receiver_id;
DROP INDEX IF EXISTS idx_messages_sender_id;

DROP INDEX IF EXISTS idx_business_reviews_user_id;
DROP INDEX IF EXISTS idx_business_reviews_business_id;

DROP INDEX IF EXISTS idx_business_ratings_user_id;
DROP INDEX IF EXISTS idx_business_ratings_business_id;

DROP INDEX IF EXISTS idx_followers_business_id;
DROP INDEX IF EXISTS idx_followers_user_id;

DROP INDEX IF EXISTS idx_rfqs_is_active;
DROP INDEX IF EXISTS idx_rfqs_business_id;

DROP INDEX IF EXISTS idx_product_reviews_user_id;
DROP INDEX IF EXISTS idx_product_reviews_product_id;

DROP INDEX IF EXISTS idx_product_ratings_user_id;
DROP INDEX IF EXISTS idx_product_ratings_product_id;

DROP INDEX IF EXISTS idx_product_images_product_id;

DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_sub_category_id;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_business_id;

DROP INDEX IF EXISTS idx_sub_categories_category_id;

DROP INDEX IF EXISTS idx_businesses_user_id;
-- +goose StatementEnd
