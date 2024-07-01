CREATE TABLE IF NOT EXISTS owners (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    login TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_login ON owners(login);

CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL
);