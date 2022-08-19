-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    username text,
    first_name text,
    last_name text,
    email text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
