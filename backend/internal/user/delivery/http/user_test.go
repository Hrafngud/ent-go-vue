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
	ts := httptest.NewServer(httpapi.New(client, db, secret, "root@example.com"))
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

	t.Run("Root-only user management CRUD", func(t *testing.T) {
		request := func(method, path, token string, body any, want int) []byte {
			t.Helper()
			payload, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(method, ts.URL+"/api"+path, bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			var result bytes.Buffer
			if _, err := result.ReadFrom(resp.Body); err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != want {
				t.Fatalf("%s %s: got %d, want %d: %s", method, path, resp.StatusCode, want, result.String())
			}
			return result.Bytes()
		}
		var login struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(request(http.MethodPost, "/auth/login", "", map[string]string{"email": "root@example.com", "password": "rotated-root-password"}, 200), &login); err != nil {
			t.Fatal(err)
		}
		var profile struct {
			Data struct {
				ID      string `json:"id"`
				IsAdmin bool   `json:"is_admin"`
			} `json:"data"`
		}
		if err := json.Unmarshal(request(http.MethodGet, "/users/me", login.Token, nil, 200), &profile); err != nil {
			t.Fatal(err)
		}
		if !profile.Data.IsAdmin {
			t.Fatal("root profile does not identify administrator access")
		}
		memberID := userID.String()
		for _, op := range []struct{ method, path string }{
			{http.MethodGet, "/admin/users"},
			{http.MethodGet, "/admin/users/" + memberID},
			{http.MethodPost, "/admin/users"},
			{http.MethodPut, "/admin/users/" + memberID},
			{http.MethodDelete, "/admin/users/" + memberID},
		} {
			request(op.method, op.path, "", nil, 401)
			request(op.method, op.path, jwtToken, nil, 403)
		}
		request(http.MethodPost, "/admin/users", login.Token, map[string]string{"name": "Duplicate", "email": "auth@test.com", "password": "password123"}, 409)
		request(http.MethodPost, "/auth/register", "", map[string]string{"name": "Duplicate", "email": "auth@test.com", "password": "password123"}, 409)
		for _, body := range []map[string]string{
			{"name": " ", "email": "valid@test.com", "password": "password123"},
			{"name": "Invalid", "email": "invalid", "password": "password123"},
			{"name": "Invalid", "email": "valid@test.com", "password": "short"},
			{"name": "Invalid", "email": "valid@test.com", "password": "😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀😀"},
		} {
			request(http.MethodPost, "/admin/users", login.Token, body, 422)
		}

		var created struct {
			Data struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"data"`
		}
		body := request(http.MethodPost, "/admin/users", login.Token, map[string]string{"name": " Managed User ", "email": " managed@test.com ", "password": "initial-password"}, 201)
		if bytes.Contains(body, []byte("password")) {
			t.Fatal("password leaked in creation response")
		}
		if err := json.Unmarshal(body, &created); err != nil {
			t.Fatal(err)
		}
		if created.Data.Name != "Managed User" || created.Data.Email != "managed@test.com" {
			t.Fatal("input was not normalized")
		}
		id := created.Data.ID
		request(http.MethodGet, "/admin/users/"+id, login.Token, nil, 200)
		request(http.MethodGet, "/admin/users", login.Token, nil, 200)
		request(http.MethodPut, "/admin/users/"+id, login.Token, map[string]string{"name": "Renamed User", "email": "renamed@test.com"}, 200)
		request(http.MethodPost, "/auth/login", "", map[string]string{"email": "renamed@test.com", "password": "initial-password"}, 200)
		request(http.MethodPut, "/admin/users/"+id, login.Token, map[string]string{"name": "Renamed User", "email": "auth@test.com"}, 409)
		request(http.MethodPut, "/admin/users/"+id, login.Token, map[string]string{"name": "Renamed User", "email": "renamed@test.com", "password": "rotated-password"}, 200)
		request(http.MethodPost, "/auth/login", "", map[string]string{"email": "renamed@test.com", "password": "initial-password"}, 401)
		request(http.MethodPost, "/auth/login", "", map[string]string{"email": "renamed@test.com", "password": "rotated-password"}, 200)
		stored, err := client.User.Get(ctx, uuid.MustParse(id))
		if err != nil || bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte("rotated-password")) != nil {
			t.Fatal("password was not hashed")
		}
		request(http.MethodDelete, "/admin/users/"+profile.Data.ID, login.Token, nil, 403)
		request(http.MethodPut, "/admin/users/"+profile.Data.ID, login.Token, map[string]string{"name": "Root", "email": "other-root@test.com"}, 403)
		request(http.MethodDelete, "/admin/users/"+id, login.Token, nil, 204)
		request(http.MethodGet, "/admin/users/"+id, login.Token, nil, 404)
		request(http.MethodPut, "/admin/users/"+id, login.Token, map[string]string{"name": "Deleted", "email": "deleted@test.com"}, 404)
		request(http.MethodDelete, "/admin/users/"+id, login.Token, nil, 404)

		// Disabling ROOT_EMAIL must deny even a valid root session.
		disabled := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+login.Token)
		httpapi.New(client, db, secret, "").ServeHTTP(disabled, req)
		if disabled.Code != http.StatusForbidden {
			t.Fatalf("disabled root policy returned %d", disabled.Code)
		}
	})
}
