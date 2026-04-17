package main

import "testing"

// TestValidateDeckName pins the rejection rules for names that become file paths under mydecks/.
// The concern is path traversal ("../") and Windows-reserved characters that would crash WriteFile.
func TestValidateDeckName(t *testing.T) {
	ok := []string{"viserai-v2", "my best deck", "deck_42", "Viserai"}
	bad := []string{"", ".", "..", "../escape", "path/with/slash", "back\\slash", "has*star", "has?mark"}

	for _, name := range ok {
		if err := validateDeckName(name); err != nil {
			t.Errorf("validateDeckName(%q) = %v, want nil", name, err)
		}
	}
	for _, name := range bad {
		if err := validateDeckName(name); err == nil {
			t.Errorf("validateDeckName(%q) = nil, want error", name)
		}
	}
}
