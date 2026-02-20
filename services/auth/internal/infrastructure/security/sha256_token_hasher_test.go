package security

import (
	"testing"

	"github.com/aligh5331/godrop/services/auth/internal/domain/entity"
)

func TestSHA256TokenHasher_Hash(t *testing.T) {
	h := SHA256TokenHasher{}

	t.Run("produces non-empty hash", func(t *testing.T) {
		hashed, err := h.Hash("my-refresh-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hashed == "" {
			t.Error("expected non-empty hash")
		}
	})

	t.Run("same input produces same hash", func(t *testing.T) {
		token := "consistent-token"
		h1, _ := h.Hash(token)
		h2, _ := h.Hash(token)
		if h1 != h2 {
			t.Errorf("same token produced different hashes: %s vs %s", h1, h2)
		}
	})

	t.Run("different inputs produce different hashes", func(t *testing.T) {
		h1, _ := h.Hash("token-a")
		h2, _ := h.Hash("token-b")
		if h1 == h2 {
			t.Error("different tokens should not produce the same hash")
		}
	})

	t.Run("empty token hashes without error", func(t *testing.T) {
		hashed, err := h.Hash("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hashed == "" {
			t.Error("expected a hash even for empty string")
		}
	})
}

func TestSHA256TokenHasher_CheckHash(t *testing.T) {
	h := SHA256TokenHasher{}

	token := "my-secret-refresh-token"
	hashed, _ := h.Hash(token)
	emptyHash, _ := h.Hash("")
	tests := []struct {
		name     string
		token    string
		hashed   entity.HashedToken
		expected bool
	}{
		{"matching token and hash", token, hashed, true},
		{"wrong token", "wrong-token", hashed, false},
		{"empty token", "", hashed, false},
		{"tampered hash", token, entity.HashedToken(string(hashed) + "x"), false},
		{"both empty", "", emptyHash, true}, // sha256("") == sha256("") is a valid match
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.CheckHash(tt.token, tt.hashed)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
