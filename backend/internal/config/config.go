package config

import (
	"errors"
	"net"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"strings"

	"backend-golang/internal/user"
)

type Config struct {
	DatabaseURL       string
	Port              string
	JWTSecret         []byte
	RootEmail         string
	RootPassword      string
	TrustedProxyCIDRs []netip.Prefix
}

func Load() (Config, error) {
	cfg := Config{DatabaseURL: os.Getenv("DATABASE_URL"), Port: value("PORT", "8080"), JWTSecret: []byte(os.Getenv("JWT_SECRET"))}
	cfg.RootEmail = strings.TrimSpace(os.Getenv("ROOT_EMAIL"))
	cfg.RootPassword = os.Getenv("ROOT_PASSWORD")
	var err error
	cfg.TrustedProxyCIDRs, err = trustedProxies(os.Getenv("TRUSTED_PROXY_CIDRS"))
	if err != nil {
		return Config{}, err
	}
	if err := validateRoot(cfg.RootEmail, cfg.RootPassword); err != nil {
		return Config{}, err
	}
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, errors.New("PORT must be an integer between 1 and 65535")
	}
	if err := validateJWTSecret(cfg.JWTSecret); err != nil {
		return Config{}, err
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL, err = legacyDatabaseURL()
		if err != nil {
			return Config{}, err
		}
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

func validateJWTSecret(secret []byte) error {
	if len(secret) < 32 || strings.HasPrefix(string(secret), "replace_with_") {
		return errors.New("JWT_SECRET must contain at least 32 characters and must not be an example placeholder")
	}
	return nil
}

func trustedProxies(value string) ([]netip.Prefix, error) {
	var prefixes []netip.Prefix
	for _, cidr := range strings.Split(value, ",") {
		if strings.TrimSpace(cidr) == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(strings.TrimSpace(cidr))
		if err != nil || prefix.Bits() == 0 {
			return nil, errors.New("TRUSTED_PROXY_CIDRS must contain specific IPv4/IPv6 CIDRs, never a catch-all network")
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

func validateRoot(email, password string) error {
	if email == "" && password == "" {
		return nil
	}
	if _, err := user.NormalizeEmail(email); err != nil {
		return errors.New("ROOT_EMAIL must be a valid email address when root credentials are configured")
	}
	if !user.ValidPassword(password, 8) || strings.HasPrefix(password, "replace_with_") {
		return errors.New("ROOT_PASSWORD must contain 8 to 72 bytes and must not be an example placeholder")
	}
	return nil
}

func legacyDatabaseURL() (string, error) {
	if os.Getenv("DB_USER") == "" || os.Getenv("DB_NAME") == "" || os.Getenv("DB_PASSWORD") == "" {
		return "", errors.New("set DATABASE_URL or DB_USER, DB_PASSWORD, and DB_NAME")
	}
	// URL encoding preserves spaces and reserved characters in credentials.
	u := url.URL{Scheme: "postgres", User: url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")), Host: net.JoinHostPort(value("DB_HOST", "localhost"), value("DB_PORT", "5432")), Path: "/" + os.Getenv("DB_NAME")}
	u.RawQuery = url.Values{"sslmode": {value("DB_SSLMODE", "disable")}}.Encode()
	return u.String(), nil
}
