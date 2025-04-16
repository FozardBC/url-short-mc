-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls(
    url TEXT NOT NULL,
    alias TEXT UNIQUE NOT NULL
);

CREATE INDEX urls_idx ON urls USING hashmap (alias);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
