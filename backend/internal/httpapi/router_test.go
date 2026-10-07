package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func TestHealthAndReadiness(t *testing.T) {
	router := http.NewServeMux()
	api := humago.NewWithPrefix(router, "/api", huma.DefaultConfig("Test", "1"))
	calls := 0
	databaseErr := error(nil)
	registerHealth(api, func(ctx context.Context) error {
		calls++
		if _, ok := ctx.Deadline(); !ok {
			t.Error("readiness probe has no deadline")
		}
		return databaseErr
	})
	check := func(path string, status int) map[string]any {
		t.Helper()
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, path, nil))
		if resp.Code != status {
			t.Fatalf("%s: got %d, want %d", path, resp.Code, status)
		}
		var body map[string]any
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		return body
	}
	if check("/api/health", 200)["status"] != "ok" || calls != 0 {
		t.Fatal("liveness must work without querying PostgreSQL")
	}
	if check("/api/ready", 200)["database"] != "connected" || calls != 1 {
		t.Fatal("readiness did not confirm PostgreSQL connectivity")
	}
	databaseErr = errors.New("private connection information")
	if check("/api/ready", 503)["detail"] != "database unavailable" {
		t.Fatal("readiness leaked database error details")
	}
	check("/api/health", 200)
}

func TestInvalidInputsAndPrivateDirectory(t *testing.T) {
	// No database is attached: any invalid input reaching a repository panics.
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"POST", "/api/auth/register", `{"name":"Name","email":"a@example.com","password":"password","is_admin":true}`, 422},
		{"POST", "/api/auth/register", `{"name":"<script>","email":"a@example.com","password":"password"}`, 422},
		{"POST", "/api/auth/register", `{"name":"Name","email":"a@example.com","password":"short"}`, 422},
		{"POST", "/api/auth/register", `{"name":null,"email":"a@example.com","password":"password"}`, 422},
		{"POST", "/api/auth/register", `{"name":12,"email":"a@example.com","password":"password"}`, 422},
		{"POST", "/api/auth/login", `{"email":"invalid","password":"password"}`, 401},
		{"POST", "/api/auth/login", `{"email":"a@example.com","password":"` + strings.Repeat("x", 73) + `"}`, 422},
		{"POST", "/api/auth/register?is_admin=true", `{"name":"Name","email":"a@example.com","password":"password"}`, 422},
		{"GET", "/api/users", "", 404},
		{"GET", "/api/users/00000000-0000-0000-0000-000000000001", "", 404},
		{"GET", "/api/admin/users", "", 401},
		{"GET", "/api/users/me", "", 401},
	} {
		t.Run(tc.method+tc.path+tc.body, func(t *testing.T) {
			handler := New(nil, nil, []byte(strings.Repeat("s", 32)), "")
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("got %d, want %d: %s", rec.Code, tc.status, rec.Body.String())
			}
		})
	}
}
