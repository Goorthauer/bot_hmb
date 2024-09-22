-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_presents (
                                       id varchar PRIMARY KEY,
                                       user_id varchar,
                                       created_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_presents;
-- +goose StatementEnd
