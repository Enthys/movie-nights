package main

import (
	"fmt"
	"movie_night/ui/layout"
	"movie_night/ui/page"
	"net/http"
	"time"

	"github.com/markbates/goth/gothic"
)

func setupHandlers(mux *http.ServeMux) {
	mux.Handle("GET /login/{provider}", notAuthenticated(googleLoginHandler))
	mux.Handle("GET /login/{provider}/callback", notAuthenticated(googleLoginCallbackHandler))
	mux.Handle("GET /logout", userAuthenticated(logoutHandler))
	mux.Handle("GET /dashboard", userAuthenticated(dashboardHandler))
	mux.Handle("GET /", notAuthenticated(indexHandler))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	layout.IndexLayout{
		Authenticated: false,
	}.Layout(page.LoginPage()).Render(r.Context(), w)
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Set("provider", "google")
	r.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(w, r)
}

func googleLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
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
	s.Values[sk_name] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)

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

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	layout.IndexLayout{
		Authenticated: true,
	}.Layout(
		page.Dashboard(
			[]page.Group{
				{Name: "Foo", MemberCount: 10},
				{Name: "Bar", MemberCount: 4},
			},
			[]page.Movie{
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
				{Name: "American Psycho", AddedDate: time.Now()},
			},
		)).Render(r.Context(), w)
}
