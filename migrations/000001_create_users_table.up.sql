CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	social_id TEXT NOT NULL,
	name text NOT NULL,
	avatar_url TEXT NULL
);
