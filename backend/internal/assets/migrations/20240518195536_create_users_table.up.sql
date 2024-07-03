CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
		email CITEXT UNIQUE NOT NULL,
		password_hash BYTEA NOT NULL,
		activated BOOL NOT NULL,
		version INTEGER NOT NULL DEFAULT 1,
		goal TEXT CHECK( goal IN ('cut', 'bulk', 'maintain') )
);