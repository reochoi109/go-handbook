package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
	Port  int
}

func main() {
	godotenv.Load()
	cfg := Config{
		DBURL: getEnv("DB_URL", "postgres://localhost:5432/db"),
		Port:  parseEnvInt("PORT", 8080),
	}

	fmt.Printf("Config loaded: %+v\n", cfg)
}

// default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// type change
func parseEnvInt(key string, fallback int) int {
	val := getEnv(key, "")
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}
