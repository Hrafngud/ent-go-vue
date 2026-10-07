package usecase

import (
	"context"
	"errors"
	"time"

	"backend-golang/internal/auth"
	"backend-golang/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo  user.Repository
	secret    []byte
	dummyHash []byte
}

func NewUsecase(userRepo user.Repository, secret []byte) auth.Usecase {
	// Compare against a real hash even for missing accounts to reduce timing leaks.
	dummyHash, err := bcrypt.GenerateFromPassword([]byte("unused-account-password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return &authUsecase{userRepo: userRepo, secret: secret, dummyHash: dummyHash}
}

func (u *authUsecase) Register(ctx context.Context, name, email, password string) error {
	input, err := user.NormalizeInput(user.UserInput{Name: name, Email: email, Password: password}, true)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	usr := &user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}
	_, err = u.userRepo.Create(ctx, usr)
	return err
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	email, err := user.NormalizeEmail(email)
	// Existing shorter passwords remain usable; new passwords require 8 bytes.
	if err != nil || !user.ValidPassword(password, 1) {
		return "", errors.New("invalid email or password")
	}
	usr, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword(u.dummyHash, []byte(password))
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(u.secret)
}
