CREATE TABLE user_movies (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	movie_id INT NOT NULL,
	CONSTRAINT fk_user_movies_users
		FOREIGN KEY(user_id)
			REFERENCES users(id),
	CONSTRAINT fk_user_movies_movies
		FOREIGN KEY(movie_id)
			REFERENCES movies(id)
)
