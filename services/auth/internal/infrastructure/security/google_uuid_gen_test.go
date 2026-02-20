package security

import (
	"testing"
)

func TestGoogleUUIDGen_NewId(t *testing.T) {
	gen := &GoogleUUIDGen{}

	t.Run("generates non-empty ID", func(t *testing.T) {
		id := gen.NewId()
		if id == "" {
			t.Error("expected non-empty ID")
		}
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		const n = 100
		seen := make(map[string]struct{}, n)
		for i := 0; i < n; i++ {
			id := gen.NewId()
			if _, exists := seen[id]; exists {
				t.Fatalf("duplicate ID generated: %s", id)
			}
			seen[id] = struct{}{}
		}
	})

	t.Run("ID is valid UUID v7 format", func(t *testing.T) {
		id := gen.NewId()
		// UUID v7 format: xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx (36 chars with dashes)
		if len(id) != 36 {
			t.Errorf("expected UUID length 36, got %d for %q", len(id), id)
		}
		if id[14] != '7' {
			t.Errorf("expected UUID version 7, got version char %q in %q", id[14], id)
		}
	})
}
