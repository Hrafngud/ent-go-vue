package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
}

var (
	ErrNotFound      = errors.New("user not found")
	ErrEmailTaken    = errors.New("email is already in use")
	ErrRootProtected = errors.New("the configured root account cannot be deleted or change its email")
	ErrInvalidInput  = errors.New("provide a plain-text name up to 255 bytes, valid email, and a password between 8 and 72 bytes")
)

type UserInput struct {
	Name     string `json:"name" minLength:"1" maxLength:"255"`
	Email    string `json:"email" maxLength:"254"`
	Password string `json:"password,omitempty" maxLength:"72"`
}

type Repository interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, u *User) (*User, error)
}

type AdminUsecase interface {
	CanManage(ctx context.Context, id uuid.UUID) (bool, error)
	ListUsers(ctx context.Context) ([]*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	CreateUser(ctx context.Context, input UserInput) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input UserInput) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type MyUsecase interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*User, error)
}
