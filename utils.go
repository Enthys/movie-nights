package main

import (
	"log"
	"os"
	"strconv"
)

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
