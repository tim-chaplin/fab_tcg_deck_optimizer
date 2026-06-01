package fabtype

import "testing"

// TestClassMatches pins that Generic is universally legal (never needs listing per hero), a
// non-class word imposes no class constraint, and a real class is legal only for a hero that
// plays it.
func TestClassMatches(t *testing.T) {
	viserai := map[string]bool{"Runeblade": true} // hero's own class only — no Generic
	cases := []struct {
		word string
		want bool
	}{
		{"Generic", true},   // universally legal
		{"Runeblade", true}, // hero's class
		{"Warrior", false},  // off-class
		{"Lightning", true}, // a talent isn't a class constraint
		{"Action", true},    // a card type isn't a class constraint
	}
	for _, c := range cases {
		if got := ClassMatches(c.word, viserai); got != c.want {
			t.Errorf("ClassMatches(%q) = %v, want %v", c.word, got, c.want)
		}
	}
}

// TestUnmodeledTypesMatchesNonDeckMinusWeapon pins the documented invariant: UnmodeledTypes is
// NonDeckTypes without Weapon (the one non-deck type the optimizer does model). Without this,
// the two maps could drift silently.
func TestUnmodeledTypesMatchesNonDeckMinusWeapon(t *testing.T) {
	for typ := range NonDeckTypes {
		want := typ != "Weapon"
		if UnmodeledTypes[typ] != want {
			t.Errorf("UnmodeledTypes[%q] = %v, want %v", typ, UnmodeledTypes[typ], want)
		}
	}
	for typ := range UnmodeledTypes {
		if !NonDeckTypes[typ] {
			t.Errorf("UnmodeledTypes has %q, which is not a NonDeckType", typ)
		}
	}
}
