package main

import (
	"fmt"
	"io"
	"movie_night/ui/components"
	"movie_night/ui/layout"
	"movie_night/ui/page"
	"movie_night/validator"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func setupHandlers(mux *http.ServeMux) {
	mux.Handle("GET /login/{provider}", notAuthenticated(googleLoginHandler))
	mux.Handle("GET /login/{provider}/callback", notAuthenticated(googleLoginCallbackHandler))
	mux.Handle("GET /logout", userAuthenticated(logoutHandler))
	mux.Handle("GET /avatar", userAuthenticated(avatarHandler))
	mux.Handle("GET /groups", userAuthenticated(groupsHandler))
	mux.Handle("GET /groups/search", userAuthenticated((searchGroupsHandler)))
	mux.Handle("GET /groups/create", userAuthenticated(createGroupFormHandler))
	mux.Handle("POST /groups", userAuthenticated(createGroupHandler))
	mux.Handle("GET /groups/{id}", userAuthenticated(viewGroupHandler))
	mux.Handle("GET /movies", userAuthenticated(myMoviesHandler))

	staticDir := "."
	mux.Handle("GET /assets/", http.FileServer(http.Dir(staticDir)))
	mux.Handle("GET /", notAuthenticated(indexHandler))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	layout.NewIndex(nil).WithBody(page.LoginPage()).Render(r.Context(), w)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func createGroupFormHandler(w http.ResponseWriter, r *http.Request) {
	layout.NewIndex(extractUser(r)).WithBody(page.GroupsCreate(validator.New())).Render(r.Context(), w)
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	name := r.FormValue("name")
	description := r.FormValue("description")

	v := validator.New()
	fmt.Println(len(name))
	v.Check(len(name) > 3 && len(name) < 25, "name", "Group name has to be at least 4 characters and at most 24 characters long.")
	v.Check(len(description) <= 300, "description", "Description should be at most 300 characters long.")

	if !v.Valid() {
		layout.NewIndex(extractUser(r)).WithBody(page.GroupsCreate(v)).Render(r.Context(), w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newGroup, err := createGroup(name, description, user.ID)

	if err != nil {
		fmt.Println(err)
	}

	if err = addUserToGroup(user.ID, newGroup.ID); err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	user := extractUser(r)
	userGroups, err := getUserGroups(user)
	if err != nil {
		internalErrorResponse(w)
		return
	}
	var groups []components.Group
	for _, userGroup := range userGroups {
		groups = append(groups, components.NewGroup(userGroup.Name, userGroup.Description, "foo"))
	}

	layout.NewIndex(extractUser(r)).WithBody(page.Groups(groups)).Render(r.Context(), w)
}

func viewGroupHandler(w http.ResponseWriter, r *http.Request) {
	layout.NewIndex(extractUser(r)).WithBody(page.ViewGroup()).Render(r.Context(), w)
}

func searchGroupsHandler(w http.ResponseWriter, r *http.Request) {
	groupNameSearch := r.URL.Query().Get("name")
	if groupNameSearch == "" {
		badRequestErrorResponse(w, map[string]string{
			"name": "name should not be empty",
		})
		return
	}
	if len(groupNameSearch) > 300 {
		badRequestErrorResponse(w, map[string]string{
			"name": "name is too long",
		})
		return
	}

	foundGroups, err := searchGroupsByName(groupNameSearch)
	if err != nil {
		internalErrorResponse(w)
		return
	}

	var groups []components.Group
	for _, foundGroup := range foundGroups {
		groups = append(groups, components.NewGroup(foundGroup.Name, foundGroup.Description, ""))
	}

	components.GroupCollection(groups, "You are not a part of any groups.").Render(r.Context(), w)
}

func myMoviesHandler(w http.ResponseWriter, r *http.Request) {
	layout.NewIndex(extractUser(r)).WithBody(page.Movies()).Render(r.Context(), w)
}

// func searchGroupsHandler(w http.ResponseWriter, r *http.Request) {
// 	name := r.FormValue("name")

// 	foundGroups, err := searchGroupsByName(name)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var groups []page.Group
// 	for _, foundGroup := range foundGroups {
// 		groups = append(groups, page.Group{
// 			Name: foundGroup.Name,
// 		})
// 	}

// 	page.GroupCollection("found-groups", groups).Render(r.Context(), w)
// }

// func dashboardHandler(w http.ResponseWriter, r *http.Request) {
// 	user := r.Context().Value(userCtxKey).(*types.User)

// 	userGroups, err := getUserGroups(user)
// 	if err != nil {
// 		fmt.Fprintln(w, "Something went terribly wrong")
// 		return
// 	}

// 	var groups []page.Group
// 	for _, userGroup := range userGroups {
// 		groups = append(groups, page.Group{
// 			Name: userGroup.Name,
// 		})
// 	}

// 	layout.IndexLayout{
// 		Authenticated: true,
// 	}.Layout(
// 		user,
// 		page.Dashboard(
// 			groups,
// 			[]page.Movie{
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 				{Name: "American Psycho", AddedDate: time.Now()},
// 			},
// 		)).Render(r.Context(), w)
// }
