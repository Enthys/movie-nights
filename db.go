package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectToDb() error {
	conn, err := sql.Open("postgres", reqEnv("MOVIENIGHT_DB_DSN"))

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
