-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS leads (
        lead_id             UUID          PRIMARY KEY DEFAULT gen_random_uuid (),
        enquirer_business_id UUID         NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
        enquire_to_id        UUID         NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
        product_id           UUID         NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        enquiry_message      TEXT           NOT NULL,
        order_quantity       NUMERIC(20, 2) NOT NULL DEFAULT 0,
        expected_price       NUMERIC(20, 2) NOT NULL DEFAULT 0,
        created_at           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT no_self_enquiry CHECK (enquirer_business_id <> enquire_to_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS leads;

-- +goose StatementEnd