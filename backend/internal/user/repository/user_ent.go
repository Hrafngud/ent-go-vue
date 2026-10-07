package repository

import (
	"context"

	"backend-golang/ent"
	entuser "backend-golang/ent/user"
	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type entRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) user.Repository {
	return &entRepository{
		client: client,
	}
}

func (r *entRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	entUser, err := r.client.User.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		Save(ctx)

	if err != nil {
		return nil, writeError(err)
	}

	return toDomainUser(entUser), nil
}

func (r *entRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	entUser, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}

	return toDomainUser(entUser), nil
}

func (r *entRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	entUser, err := r.client.User.Query().Where(entuser.Email(email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}

	return toDomainUser(entUser), nil
}

func (r *entRepository) List(ctx context.Context) ([]*user.User, error) {
	entUsers, err := r.client.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var users []*user.User
	for _, eu := range entUsers {
		users = append(users, toDomainUser(eu))
	}

	return users, nil
}

func (r *entRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return user.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *entRepository) Update(ctx context.Context, id uuid.UUID, u *user.User) (*user.User, error) {
	update := r.client.User.UpdateOneID(id).SetName(u.Name).SetEmail(u.Email)
	if u.Password != "" {
		update.SetPassword(u.Password)
	}
	updated, err := update.Save(ctx)
	if err != nil {
		return nil, writeError(err)
	}
	return toDomainUser(updated), nil
}

func writeError(err error) error {
	if ent.IsNotFound(err) {
		return user.ErrNotFound
	}
	if ent.IsConstraintError(err) {
		return user.ErrEmailTaken
	}
	return err
}

func toDomainUser(u *ent.User) *user.User {
	if u == nil {
		return nil
	}
	return &user.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}
