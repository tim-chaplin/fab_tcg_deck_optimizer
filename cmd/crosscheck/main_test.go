package main

import "testing"

// TestInScope pins the typebox filter: only no-talent Runeblade/Generic main-deck cards pass.
// The tricky cases are multi-word card types (no " - " separator) and double-faced cards whose
// other face is talented.
func TestInScope(t *testing.T) {
	cases := []struct {
		typebox string
		want    bool
	}{
		{"Runeblade Action - Attack", true},
		{"Generic Action - Attack", true},
		{"Generic Instant", true},
		{"Runeblade Action", true},
		{"Generic Defense Reaction", true}, // multi-word type, no " - "
		{"Generic Attack Reaction", true},
		{"Runeblade Aura", true},
		// Talented — Viserai can't run these.
		{"Lightning Runeblade Action - Attack", false},
		{"Shadow Runeblade Aura", false},
		{"Earth Runeblade Action", false},
		// Double-faced card with a talented second face.
		{"Runeblade Action||Lightning Instant", false},
		{"Runeblade Action||Earth Instant", false},
		// Off-class.
		{"Warrior Action - Attack", false},
		{"Runeblade Wizard Action", false},
		// Non-deck card types.
		{"Runeblade Weapon - Sword", false},
		{"Generic Equipment - Arms", false},
		{"Runeblade Hero - Young", false},
		{"Generic Token Action - Attack", false},
		// Royal is a supertype, not a talent — Royal Generics stay legal.
		{"Generic Action - Attack Reaction", true},
	}
	for _, c := range cases {
		if got := inScope(c.typebox); got != c.want {
			t.Errorf("inScope(%q) = %v, want %v", c.typebox, got, c.want)
		}
	}
}
