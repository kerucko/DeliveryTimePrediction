-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists results (
    id uuid primary key,
    delivery_time float
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE results
-- +goose StatementEnd
