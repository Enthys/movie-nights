package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"movie_night/types"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func getMovieFromIMDB(link string) (*types.Movie, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request. %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed. %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code returned: %d", resp.StatusCode)
	}

	movie := types.Movie{}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	if sel := doc.Find(".hero__primary-text"); sel.Length() == 0 {
		return nil, fmt.Errorf("failed to retrieve movie title")
	} else {
		movie.Name = sel.Text()
	}

	if sel := doc.Find("span[data-testid=\"plot-xs_to_m\"]"); sel.Length() == 0 {
		return nil, fmt.Errorf("failed to retrieve movie description")
	} else {
		movie.Description = sel.Text()
	}

	if sel := doc.Find(`div[data-testid="genres"] span.ipc-chip__text`); sel.Length() == 0 {
		return nil, fmt.Errorf("failed to retrieve movie genres")
	} else {
		sel.Map(func(i int, s *goquery.Selection) string {
			fmt.Println(i, s.Text())
			return s.Text()
		})
		movie.Genres = append(movie.Genres, sel.Map(func(i int, s *goquery.Selection) string { return s.Text() })...)
	}

	posterReq, err := http.NewRequest(http.MethodGet, "https://api.themoviedb.org/3/search/movie?api_key=15d2ea6d0dc1d476efbca3eba2b9bbfb&query="+url.QueryEscape(movie.Name), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie poster retrieval request. %w", err)
	}

	posterResp, err := http.DefaultClient.Do(posterReq)
	if err != nil {
		return nil, fmt.Errorf("request for poster retrieval failed. %w", err)
	}

	var result struct {
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
	}

	posterBody, _ := io.ReadAll(posterResp.Body)
	if err = json.Unmarshal(posterBody, &result); err != nil {
		return nil, fmt.Errorf("failed to retrieve poster address. %w", err)
	}

	movie.AvatarLink = "http://image.tmdb.org/t/p/w500" + result.Results[0].PosterPath

	return &movie, nil
}
