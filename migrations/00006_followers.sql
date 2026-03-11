-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS followers (
        business_id UUID REFERENCES businesses (id) ON DELETE CASCADE,
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (business_id, user_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS followers ();

-- +goose StatementEnd