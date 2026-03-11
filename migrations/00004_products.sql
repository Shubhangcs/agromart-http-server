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

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS product_images;

-- +goose StatementEnd