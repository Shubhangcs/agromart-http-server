-- +goose Up
-- +goose StatementBegin

-- Indexes for FK columns to improve JOIN and WHERE performance
CREATE INDEX IF NOT EXISTS idx_businesses_user_id ON businesses(user_id);
CREATE INDEX IF NOT EXISTS idx_products_business_id ON products(business_id);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_sub_category_id ON products(sub_category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_product_active);
CREATE INDEX IF NOT EXISTS idx_sub_categories_category_id ON sub_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_business_id ON rfqs(business_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_is_active ON rfqs(is_rfq_active);
CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id);
CREATE INDEX IF NOT EXISTS idx_followers_user_id ON followers(user_id);
CREATE INDEX IF NOT EXISTS idx_followers_business_id ON followers(business_id);
CREATE INDEX IF NOT EXISTS idx_business_ratings_business_id ON business_ratings(business_id);

-- Fix business_ratings: drop the incorrect per-user unique constraint and add
-- the correct per-(business, user) unique constraint so a user can rate many
-- businesses but only once per business.
DO $$
BEGIN
    -- Drop the old unique constraint on user_id alone if it exists
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'business_ratings_user_id_key'
          AND conrelid = 'business_ratings'::regclass
    ) THEN
        ALTER TABLE business_ratings DROP CONSTRAINT business_ratings_user_id_key;
    END IF;

    -- Add the correct composite unique constraint
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'unique_user_business_rating'
          AND conrelid = 'business_ratings'::regclass
    ) THEN
        ALTER TABLE business_ratings
            ADD CONSTRAINT unique_user_business_rating UNIQUE (business_id, user_id);
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_businesses_user_id;
DROP INDEX IF EXISTS idx_products_business_id;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_sub_category_id;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_sub_categories_category_id;
DROP INDEX IF EXISTS idx_rfqs_business_id;
DROP INDEX IF EXISTS idx_rfqs_is_active;
DROP INDEX IF EXISTS idx_product_images_product_id;
DROP INDEX IF EXISTS idx_followers_user_id;
DROP INDEX IF EXISTS idx_followers_business_id;
DROP INDEX IF EXISTS idx_business_ratings_business_id;

ALTER TABLE business_ratings DROP CONSTRAINT IF EXISTS unique_user_business_rating;

-- +goose StatementEnd
