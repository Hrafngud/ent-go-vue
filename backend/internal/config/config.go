package config

import (
	"errors"
	"net"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL  string
	Port         string
	JWTSecret    []byte
	RootEmail    string
	RootPassword string
}

func Load() (Config, error) {
	cfg := Config{DatabaseURL: os.Getenv("DATABASE_URL"), Port: value("PORT", "8080"), JWTSecret: []byte(os.Getenv("JWT_SECRET"))}
	cfg.RootEmail = strings.TrimSpace(os.Getenv("ROOT_EMAIL"))
	cfg.RootPassword = os.Getenv("ROOT_PASSWORD")
	if cfg.RootEmail != "" || cfg.RootPassword != "" {
		address, err := mail.ParseAddress(cfg.RootEmail)
		if err != nil || address.Address != cfg.RootEmail {
			return Config{}, errors.New("ROOT_EMAIL must be a valid email address when root credentials are configured")
		}
		if len(cfg.RootPassword) < 8 || len(cfg.RootPassword) > 72 || strings.TrimSpace(cfg.RootPassword) == "" || strings.HasPrefix(cfg.RootPassword, "replace_with_") {
			return Config{}, errors.New("ROOT_PASSWORD must contain 8 to 72 bytes and must not be an example placeholder")
		}
	}
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
