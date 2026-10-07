package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   []byte
}

func Load() (Config, error) {
	cfg := Config{DatabaseURL: os.Getenv("DATABASE_URL"), Port: value("PORT", "8080"), JWTSecret: []byte(os.Getenv("JWT_SECRET"))}
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, errors.New("PORT must be an integer between 1 and 65535")
	}
	if len(cfg.JWTSecret) < 32 || strings.HasPrefix(string(cfg.JWTSecret), "replace_with_") {
		return Config{}, errors.New("JWT_SECRET must contain at least 32 characters and must not be an example placeholder")
	}
	if cfg.DatabaseURL == "" {
		if os.Getenv("DB_USER") == "" || os.Getenv("DB_NAME") == "" || os.Getenv("DB_PASSWORD") == "" {
			return Config{}, errors.New("set DATABASE_URL or DB_USER, DB_PASSWORD, and DB_NAME")
		}
		// URL encoding preserves spaces and reserved characters in credentials.
		u := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
			Host:   net.JoinHostPort(value("DB_HOST", "localhost"), value("DB_PORT", "5432")),
			Path:   "/" + os.Getenv("DB_NAME"),
		}
		q := url.Values{"sslmode": {value("DB_SSLMODE", "disable")}}
		u.RawQuery = q.Encode()
		cfg.DatabaseURL = u.String()
	}
	u, err := url.Parse(cfg.DatabaseURL)
	if err != nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Host == "" || u.Path == "" || u.Path == "/" {
		return Config{}, errors.New("DATABASE_URL must be a PostgreSQL URL with a host and database name")
	}
	return cfg, nil
}

func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
