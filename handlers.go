package main

import (
	"fmt"
	"net/http"

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
	t, ok := templates["index.page.html"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "could not find template"}`))
		return
	}

	if err := t.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "could not execute. %s"}`, err)))
		return
	}
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

	render(w, r, "dashboard.page.html", nil)
}
