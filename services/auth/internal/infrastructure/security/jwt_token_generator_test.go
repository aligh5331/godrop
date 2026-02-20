package security

import (
	"testing"
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/application/dto"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTTokenGenerator_GenerateAccessToken(t *testing.T) {
	secret := []byte("test-secret-key")
	gen := NewJWTTokenGenerator(secret)

	metadata := dto.SessionMetadataDTO{
		IP:          "192.168.1.1",
		ClientAgent: "Mozilla/5.0",
	}

	t.Run("generates valid signed token", func(t *testing.T) {
		tokenStr, err := gen.GenerateAccessToken("user-123", metadata, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tokenStr == "" {
			t.Error("expected non-empty token")
		}

		// Parse and verify claims
		token, err := jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			t.Fatalf("token parse error: %v", err)
		}
		claims, ok := token.Claims.(*accessClaims)
		if !ok || !token.Valid {
			t.Fatal("invalid token claims")
		}
		if claims.UserID != "user-123" {
			t.Errorf("expected UserID user-123, got %s", claims.UserID)
		}
		if claims.IP != metadata.IP {
			t.Errorf("expected IP %s, got %s", metadata.IP, claims.IP)
		}
		if claims.ClientAgent != metadata.ClientAgent {
			t.Errorf("expected ClientAgent %s, got %s", metadata.ClientAgent, claims.ClientAgent)
		}
	})

	t.Run("token expires correctly", func(t *testing.T) {
		tokenStr, err := gen.GenerateAccessToken("user-456", metadata, -time.Second) // already expired
		if err != nil {
			t.Fatalf("unexpected error generating: %v", err)
		}

		_, err = jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err == nil {
			t.Error("expected error for expired token")
		}
	})
}

func TestJWTTokenGenerator_GenerateRefreshToken(t *testing.T) {
	secret := []byte("test-secret-key")
	gen := NewJWTTokenGenerator(secret)

	t.Run("generates valid refresh token with session claims", func(t *testing.T) {
		tokenStr, err := gen.GenerateRefreshToken("user-123", "session-abc", 24*time.Hour)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		token, err := jwt.ParseWithClaims(tokenStr, &refreshClaims{}, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			t.Fatalf("token parse error: %v", err)
		}
		claims, ok := token.Claims.(*refreshClaims)
		if !ok || !token.Valid {
			t.Fatal("invalid token claims")
		}
		if claims.UserID != "user-123" {
			t.Errorf("expected UserID user-123, got %s", claims.UserID)
		}
		if claims.SessionID != "session-abc" {
			t.Errorf("expected SessionID session-abc, got %s", claims.SessionID)
		}
	})

	t.Run("wrong secret fails verification", func(t *testing.T) {
		tokenStr, _ := gen.GenerateRefreshToken("user-123", "session-abc", time.Hour)
		_, err := jwt.ParseWithClaims(tokenStr, &refreshClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte("wrong-secret"), nil
		})
		if err == nil {
			t.Error("expected signature verification to fail with wrong secret")
		}
	})
}
