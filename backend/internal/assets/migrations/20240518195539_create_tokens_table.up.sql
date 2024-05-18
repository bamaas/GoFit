CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    expiry timestamp(0) NOT NULL,
    scope STRING NOT NULL
);