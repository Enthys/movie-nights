package main

import (
	"log"
	"net/http"
)

func notAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := store.Get(r, sessionCookieKey)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("something went wrong!"))
			log.Println("failed to retrieve user session. ", err)
			return
		}

		if !u.IsNew {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func userAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := store.Get(r, sessionCookieKey)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("something went wrong!"))
			log.Println("failed to retrieve user session. ", err)
			return
		}

		if _, ok := u.Values[sk_authenticated]; !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}
