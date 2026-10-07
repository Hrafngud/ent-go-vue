package http

import (
	"context"
	"errors"
	"net/http"

	"backend-golang/internal/auth/middleware"
	"backend-golang/internal/user"

	"github.com/danielgtaylor/huma/v2"
)

type MyUserRequest struct {
	Authorization string `header:"Authorization" doc:"Bearer token"`
}

type MyUserResponse struct {
	Body struct {
		Data *user.User `json:"data"`
	}
}

func RegisterMyUserRoutes(api huma.API, uc user.MyUsecase, secret []byte) {
	huma.Register(api, huma.Operation{
		OperationID: "get-my-profile",
		Method:      http.MethodGet,
		Path:        "/users/me",
		Summary:     "Get current user profile",
		Tags:        []string{"Users"},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Middlewares: huma.Middlewares{middleware.RequireJWT(api, secret)},
	}, func(ctx context.Context, input *MyUserRequest) (*MyUserResponse, error) {
		usr, err := uc.GetProfile(ctx, middleware.UserID(ctx))
		if err != nil {
			if errors.Is(err, user.ErrNotFound) {
				return nil, huma.Error404NotFound("user not found")
			}
			return nil, huma.Error500InternalServerError("could not load profile")
		}
		resp := &MyUserResponse{}
		resp.Body.Data = usr
		return resp, nil
	})
}
