package main

import (
	"context"
	"errors"
	"log"
	"movie_night/types"
	"net/http"

	"github.com/gorilla/securecookie"
)

func notAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := store.Get(r, sessionCookieKey)
		if err != nil {
			if errors.Is(err, securecookie.ErrMacInvalid) {
				u.Options.MaxAge = -1
				if saveErr := u.Save(r, w); saveErr != nil {
					log.Println("failed to clear session.", saveErr)
				}
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong!"))
			log.Println("failed to retrieve user session. ", err)
			return
		}

		if !u.IsNew {
			http.Redirect(w, r, "/groups", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func userAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := store.Get(r, sessionCookieKey)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, envelope{"error": "something went wrong"}, nil)
			log.Println("failed to retrieve user session. ", err)
			return
		}

		if _, ok := u.Values[sk_authenticated]; !ok {
			writeJSON(w, http.StatusUnauthorized, envelope{
				"error": "not logged in",
			}, nil)
			return
		}

		user := &types.User{
			ID:        u.Values[sk_id].(int),
			Name:      u.Values[sk_name].(string),
			AvatarURL: u.Values[sk_avatar].(string),
			SocialId:  u.Values[sk_socialId].(string),
		}

		r = r.WithContext(context.WithValue(r.Context(), UserCtxKey, user))
		next(w, r)
	}
}
