CREATE TABLE IF NOT EXISTS group_users(
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	group_id INT not NULL,
	CONSTRAINT fk_users
		FOREIGN KEY(user_id)
			REFERENCES users(id),
	CONSTRAINT fk_groups
		FOREIGN KEY(group_id)
			REFERENCES groups(id)
)
