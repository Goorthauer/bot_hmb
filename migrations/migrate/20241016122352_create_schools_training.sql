-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS schools_training
(
    school_id varchar,
    price int,
    schedule jsonb,
    description varchar,
    created_at timestamptz,
    PRIMARY KEY (created_at, school_id)
);

INSERT INTO schools_training (school_id, price, schedule, description, created_at)
VALUES ('da68b3a4-a310-43ab-805f-159090d8cf55',
        2500,
        '[
          {
            "day": "пн",
            "time": {
              "open": "20-00",
              "closed": "22-00"
            },
            "description": "отработка приемов"
          },
          {
            "day": "ср",
            "time": {
              "open": "20-00",
              "closed": "22-00"
            },
            "description": "спаринги и железо"
          },
          {
            "day": "пт",
            "time": {
              "open": "20-00",
              "closed": "22-00"
            },
            "description": "кроссфит"
          }
        ]',
        'Тренер многократный чемпион различных соревнований по ИСБ, Иванов Иван Иваныч.',
        '2024-09-19 07:19:29.158970 +00:00');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schools_training;
-- +goose StatementEnd
