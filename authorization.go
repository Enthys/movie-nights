package main

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

const (
	sessionCookieKey = "sid"
	sk_authenticated = "authenticated"
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
