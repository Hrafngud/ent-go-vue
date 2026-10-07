package http_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"backend-golang/ent"
	entuser "backend-golang/ent/user"
	"backend-golang/internal/httpapi"
	userrepo "backend-golang/internal/user/repository"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

func TestUserAndAuthAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test")
	}

	ctx := context.Background()

	migrations, err := filepath.Glob("../../../../ent/migrate/migrations/*.sql")
	if err != nil || len(migrations) == 0 {
		t.Fatal("migration files unavailable")
	}
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("test-user"),
		postgres.WithPassword("test-pass"),
		postgres.WithInitScripts(migrations...),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	secret := []byte("integration-test-secret-at-least-32-characters")
	ts := httptest.NewServer(httpapi.New(client, db, secret))
	defer ts.Close()

	t.Run("Health and database readiness", func(t *testing.T) {
		for _, path := range []string{"/api/health", "/api/ready", "/api/openapi.json", "/api/users"} {
			resp, err := http.Get(ts.URL + path)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("%s: expected 200, got %d", path, resp.StatusCode)
			}
		}
	})

	t.Run("Root provisioning and password rotation", func(t *testing.T) {
		const email = "root@example.com"
		const password = "initial-root-password"
		if err := userrepo.SyncRootUser(ctx, client, "", ""); err != nil {
			t.Fatal(err)
		}
		count, err := client.User.Query().Count(ctx)
		if err != nil || count != 0 {
			t.Fatal("disabled provisioning created an account")
		}
		if err := userrepo.SyncRootUser(ctx, client, email, password); err != nil {
			t.Fatal(err)
		}
		initial, err := client.User.Query().Where(entuser.Email(email)).Only(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if initial.Name != "Root" || initial.Password == password || bcrypt.CompareHashAndPassword([]byte(initial.Password), []byte(password)) != nil {
			t.Fatal("root account was not created with a hashed password")
		}
		if err := userrepo.SyncRootUser(ctx, client, email, password); err != nil {
			t.Fatal(err)
		}
		unchanged, err := client.User.Get(ctx, initial.ID)
		if err != nil || unchanged.Password != initial.Password {
			t.Fatal("unchanged credentials rewrote the password")
		}
		const rotated = "rotated-root-password"
		if err := userrepo.SyncRootUser(ctx, client, email, rotated); err != nil {
			t.Fatal(err)
		}
		current, err := client.User.Query().Where(entuser.Email(email)).Only(ctx)
		if err != nil || current.ID != initial.ID || !current.CreatedAt.Equal(initial.CreatedAt) {
			t.Fatal("password rotation replaced the root account")
		}
		for _, credentials := range []struct {
			password string
			status   int
		}{{password, http.StatusUnauthorized}, {rotated, http.StatusOK}} {
			body, err := json.Marshal(map[string]string{"email": email, "password": credentials.password})
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.Post(ts.URL+"/api/auth/login", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()
			if resp.StatusCode != credentials.status {
				t.Fatalf("root login: got %d, want %d", resp.StatusCode, credentials.status)
			}
		}
		count, err = client.User.Query().Count(ctx)
		if err != nil || count != 1 {
			t.Fatal("repeated provisioning duplicated the root account")
		}
	})

	var jwtToken string

	t.Run("Register User", func(t *testing.T) {
		reqBody := []byte(`{"name":"Auth User","email":"auth@test.com","password":"securepassword123"}`)
		resp, err := http.Post(ts.URL+"/api/auth/register", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status OK/Created, got %d", resp.StatusCode)
		}
	})

	t.Run("Login User", func(t *testing.T) {
		reqBody := []byte(`{"email":"auth@test.com","password":"securepassword123"}`)
		resp, err := http.Post(ts.URL+"/api/auth/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}

		var actualResp struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&actualResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if actualResp.Token == "" {
			t.Fatalf("expected token, got empty string")
		}
		jwtToken = actualResp.Token
	})

	var userID uuid.UUID

	t.Run("Get My Profile (Protected)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}

		var fetchResp struct {
			Data struct {
				ID    uuid.UUID `json:"id"`
				Email string    `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&fetchResp); err != nil {
			t.Fatalf("failed to decode get response: %v", err)
		}
		if fetchResp.Data.Email != "auth@test.com" {
			t.Errorf("expected email auth@test.com, got %s", fetchResp.Data.Email)
		}
		userID = fetchResp.Data.ID
	})

	t.Run("Get My Profile (Unauthorized)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/users/me", nil)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Get User Detail (Public)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/users/" + userID.String())
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}

		var detailResp struct {
			Data struct {
				ID    uuid.UUID `json:"id"`
				Email string    `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
			t.Fatalf("failed to decode get response: %v", err)
		}
		if detailResp.Data.Email != "auth@test.com" {
			t.Errorf("expected email auth@test.com, got %s", detailResp.Data.Email)
		}
	})

	t.Run("List Users (Public)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/users")
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}
	})
}
