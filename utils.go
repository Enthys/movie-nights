package main

import (
	"log"
	"movie_night/types"
	"net/http"
	"os"
	"strconv"
)

func internalErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Something went wrong!"))
}

func extractUser(r *http.Request) *types.User {
	user, ok := r.Context().Value(UserCtxKey).(*types.User)

	if !ok {
		panic("request with no user session passed authorization middleware")
	}

	return user
}

func reqEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("environment variable %s missing", key)
	}

	return val
}

func reqEnvInt(key string) int {
	val := reqEnv(key)

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("environment variable is not a valid integer. Value: %s", val)
	}

	return intVal
}

func reqEnvBool(key string) bool {
	val := reqEnv(key)

	return val == "true"
}
