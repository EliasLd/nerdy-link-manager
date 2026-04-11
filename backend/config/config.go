package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	JWTSecret            string
	DBPath               string
	InitialAdminEmail    string
	InitialAdminPassword string
}

// Reads environment file and returns config
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found, using environment variables.")
	} else {
		log.Println("[INFO] Found .env file, using it.")
	}

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		DBPath:               getEnv("DB_PATH", "internal/db/nedry_link_manaer.db"),
		InitialAdminEmail:    getEnv("INITIAL_ADMIN_EMAIL", ""),
		InitialAdminPassword: getEnv("INITIAL_ADMIN_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
