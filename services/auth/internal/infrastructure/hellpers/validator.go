package hellpers

import (
	"net"
	"net/mail"
	"strings"
	"unicode"

	"github.com/aligh5331/godrop/services/auth/internal/domain"
)

type Validator struct{}

func (v Validator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return domain.ErrPasswordTooShort
	}
	if len(password) > 72 { // bcrypt limit
		return domain.ErrPasswordTooLong
	}

	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return domain.ErrPasswordTooWeak
	}
	return nil
}

func (v Validator) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return domain.ErrInvalidEmail
	}
	// ParseAddress accepts "Name <email>" format, we want plain email only
	if strings.ContainsAny(email, "<>") {
		return domain.ErrInvalidEmail
	}
	return nil
}

func (v Validator) ValidateIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return domain.ErrInvalidIP
	}
	return nil
}
