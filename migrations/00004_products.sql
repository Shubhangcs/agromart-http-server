-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS products (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        business_id UUID REFERENCES businesses (id) ON DELETE CASCADE,
        category_id UUID REFERENCES categories (id) ON DELETE CASCADE,
        sub_category_id UUID REFERENCES sub_categories (id) ON DELETE CASCADE,
        name VARCHAR(200) NOT NULL,
        description VARCHAR(1000) NOT NULL,
        quantity NUMERIC(20, 2) NOT NULL,
        unit VARCHAR(50) NOT NULL,
        price NUMERIC(20, 2) NOT NULL,
        moq VARCHAR(200) NOT NULL,
        is_product_active BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS product_images (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        image_index INTEGER NOT NULL,
        product_id UUID REFERENCES products (id) ON DELETE CASCADE,
        image TEXT NOT NULL,
        UNIQUE (product_id, image_index),
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS product_ratings (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        rating NUMERIC(2, 1) NOT NULL CHECK (
            rating >= 0.5
            AND rating <= 5.0
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (product_id, user_id)
    );

CREATE TABLE
    IF NOT EXISTS product_reviews (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        review TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_ratings;

DROP TABLE IF EXISTS product_reviews;

DROP TABLE IF EXISTS product_images;

DROP TABLE IF EXISTS products;

-- +goose StatementEnd
