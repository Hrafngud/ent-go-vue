package user

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Use the same practical email syntax as the browser. Display names, quoted
// local parts, and Unicode addresses are not supported by this application.
var emailPattern = regexp.MustCompile("^[A-Za-z0-9.!#$%&'*+/=?^_`{|}~-]+@[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?(?:\\.[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?)+$")

func NormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if len(email) > 254 || !emailPattern.MatchString(email) {
		return "", ErrInvalidInput
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts[0]) > 64 || strings.HasPrefix(parts[0], ".") || strings.HasSuffix(parts[0], ".") || strings.Contains(parts[0], "..") {
		return "", ErrInvalidInput
	}
	for _, label := range strings.Split(parts[1], ".") {
		if len(label) > 63 {
			return "", ErrInvalidInput
		}
	}
	// Preserve case to keep existing account identities and root configuration.
	return email, nil
}

func ValidPassword(password string, minimum int) bool {
	return utf8.ValidString(password) && len(password) >= minimum && len(password) <= 72 && strings.TrimSpace(password) != "" && !strings.ContainsRune(password, '\x00')
}

// NormalizeInput also validates password byte length because bcrypt accepts at most 72 bytes.
func NormalizeInput(input UserInput, requirePassword bool) (UserInput, error) {
	if !utf8.ValidString(input.Name) {
		return input, ErrInvalidInput
	}
	input.Name = norm.NFC.String(strings.TrimSpace(input.Name))
	var err error
	input.Email, err = NormalizeEmail(input.Email)
	if input.Name == "" || len(input.Name) > 255 || err != nil || strings.ContainsAny(input.Name, "<>") || strings.ContainsFunc(input.Name, unicode.IsControl) {
		return input, ErrInvalidInput
	}
	if (requirePassword || input.Password != "") && !ValidPassword(input.Password, 8) {
		return input, ErrInvalidInput
	}
	return input, nil
}
