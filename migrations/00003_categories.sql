-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS categories (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        category_image TEXT,
        name VARCHAR(100) NOT NULL,
        description VARCHAR(500) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS sub_categories (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        category_id UUID REFERENCES categories (id) ON DELETE CASCADE,
        sub_category_image TEXT,
        name VARCHAR(100) NOT NULL,
        description VARCHAR(500) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sub_categories;

DROP TABLE IF EXISTS categories;

-- +goose StatementEnd
