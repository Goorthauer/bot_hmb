-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS users (
                                     id varchar PRIMARY KEY,
                                     is_master boolean,
                                     is_deleted bool,
                                     is_activated bool,
                                     firstname_encrypted varchar,
                                     lastname_encrypted varchar,
                                     phone varchar,
                                     username varchar, -- index
                                     deleted_at timestamp,
                                     deleted_by varchar,
                                     registered_at timestamp,
                                     pd_encryption_key varchar
);
CREATE INDEX IF NOT EXISTS users_username_index ON users(username);

CREATE TABLE IF NOT EXISTS telegram_accounts (
                                                 user_id varchar,
                                                 chat_id bigint,
                                                 is_active bool,
                                                 created_at timestamptz,
                                                 updated_at timestamptz,
                                                 PRIMARY KEY (user_id, chat_id)
    );
CREATE INDEX IF NOT EXISTS telegram_accounts_chat_id_idx ON telegram_accounts(chat_id);
CREATE UNIQUE INDEX IF NOT EXISTS telegram_accounts_chat_id_is_active_uniq_idx ON telegram_accounts(chat_id) WHERE is_active;
CREATE TABLE IF NOT EXISTS schools (
                                                     id varchar primary key,
                                                     name varchar,
                                                     city varchar,
                                                     address varchar,
                                                     contact varchar,
                                                     vk_link varchar,
                                                     region int
);
CREATE TABLE IF NOT EXISTS telegram_auth_tickets (
                                                     token varchar primary key,
                                                     user_id varchar,
                                                     created_at timestamptz,
                                                     updated_at timestamptz,
                                                     expires_at timestamptz,
                                                     is_spent bool,
                                                     is_blocked bool,
                                                     spent_at timestamptz
);
CREATE INDEX IF NOT EXISTS telegram_auth_tickets_user_id_created_at ON telegram_auth_tickets(user_id, created_at DESC);

INSERT INTO schools (id, name,vk_link,city,contact,address,region)
VALUES ('da68b3a4-a310-43ab-805f-159090d8cf55',
        'Школа ИСБ Ставрополь',
        'https://vk.com/hmbschoolstv',
        'Ставрополь',
        '+7-905-493-34-92',
        'Рогожникова 1',26);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS telegram_auth_tickets_user_id_created_at;
DROP TABLE IF EXISTS telegram_auth_tickets;
DROP INDEX IF EXISTS telegram_accounts_chat_id_is_active_uniq_idx;
DROP INDEX IF EXISTS telegram_accounts_chat_id_idx;
DROP TABLE IF EXISTS telegram_accounts;
DROP TABLE IF EXISTS schools;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS pg_trgm;
-- +goose StatementEnd
