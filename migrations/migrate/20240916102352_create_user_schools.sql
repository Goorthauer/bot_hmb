-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_schools
(
    user_id   varchar PRIMARY KEY,
    school_id varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_schools;
-- +goose StatementEnd
