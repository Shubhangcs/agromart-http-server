-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id          UUID      PRIMARY KEY DEFAULT gen_random_uuid (),
    sender_id   UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    receiver_id UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content     TEXT      NOT NULL,
    is_read     BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;

-- +goose StatementEnd
