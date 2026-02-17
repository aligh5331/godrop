package security

import (
	"auth/internal/domain"
	"auth/internal/domain/entity"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewHasher(t *testing.T) {
	tests := []struct {
		name     string
		cost     int
		expected int
	}{
		{
			name:     "Valid cost within range",
			cost:     12,
			expected: 12,
		},
		{
			name:     "Cost below minimum",
			cost:     bcrypt.MinCost - 1,
			expected: bcrypt.DefaultCost,
		},
		{
			name:     "Cost above maximum",
			cost:     bcrypt.MaxCost + 1,
			expected: bcrypt.DefaultCost,
		},
		{
			name:     "Default cost",
			cost:     bcrypt.DefaultCost,
			expected: bcrypt.DefaultCost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewHasher(tt.cost)
			bcryptHasher, ok := hasher.(*BcryptHasher)
			if !ok {
				t.Fatalf("expected *BcryptHasher, got %T", hasher)
			}
			if bcryptHasher.cost != tt.expected {
				t.Errorf("expected cost %d, got %d", tt.expected, bcryptHasher.cost)
			}
		})
	}
}

func TestBcryptHasher_Hash(t *testing.T) {
	hasher := NewHasher(bcrypt.DefaultCost)

	tests := []struct {
		name        string
		password    string
		expectError error
	}{
		{
			name:        "Successful hash",
			password:    "validpassword123",
			expectError: nil,
		},
		{
			name:        "Password too long",
			password:    string(make([]byte, 73)), // bcrypt max is 72 bytes
			expectError: domain.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := hasher.Hash(tt.password)
			if tt.expectError != nil {
				if err == nil || err.Error() != tt.expectError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectError, err)
				}
				if hashed != "" {
					t.Errorf("expected empty hash on error, got %s", hashed)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if hashed == "" {
					t.Error("expected non-empty hash")
				}
				// Verify it's a valid bcrypt hash
				err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(tt.password))
				if err != nil {
					t.Errorf("hash does not match password: %v", err)
				}
			}
		})
	}
}

func TestBcryptHasher_CheckPasswordHash(t *testing.T) {
	hasher := NewHasher(bcrypt.DefaultCost)

	// Generate a valid hash for testing
	validPassword := "validpassword123"
	hashed, err := hasher.Hash(validPassword)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     entity.HashedPassword
		expected bool
	}{
		{
			name:     "Matching password and hash",
			password: validPassword,
			hash:     hashed,
			expected: true,
		},
		{
			name:     "Non-matching password",
			password: "wrongpassword",
			hash:     hashed,
			expected: false,
		},
		{
			name:     "Empty hash",
			password: validPassword,
			hash:     "",
			expected: false,
		},
		{
			name:     "Password longer than 72 bytes",
			password: string(make([]byte, 73)),
			hash:     hashed,
			expected: false,
		},
		{
			name:     "Empty password with valid hash",
			password: "",
			hash:     hashed,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasher.CheckPasswordHash(tt.password, tt.hash)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
