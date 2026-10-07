package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestRequireJWT(t *testing.T) {
	secret := []byte("test-secret-with-at-least-32-characters")
	id := uuid.New()
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("Test", "1"))
	huma.Register(api, huma.Operation{
		OperationID: "protected", Method: http.MethodGet, Path: "/protected",
		Middlewares: huma.Middlewares{RequireJWT(api, secret)},
	}, func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		if UserID(ctx) != id {
			t.Error("authenticated user was not passed to the handler")
		}
		return &struct{}{}, nil
	})
	for _, tc := range []struct {
		name   string
		key    []byte
		method jwt.SigningMethod
		exp    *time.Time
		status int
	}{
		{"valid", secret, jwt.SigningMethodHS256, timePtr(time.Now().Add(time.Hour)), 204},
		{"wrong secret", []byte("wrong-key"), jwt.SigningMethodHS256, timePtr(time.Now().Add(time.Hour)), 401},
		{"expired", secret, jwt.SigningMethodHS256, timePtr(time.Now().Add(-time.Hour)), 401},
		{"missing expiration", secret, jwt.SigningMethodHS256, nil, 401},
		{"wrong algorithm", secret, jwt.SigningMethodHS512, timePtr(time.Now().Add(time.Hour)), 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			claims := jwt.MapClaims{"sub": id.String()}
			if tc.exp != nil {
				claims["exp"] = tc.exp.Unix()
			}
			token, err := jwt.NewWithClaims(tc.method, claims).SignedString(tc.key)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			if resp.Code != tc.status {
				t.Fatalf("got %d, want %d", resp.Code, tc.status)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time { return &t }
