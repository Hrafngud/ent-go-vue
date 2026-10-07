package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
