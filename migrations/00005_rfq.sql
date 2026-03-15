-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS rfqs (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        business_id UUID REFERENCES businesses (id) ON DELETE CASCADE,
        category_id UUID REFERENCES categories (id) ON DELETE CASCADE,
        sub_category_id UUID REFERENCES sub_categories (id) ON DELETE CASCADE,
        product_name TEXT NOT NULL,
        quantity NUMERIC(20, 1) NOT NULL,
        unit VARCHAR(50) NOT NULL,
        price NUMERIC(20, 2) NOT NULL,
        is_rfq_active BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rfqs;

-- +goose StatementEnd
