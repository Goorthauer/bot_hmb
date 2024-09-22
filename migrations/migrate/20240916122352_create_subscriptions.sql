-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions
(
    id          varchar primary key,
    user_id     varchar,
    school_id   varchar,
    price       varchar,
    days        int,
    created_at  timestamptz,
    deadline_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
