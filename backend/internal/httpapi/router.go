package httpapi

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"backend-golang/ent"
	authhttp "backend-golang/internal/auth/delivery/http"
	authuc "backend-golang/internal/auth/usecase"
	userhttp "backend-golang/internal/user/delivery/http"
	userrepo "backend-golang/internal/user/repository"
	useruc "backend-golang/internal/user/usecase"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// New wires the existing feature modules into the public /api namespace.
func New(client *ent.Client, db *sql.DB, secret []byte) http.Handler {
	router := http.NewServeMux()
	config := huma.DefaultConfig("Ent Go Vue API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
	}
	api := humago.NewWithPrefix(router, "/api", config)
	repo := userrepo.NewEntRepository(client)
	authhttp.RegisterRoutes(api, authuc.NewUsecase(repo, secret))
	userhttp.RegisterMyUserRoutes(api, useruc.NewMyUsecase(repo), secret)
	userhttp.RegisterPublicUserRoutes(api, useruc.NewPublicUsecase(repo))
	registerHealth(api, db.PingContext)
	return router
}

type healthResponse struct {
	Body struct {
		Status string `json:"status"`
	}
}

type readyResponse struct {
	Body struct {
		Status   string `json:"status"`
		Database string `json:"database"`
	}
}

func registerHealth(api huma.API, ping func(context.Context) error) {
	huma.Register(api, huma.Operation{
		OperationID: "health", Method: http.MethodGet, Path: "/health",
		Summary: "API liveness (no database work)", Tags: []string{"Health"},
	}, func(context.Context, *struct{}) (*healthResponse, error) {
		resp := &healthResponse{}
		resp.Body.Status = "ok"
		return resp, nil
	})
	huma.Register(api, huma.Operation{
		OperationID: "ready", Method: http.MethodGet, Path: "/ready",
		Summary: "API readiness and PostgreSQL connectivity", Tags: []string{"Health"},
		Errors: []int{http.StatusServiceUnavailable},
	}, func(ctx context.Context, _ *struct{}) (*readyResponse, error) {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := ping(ctx); err != nil {
			return nil, huma.Error503ServiceUnavailable("database unavailable")
		}
		resp := &readyResponse{}
		resp.Body.Status, resp.Body.Database = "ok", "connected"
		return resp, nil
	})
}
