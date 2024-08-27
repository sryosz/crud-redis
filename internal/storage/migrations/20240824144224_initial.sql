-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email TEXT,
    password bytea
);

-- +goose Down
DROP TABLE IF EXISTS users
