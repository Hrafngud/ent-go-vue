package usecase

import (
	"context"

	"backend-golang/internal/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type adminUsecase struct {
	repo      user.Repository
	rootEmail string
}

func NewAdminUsecase(repo user.Repository, rootEmail string) user.AdminUsecase {
	return &adminUsecase{repo: repo, rootEmail: rootEmail}
}

func (u *adminUsecase) isRoot(usr *user.User) bool {
	return u.rootEmail != "" && usr.Email == u.rootEmail
}

func (u *adminUsecase) CanManage(ctx context.Context, id uuid.UUID) (bool, error) {
	usr, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	return u.isRoot(usr), nil
}

func (u *adminUsecase) ListUsers(ctx context.Context) ([]*user.User, error) {
	users, err := u.repo.List(ctx)
	for _, usr := range users {
		usr.IsAdmin = u.isRoot(usr)
	}
	return users, err
}

func (u *adminUsecase) GetUser(ctx context.Context, id uuid.UUID) (*user.User, error) {
	usr, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	usr.IsAdmin = u.isRoot(usr)
	return usr, nil
}

func (u *adminUsecase) CreateUser(ctx context.Context, input user.UserInput) (*user.User, error) {
	input, err := user.NormalizeInput(input, true)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return u.repo.Create(ctx, &user.User{Name: input.Name, Email: input.Email, Password: string(hash)})
}

func (u *adminUsecase) UpdateUser(ctx context.Context, id uuid.UUID, input user.UserInput) (*user.User, error) {
	input, err := user.NormalizeInput(input, false)
	if err != nil {
		return nil, err
	}
	current, err := u.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.isRoot(current) && input.Email != current.Email {
		return nil, user.ErrRootProtected
	}
	password := ""
	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(hash)
	}
	updated, err := u.repo.Update(ctx, id, &user.User{Name: input.Name, Email: input.Email, Password: password})
	if err != nil {
		return nil, err
	}
	updated.IsAdmin = u.isRoot(updated)
	return updated, nil
}

func (u *adminUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	current, err := u.GetUser(ctx, id)
	if err != nil {
		return err
	}
	if u.isRoot(current) {
		return user.ErrRootProtected
	}
	return u.repo.Delete(ctx, id)
}
