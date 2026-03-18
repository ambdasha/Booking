package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTP struct {
		Addr string
	}
	DB struct {
		DSN string
	}
	Auth struct {
		JWTSecret string
	}
	Log struct {
		Level string
	}
}

func MustLoad() Config {
	var c Config

	c.HTTP.Addr = getEnv("HTTP_ADDR", ":8080")
	c.DB.DSN = getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/booking?sslmode=disable")
	c.Auth.JWTSecret = getEnv("JWT_SECRET", "dev_secret_change_me")
	c.Log.Level = getEnv("LOG_LEVEL", "info")

	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		log.Fatal("JWT_SECRET is required")
	}
	return c
}

func (c Config) LogLevel() slog.Level {
	switch strings.ToLower(c.Log.Level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getEnv(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}