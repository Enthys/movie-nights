package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movie_night/types"
	"time"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrGroupNotFound = errors.New("group not found")
	ErrConflict      = errors.New("conflict occurred")
)

func getOrCreateUser(socialId, firstName, lastName, avatarURL string) (*types.User, error) {
	user, err := getUserBySocial(socialId)

	if err != nil {
		switch err {
		case ErrUserNotFound:
			newUser := types.User{
				SocialId:  socialId,
				Name:      fmt.Sprintf("%s %s", firstName, lastName),
				AvatarURL: avatarURL,
			}

			return createUser(newUser)
		default:
			return nil, err
		}
	}

	return user, nil
}

func getUserBySocial(id string) (*types.User, error) {
	var user types.User

	query := `SELECT id, name, social_id, avatar_url FROM users WHERE social_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.SocialId,
		&user.AvatarURL,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func createUser(user types.User) (*types.User, error) {
	query := `INSERT INTO users (social_id, name, avatar_url) VALUES ($1, $2, $3) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.QueryRowContext(ctx, query, user.SocialId, user.Name, user.AvatarURL).Scan(&user.ID); err != nil {
		return nil, err
	}

	return &user, nil
}

func createGroup(name, description string, creatorId int) (*types.Group, error) {
	group := types.Group{
		Name:        name,
		Description: description,
		CreatedBy:   creatorId,
	}

	query := `INSERT INTO groups (name, description, created_by) VALUES ($1, $2, $3) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.QueryRowContext(ctx, query, name, description, creatorId).Scan(&group.ID); err != nil {
		return nil, err
	}

	return &group, nil
}

func getGroupById(id int) (*types.Group, error) {
	var group types.Group

	query := `SELECT id, name, description, created_by FROM groups WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.QueryRowContext(ctx, query, id).Scan(&group.ID, &group.Name, &group.Description, &group.CreatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}

		return nil, err
	}

	return &group, nil
}

func getGroupByName(name string) (*types.Group, error) {
	var group types.Group

	query := `SELECT id, name, description, created_by FROM groups WHERE name = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.QueryRowContext(ctx, query, name).Scan(&group.ID, &group.Name, &group.Description, &group.CreatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}

		return nil, err
	}

	return &group, nil
}

func addUserToGroup(userId, groupId int) error {
	query := `INSERT INTO group_users(user_id, group_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, userId, groupId)

	return err
}

func getUserGroups(user *types.User) ([]*types.Group, error) {
	var groups []*types.Group

	query := `
		SELECT groups.id, groups.name, groups.description, groups.created_by
		FROM groups as groups
		JOIN group_users AS gu ON gu.group_id = groups.id
		WHERE gu.user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group types.Group

		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedBy); err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return groups, nil
}

func searchGroupsByName(name string) ([]*types.Group, error) {
	var groups []*types.Group

	query := `
		SELECT groups.id, groups.name, groups.description, groups.created_by
		FROM groups as groups
		LEFT JOIN group_users AS gu ON gu.group_id = groups.id
		WHERE groups.name ILIKE $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var group types.Group

		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedBy); err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return groups, nil
}

func getUserMovies(user *types.User, page int, nameSearch string) ([]*types.Movie, error) {
	movies := []*types.Movie{}

	query := `
		SELECT movies.id, movies.movie_name, movies.movie_description, movies.imdb_link, movies.genres FROM movies as movies
		JOIN user_movies as user_movies ON user_movies.movie_id = movies.id
		WHERE user_movies.user_id = $1
		AND (to_tsvector('simple', movie_name) @@ plainto_tsquery('simple', $2) OR $2 = '')
		LIMIT $3 OFFSET $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, user.ID, nameSearch, 10, page*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie types.Movie
		if err = rows.Scan(&movie.ID, &movie.Name, &movie.Description, &movie.IMDBLink, pq.Array(&movie.Genres)); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return movies, nil
}
