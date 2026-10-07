package http

import (
	"context"
	"errors"
	"net/http"

	"backend-golang/internal/auth/middleware"
	"backend-golang/internal/user"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreateUserRequest struct{ Body user.UserInput }
type UpdateUserRequest struct {
	ID   uuid.UUID `path:"id"`
	Body user.UserInput
}

func RegisterAdminUserRoutes(api huma.API, uc user.AdminUsecase, secret []byte) {
	requireAdmin := func(ctx huma.Context, next func(huma.Context)) {
		allowed, err := uc.CanManage(ctx.Context(), middleware.UserID(ctx.Context()))
		if errors.Is(err, user.ErrNotFound) {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "account no longer exists")
			return
		}
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "could not verify administrator access")
			return
		}
		if !allowed {
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "administrator access required")
			return
		}
		next(ctx)
	}
	operation := func(id, method, path, summary string) huma.Operation {
		return huma.Operation{
			OperationID: id, Method: method, Path: path, Summary: summary, Tags: []string{"Admin users"},
			Security:    []map[string][]string{{"bearerAuth": {}}},
			Middlewares: huma.Middlewares{middleware.RequireJWT(api, secret), requireAdmin},
			Errors:      []int{400, 401, 403, 404, 409, 422, 500},
		}
	}
	huma.Register(api, operation("admin-list-users", http.MethodGet, "/admin/users", "List users as administrator"),
		func(ctx context.Context, _ *struct{}) (*ListUsersResponse, error) {
			users, err := uc.ListUsers(ctx)
			if err != nil {
				return nil, adminError(err)
			}
			resp := &ListUsersResponse{}
			resp.Body.Data = users
			if users == nil {
				resp.Body.Data = []*user.User{}
			}
			return resp, nil
		})
	huma.Register(api, operation("admin-get-user", http.MethodGet, "/admin/users/{id}", "Get user as administrator"),
		func(ctx context.Context, input *GetUserRequest) (*UserResponse, error) {
			usr, err := uc.GetUser(ctx, input.ID)
			if err != nil {
				return nil, adminError(err)
			}
			resp := &UserResponse{}
			resp.Body.Data = usr
			return resp, nil
		})
	createOperation := operation("admin-create-user", http.MethodPost, "/admin/users", "Create a user")
	createOperation.DefaultStatus = http.StatusCreated
	huma.Register(api, createOperation, func(ctx context.Context, input *CreateUserRequest) (*UserResponse, error) {
		usr, err := uc.CreateUser(ctx, input.Body)
		if err != nil {
			return nil, adminError(err)
		}
		resp := &UserResponse{}
		resp.Body.Data = usr
		return resp, nil
	})
	huma.Register(api, operation("admin-update-user", http.MethodPut, "/admin/users/{id}", "Update a user; omit password to keep it"),
		func(ctx context.Context, input *UpdateUserRequest) (*UserResponse, error) {
			usr, err := uc.UpdateUser(ctx, input.ID, input.Body)
			if err != nil {
				return nil, adminError(err)
			}
			resp := &UserResponse{}
			resp.Body.Data = usr
			return resp, nil
		})
	deleteOperation := operation("admin-delete-user", http.MethodDelete, "/admin/users/{id}", "Delete a user")
	deleteOperation.DefaultStatus = http.StatusNoContent
	huma.Register(api, deleteOperation, func(ctx context.Context, input *GetUserRequest) (*struct{}, error) {
		if err := uc.DeleteUser(ctx, input.ID); err != nil {
			return nil, adminError(err)
		}
		return nil, nil
	})
}

func adminError(err error) error {
	switch {
	case errors.Is(err, user.ErrNotFound):
		return huma.Error404NotFound("user not found")
	case errors.Is(err, user.ErrEmailTaken):
		return huma.Error409Conflict("email is already in use")
	case errors.Is(err, user.ErrRootProtected):
		return huma.Error403Forbidden(user.ErrRootProtected.Error())
	case errors.Is(err, user.ErrInvalidInput):
		return huma.Error422UnprocessableEntity(user.ErrInvalidInput.Error())
	default:
		return huma.Error500InternalServerError("could not save user changes")
	}
}
