package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectToDb() error {
	username := reqEnv("MOVIE_NIGHTS_DB_USER")
	password := url.QueryEscape(reqEnv("MOVIE_NIGHTS_DB_PASS"))
	host := reqEnv("MOVIE_NIGHTS_DB_HOST")
	port := reqEnvInt("MOVIE_NIGHTS_DB_PORT")
	name := reqEnv("MOVIE_NIGHTS_DB_NAME")
	args := reqEnv("MOVIE_NIGHTS_DB_ARGS")

	conn, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		username,
		password,
		host,
		port,
		name,
		args,
	))

	if err != nil {
		return fmt.Errorf("failed to connect to database. %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("faulty connection received. %w", err)
	}

	db = conn

	return nil
}
