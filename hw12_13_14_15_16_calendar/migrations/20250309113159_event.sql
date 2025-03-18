-- +goose Up
-- +goose StatementBegin
create table if not exists event (
    id serial PRIMARY KEY,
    title text NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    description text NOT NULL,
    user_id integer NOT NULL,
    notify_on integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists event;
-- +goose StatementEnd
