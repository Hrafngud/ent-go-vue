package config

import (
	"net/url"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	for _, key := range []string{"DATABASE_URL", "PORT", "JWT_SECRET", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "ROOT_EMAIL", "ROOT_PASSWORD"} {
		t.Setenv(key, "")
	}
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
	t.Setenv("DB_USER", "user @test")
	t.Setenv("DB_PASSWORD", "p@ss word:/?#")
	t.Setenv("DB_NAME", "app")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	password, _ := u.User.Password()
	if password != "p@ss word:/?#" || u.User.Username() != "user @test" || u.Host != "localhost:5432" || cfg.Port != "8080" {
		t.Fatal("legacy configuration did not preserve credentials or defaults")
	}
	t.Setenv("DATABASE_URL", "postgres://override:secret@postgres:5432/override?sslmode=disable")
	cfg, err = Load()
	if err != nil || cfg.DatabaseURL != "postgres://override:secret@postgres:5432/override?sslmode=disable" {
		t.Fatal("DATABASE_URL did not take precedence")
	}
	t.Setenv("PORT", "70000")
	if _, err := Load(); err == nil {
		t.Fatal("invalid port accepted")
	}
	t.Setenv("PORT", "8080")
	t.Setenv("JWT_SECRET", "replace_with_at_least_32_random_characters")
	if _, err := Load(); err == nil {
		t.Fatal("example secret accepted")
	}
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
	t.Setenv("DATABASE_URL", "invalid://user:private-password@localhost/app")
	if _, err := Load(); err == nil || strings.Contains(err.Error(), "private-password") {
		t.Fatal("invalid database URL accepted or leaked credentials")
	}
}

func TestRootCredentials(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:secret@localhost/app")
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
	t.Setenv("PORT", "8080")
	for _, tc := range []struct {
		name, email, password string
		valid                 bool
	}{
		{"disabled", "", "", true},
		{"configured", "root@example.com", "root password", true},
		{"trim email", " root@example.com ", "root password", true},
		{"email only", "root@example.com", "", false},
		{"password only", "", "root password", false},
		{"invalid email", "root", "root password", false},
		{"display name", "Root <root@example.com>", "root password", false},
		{"short password", "root@example.com", "short", false},
		{"long password", "root@example.com", strings.Repeat("a", 73), false},
		{"blank password", "root@example.com", strings.Repeat(" ", 8), false},
		{"placeholder", "root@example.com", "replace_with_root_password", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ROOT_EMAIL", tc.email)
			t.Setenv("ROOT_PASSWORD", tc.password)
			cfg, err := Load()
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v, error=%v", tc.valid, err)
			}
			if tc.valid && (cfg.RootEmail != strings.TrimSpace(tc.email) || cfg.RootPassword != tc.password) {
				t.Fatal("root credentials were not preserved")
			}
			if err != nil && strings.Contains(err.Error(), tc.password) && tc.password != "" {
				t.Fatal("validation error leaked password")
			}
		})
	}
}
