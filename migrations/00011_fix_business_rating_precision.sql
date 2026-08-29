-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_ratings ALTER COLUMN rating TYPE NUMERIC(2,1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_ratings ALTER COLUMN rating TYPE NUMERIC(1,1);
-- +goose StatementEnd
