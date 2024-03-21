package main

import (
	"encoding/json"
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

func badRequestErrorResponse(w http.ResponseWriter, err map[string]string) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		log.Println(err)
	}
}

func conflictErrorResponse(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusConflict)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": reason}); err != nil {
		log.Println()
	}
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
