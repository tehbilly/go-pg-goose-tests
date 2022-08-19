-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id int NOT NULL PRIMARY KEY,
    user_id int NOT NULL,
    title text NOT NULL,
    body text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
