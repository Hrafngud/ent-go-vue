package http

import (
	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type GetUserRequest struct {
	ID uuid.UUID `path:"id" doc:"User ID"`
}

type UserResponse struct {
	Body struct {
		Data *user.User `json:"data"`
	}
}

type ListUsersResponse struct {
	Body struct {
		Data []*user.User `json:"data"`
	}
}
