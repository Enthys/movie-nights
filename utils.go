package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"movie_night/types"
	"movie_night/validator"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type envelope map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func readJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	var maxBytes int64 = 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func internalErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Something went wrong!"))
}

func badRequestErrorResponse(w http.ResponseWriter, err envelope) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(map[string]envelope{"error": err}); err != nil {
		log.Println(err)
	}
}

func notFoundResponse(w http.ResponseWriter, err envelope) {
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(map[string]envelope{"error": err}); err != nil {
		log.Println(err)
	}
}

func validationErrorResponse(w http.ResponseWriter, v *validator.Validator) {
	w.WriteHeader(http.StatusBadRequest)
	m := map[string]envelope{
		"errors": {},
	}

	for key, err := range v.Errors {
		m["errors"][key] = err
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
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

func queryVal(r *http.Request, key string, defaultVal string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}

	return val
}

func queryIntVal(r *http.Request, key string, defaultVal int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal, nil
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal, err
	}

	return intVal, nil
}

func pathArgIntVal(r *http.Request, key string) (int, error) {
	val := r.PathValue(key)
	if val == "" {
		return 0, fmt.Errorf("missing path variable")
	}

	intVal, err := strconv.Atoi(val)

	return intVal, err
}
