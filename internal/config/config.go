package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT              string
	DB_URL            string
	REDIS_URL         string
	RABBITMQ_URL      string
	JWT_SECRET        string
	RESEND_API_KEY    string
	WORKER_COUNT      int
	FRONTEND_URL      string
	RESEND_EMAIL_FROM string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using system env")
	}

	workerCount, _ := strconv.Atoi(getEnv("WORKER_COUNT", "10"))
	return &Config{
		PORT:              getEnv("Port", ""),
		DB_URL:            getEnv("GOOSE_DBSTRING", ""),
		REDIS_URL:         getEnv("REDIS_URL", ""),
		RABBITMQ_URL:      getEnv("RABBITMQ_URL", ""),
		JWT_SECRET:        getEnv("JWT_SECRET", ""),
		RESEND_API_KEY:    getEnv("RESEND_API_KEY", ""),
		WORKER_COUNT:      workerCount,
		FRONTEND_URL:      getEnv("FRONTEND_URL", ""),
		RESEND_EMAIL_FROM: getEnv("RESEND_EMAIL_FROM", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
