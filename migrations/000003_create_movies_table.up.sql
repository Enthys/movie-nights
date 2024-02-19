CREATE TABLE IF NOT EXISTS movies(
	id SERIAL PRIMARY KEY,
	imdb_link TEXT NOT NULL,
	genres TEXT[],
	added_by INT NOT NULL
)
