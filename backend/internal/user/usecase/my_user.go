package usecase

import (
	"context"

	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type myUsecaseImpl struct {
	repo      user.Repository
	rootEmail string
}

func NewMyUsecase(repo user.Repository, rootEmail string) user.MyUsecase {
	return &myUsecaseImpl{repo: repo, rootEmail: rootEmail}
}

func (u *myUsecaseImpl) GetProfile(ctx context.Context, id uuid.UUID) (*user.User, error) {
	usr, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	usr.IsAdmin = u.rootEmail != "" && usr.Email == u.rootEmail
	return usr, nil
}
