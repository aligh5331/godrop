package entity

import (
	"auth/internal/domain"
	"errors"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	return time.Date(2024, 1, 1, 10, 0, 0, 0, time.FixedZone("TEST", 3*60*60))
}

func TestNewUser_Success(t *testing.T) {
	now := fixedTime()

	u, err := NewUser(
		"  id-1  ",
		"  Ali  ",
		"  ali@test.com  ",
		HashedPassword("hashed"),
		now,
		now,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Id() != "id-1" {
		t.Errorf("expected trimmed id, got %q", u.Id())
	}

	if u.Name() != "Ali" {
		t.Errorf("expected trimmed name, got %q", u.Name())
	}

	if u.Email() != "ali@test.com" {
		t.Errorf("expected trimmed email, got %q", u.Email())
	}

	if !u.CreatedAt().Equal(now.UTC()) {
		t.Errorf("expected createdAt in UTC")
	}

	if !u.UpdatedAt().Equal(now.UTC()) {
		t.Errorf("expected updatedAt in UTC")
	}
}

func TestNewUser_ValidationErrors(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name string
		id   string
		n    string
		e    string
		p    HashedPassword
		err  error
	}{
		{"empty id", "", "Ali", "a@test.com", "hash", domain.ErrEmptyId},
		{"empty name", "id", "", "a@test.com", "hash", domain.ErrEmptyName},
		{"empty email", "id", "Ali", "", "hash", domain.ErrEmptyEmail},
		{"empty password", "id", "Ali", "a@test.com", "", domain.ErrEmptyPassword},
		{"whitespace id", "   ", "Ali", "a@test.com", "hash", domain.ErrEmptyId},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.id, tt.n, tt.e, tt.p, now, now)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestSetHashedPassword(t *testing.T) {
	now := fixedTime()

	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	newTime := now.Add(time.Hour)
	u.SetHashedPassword("newhash", newTime)

	if u.Password() != "newhash" {
		t.Errorf("password not updated")
	}

	if !u.UpdatedAt().Equal(newTime.UTC()) {
		t.Errorf("updatedAt not set correctly")
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	updateTime := now.Add(2 * time.Hour)
	err := u.UpdateProfile("  NewName  ", updateTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Name() != "NewName" {
		t.Errorf("name not updated or not trimmed")
	}

	if !u.UpdatedAt().Equal(updateTime.UTC()) {
		t.Errorf("updatedAt not updated correctly")
	}
}

func TestUpdateProfile_InvalidName(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	oldName := u.Name()

	err := u.UpdateProfile("   ", now.Add(time.Hour))
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}

	// Ensure state was NOT mutated
	if u.Name() != oldName {
		t.Errorf("user mutated despite validation error")
	}
}

func TestChangePassword(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	before := u.UpdatedAt()
	err := u.ChangePassword("newhash", now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Password() != "newhash" {
		t.Errorf("password not updated")
	}

	if !u.UpdatedAt().After(before) {
		t.Errorf("updatedAt should be updated to a newer time")
	}
}

func TestChangePassword_Empty(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	before := u.UpdatedAt()
	err := u.ChangePassword("   ", now.Add(1*time.Hour))
	if !errors.Is(err, domain.ErrEmptyPassword) {
		t.Fatalf("expected ErrEmptyPassword")
	}

	if u.UpdatedAt().After(before) {
		t.Errorf("updatedAt should not be updated to a newer time")
	}

}

func TestChangeName(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	before := u.UpdatedAt()
	err := u.ChangeName("  New  ", now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Name() != "New" {
		t.Errorf("name not trimmed correctly")
	}

	if !u.UpdatedAt().After(before) {
		t.Errorf("updatedAt should be updated to a newer time")
	}

}

func TestChangeEmail(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	before := u.UpdatedAt()
	err := u.ChangeEmail("  new@test.com  ", now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Email() != "new@test.com" {
		t.Errorf("email not trimmed correctly")
	}

	if !u.UpdatedAt().After(before) {
		t.Errorf("updatedAt should be updated to a newer time")
	}
}

func TestChangeEmail_Empty(t *testing.T) {
	now := fixedTime()
	u, _ := NewUser("id", "Ali", "a@test.com", "hash", now, now)

	err := u.ChangeEmail(strings.Repeat(" ", 5), now.Add(1*time.Hour))
	if !errors.Is(err, domain.ErrEmptyEmail) {
		t.Fatalf("expected ErrEmptyEmail")
	}
}
