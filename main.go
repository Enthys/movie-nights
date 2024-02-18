package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func loadGoogleAuthentication() {
	goth.UseProviders(
		google.New(reqEnv("GOOGLE_API_ID"), reqEnv("GOOGLE_API_SECRET"), reqEnv("GOOGLE_API_REDIRECT"), "profile"),
	)
}

func main() {
	setupSessionStore()
	loadGoogleAuthentication()

	if err := loadTemplates(); err != nil {
		log.Fatalf("failed to load templates. %s", err)
	}

	if err := connectToDb(); err != nil {
		log.Fatalf("failed to establish database connection. %s", err)
	}

	mux := http.NewServeMux()

	setupHandlers(mux)

	server := http.Server{
		Addr:    "0.0.0.0:5000",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
