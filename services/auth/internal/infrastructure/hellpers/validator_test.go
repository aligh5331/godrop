package hellpers

import (
	"auth/internal/domain"
	"errors"
	"strings"
	"testing"
)

func TestValidator_ValidatePassword(t *testing.T) {
	v := Validator{}

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"valid password", "Password1", nil},
		{"valid long password", "MyPassword123!", nil},
		{"too short", "Pass1", domain.ErrPasswordTooShort},
		{"too long", strings.Repeat("A", 73), domain.ErrPasswordTooLong},
		{"no uppercase", "password1", domain.ErrPasswordTooWeak},
		{"no lowercase", "PASSWORD1", domain.ErrPasswordTooWeak},
		{"no digit", "PasswordA", domain.ErrPasswordTooWeak},
		{"exactly 8 chars valid", "Pass123!", nil},
		{"exactly 72 chars valid", strings.Repeat("A", 35) + strings.Repeat("a", 35) + "12", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidatePassword(tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidator_ValidateEmail(t *testing.T) {
	v := Validator{}

	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{"valid email", "user@example.com", nil},
		{"valid subdomain", "user@mail.example.com", nil},
		{"empty email", "", domain.ErrInvalidEmail},
		{"no at sign", "userexample.com", domain.ErrInvalidEmail},
		{"no domain", "user@", domain.ErrInvalidEmail},
		{"name format rejected", "Name <user@example.com>", domain.ErrInvalidEmail},
		{"angle bracket in email", "<user@example.com>", domain.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateEmail(tt.email)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidator_ValidateIP(t *testing.T) {
	v := Validator{}

	tests := []struct {
		name    string
		ip      string
		wantErr error
	}{
		{"valid IPv4", "192.168.1.1", nil},
		{"valid IPv4 loopback", "127.0.0.1", nil},
		{"valid IPv6", "::1", nil},
		{"valid IPv6 full", "2001:db8::1", nil},
		{"empty ip", "", domain.ErrInvalidIP},
		{"invalid ip", "999.999.999.999", domain.ErrInvalidIP},
		{"text", "notanip", domain.ErrInvalidIP},
		{"partial ip", "192.168.1", domain.ErrInvalidIP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateIP(tt.ip)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
