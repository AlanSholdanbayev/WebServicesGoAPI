package config

import (
	"fmt"
	"os"
)

type LoggerConfig struct {
	SeqURL    string
	SeqAPIKey string
}

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
	Logger    LoggerConfig
}

// Load загружает конфигурацию из переменных окружения или задаёт значения по умолчанию
func Load() (*Config, error) {
	port := getenv("PORT", "8080")

	dbHost := getenv("DB_HOST", "localhost")
	dbPort := getenv("DB_PORT", "5432")
	dbUser := getenv("DB_USER", "postgres")
	dbPass := getenv("DB_PASS", "postgres")
	dbName := getenv("DB_NAME", "final_project")

	jwtSecret := getenv("JWT_SECRET", "secret")

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	seqURL := getenv("SEQ_URL", "http://localhost:5341") // локальный Seq
	seqAPIKey := getenv("SEQ_API_KEY", "")               // если ключ не используется

	cfg := &Config{
		Port:      port,
		DBUrl:     dbUrl,
		JWTSecret: jwtSecret,
		Logger: LoggerConfig{
			SeqURL:    seqURL,
			SeqAPIKey: seqAPIKey,
		},
	}

	return cfg, nil
}

// getenv берёт переменную окружения или возвращает значение по умолчанию
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
