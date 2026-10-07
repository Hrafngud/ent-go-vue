package user

import (
	"net/mail"
	"strings"
)

// NormalizeInput also validates password byte length because bcrypt accepts at most 72 bytes.
func NormalizeInput(input UserInput, requirePassword bool) (UserInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	address, err := mail.ParseAddress(input.Email)
	if input.Name == "" || len(input.Name) > 255 || len(input.Email) > 255 || err != nil || address.Address != input.Email {
		return input, ErrInvalidInput
	}
	if (requirePassword || input.Password != "") && (len(input.Password) < 6 || len(input.Password) > 72) {
		return input, ErrInvalidInput
	}
	return input, nil
}
