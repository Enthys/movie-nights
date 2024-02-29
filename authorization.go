package main

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

type RequestCtxKey string

const (
	sessionCookieKey = "sid"
	sk_authenticated = "user_authenticated"
	sk_id            = "user_id"
	sk_name          = "user_name"
	sk_avatar        = "user_avatar"
	sk_socialId      = "user_socialId"

	UserCtxKey = RequestCtxKey("user")
)

var (
	sessionKey = []byte(reqEnv("SESSION_SECRET"))
	store      = sessions.NewCookieStore(sessionKey)
	maxAge     = 86400 * 30 // 30 days
)

func setupSessionStore() {
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = reqEnv("ENVIRONMENT") == "production"

	gothic.Store = store
}
