package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"movie_night/validator"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/markbates/goth/gothic"
)

func setupHandlers(mux *http.ServeMux) {
	mux.Handle("GET /api/login/{provider}", notAuthenticated(googleLoginHandler))
	mux.Handle("GET /api/login/{provider}/callback", notAuthenticated(googleLoginCallbackHandler))
	mux.Handle("GET /api/logout", userAuthenticated(logoutHandler))
	mux.Handle("GET /api/profile", userAuthenticated((profileHandler)))
	mux.Handle("GET /api/avatar", userAuthenticated(avatarHandler))
	mux.Handle("GET /api/groups", userAuthenticated(groupsHandler))
	mux.Handle("GET /api/groups/search", userAuthenticated((searchGroupsHandler)))
	mux.Handle("POST /api/groups", userAuthenticated(createGroupHandler))
	mux.Handle("GET /api/groups/{id}", userAuthenticated(viewGroupHandler))
	mux.Handle("GET /api/movies", userAuthenticated(getMoviesHandler))
	mux.Handle("POST /api/movies", userAuthenticated(addMovieHandler))

	staticDir := "./client/dist/"
	fs := http.FileServer(http.Dir(staticDir))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			fullPath := staticDir + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			fmt.Println(fullPath)
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				r.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, r)
	})
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Set("provider", "google")
	r.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(w, r)
}

func googleLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	user, err := getOrCreateUser(gothUser.UserID, gothUser.FirstName, gothUser.LastName, gothUser.AvatarURL)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	s, err := store.Get(r, sessionCookieKey)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	s.Values[sk_authenticated] = true
	s.Values[sk_id] = user.ID
	s.Values[sk_socialId] = user.SocialId
	s.Values[sk_name] = user.Name
	s.Values[sk_avatar] = user.AvatarURL

	if err = s.Save(r, w); err != nil {
		fmt.Fprintln(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "loggedIn",
		Value:    "1",
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := store.Get(r, sessionCookieKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to retrieve session"))
		return
	}
	cookie.Options.MaxAge = -1
	cookie.Save(r, w)
	http.SetCookie(w, &http.Cookie{
		Name:     "loggedIn",
		Value:    "1",
		MaxAge:   -1,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	if err := writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		log.Println(err)
		internalErrorResponse(w)
	}
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)

	req, _ := http.NewRequest("GET", user.AvatarURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "Something went wrong")
		return
	}

	io.Copy(w, resp.Body)
}

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	var req CreateGroupRequest
	if err := readJSON(w, r, &req); err != nil {
		internalErrorResponse(w)
		return
	}

	v := validator.New()
	v.Check(len(req.Name) > 3 && len(req.Name) < 25, "name", "Group name has to be at least 4 characters and at most 24 characters long.")
	v.Check(len(req.Description) <= 300, "description", "Description should be at most 300 characters long.")

	if !v.Valid() {
		validationErrorResponse(w, v)
		return
	}

	existingGroup, err := getGroupByName(req.Name)
	if err != nil && !errors.Is(err, ErrGroupNotFound) {
		v.AddError("internal", "Something went wrong, try again later.")
		validationErrorResponse(w, v)
		return
	}

	if existingGroup != nil {
		v.AddError("name", "Group name is already taken.")
		validationErrorResponse(w, v)
		return
	}

	newGroup, err := createGroup(req.Name, req.Description, user.ID)
	if err != nil {
		v.AddError("internal", "Something went wrong, try again later.")
		validationErrorResponse(w, v)
		return
	}

	if err = addUserToGroup(user.ID, newGroup.ID); err != nil {
		// Todo: rollback group creation if user was not added to group.
		fmt.Println(err)
	}

	if err = writeJSON(w, http.StatusCreated, envelope{"group": newGroup}, nil); err != nil {
		fmt.Println(err)
		internalErrorResponse(w)
	}
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	userGroups, err := getUserGroups(user)
	if err != nil {
		fmt.Println(err)
		internalErrorResponse(w)
		return
	}

	if err = writeJSON(w, http.StatusOK, envelope{"groups": userGroups}, nil); err != nil {
		fmt.Println(err)
		internalErrorResponse(w)
		return
	}
}

func viewGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupId, err := pathArgIntVal(r, "id")
	if err != nil {
		badRequestErrorResponse(w, envelope{"error": "invalid path argument"})
		return
	}

	group, err := getGroupById(groupId)
	if err != nil {
		switch err {
		case ErrGroupNotFound:
			notFoundResponse(w, envelope{"group": "group does not exist"})
		default:
			internalErrorResponse(w)
		}
		return
	}

	if err = writeJSON(w, http.StatusOK, envelope{"group": group}, nil); err != nil {
		internalErrorResponse(w)
	}
}

func searchGroupsHandler(w http.ResponseWriter, r *http.Request) {
	groupNameSearch := r.URL.Query().Get("name")
	if groupNameSearch == "" {
		badRequestErrorResponse(w, envelope{
			"name": "name should not be empty",
		})
		return
	}
	if len(groupNameSearch) > 300 {
		badRequestErrorResponse(w, envelope{
			"name": "name is too long",
		})
		return
	}

	foundGroups, err := searchGroupsByName(groupNameSearch)
	if err != nil {
		internalErrorResponse(w)
		return
	}

	if err = writeJSON(w, http.StatusOK, envelope{
		"groups": foundGroups,
	}, nil); err != nil {
		fmt.Println(err)
		internalErrorResponse(w)
		return
	}
}

func getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	page, err := queryIntVal(r, "page", 1)
	if err != nil {
		badRequestErrorResponse(w, envelope{"error": "invalid page value"})
		return
	}
	search := queryVal(r, "name", "")

	movies, err := getUserMovies(user, page, search)
	if err != nil {
		internalErrorResponse(w)
		return
	}

	if err = writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil); err != nil {
		internalErrorResponse(w)
		return
	}
}

type AddMovieRequest struct {
	MovieLink string `json:"link"`
	GroupId   *int   `json:"groupId"`
}

func addMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req AddMovieRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequestErrorResponse(w, envelope{"error": "invalid request"})
		return
	}
	link := req.MovieLink

	v := validator.New()
	v.Check(len(link) > 0, "link", "link is required")
	v.Check(strings.HasPrefix(link, "https://www.imdb.com/title") || strings.HasPrefix(link, "https://imdb.com/title"), "link", "invalid IMDB link provided")

	movieUrl, err := url.Parse(link)
	if err != nil {
		v.Check(false, "link", "invalid IMDB link provided")
	}

	if !v.Valid() {
		validationErrorResponse(w, v)
		return
	}

	link = strings.Trim("https://"+movieUrl.Host+movieUrl.Path, "/")

	existingMovie, err := getMovieByLink(link)
	if err != nil && !errors.Is(err, ErrMovieNotFound) {
		log.Println(err)
		internalErrorResponse(w)
		return
	}

	user := extractUser(r)

	if existingMovie != nil {
		if err = addToUserMovies(user.ID, existingMovie.ID); err != nil {
			log.Println(err)
			internalErrorResponse(w)
			return
		}

		if err := writeJSON(w, http.StatusOK, envelope{"movie": existingMovie}, nil); err != nil {
			log.Println(err)
			internalErrorResponse(w)
		}
		return
	}

	movie, err := getMovieFromIMDB(link)
	if err != nil {
		log.Println(err)
		internalErrorResponse(w)
		return
	}

	movie.AddedBy = user.ID
	if err = saveMovie(movie); err != nil {
		log.Println(err)
		internalErrorResponse(w)
		return
	}

	if err = addToUserMovies(user.ID, movie.ID); err != nil {
		log.Println(err)
		internalErrorResponse(w)
		return
	}

	if err := writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		log.Println(err)
		internalErrorResponse(w)
	}
}
