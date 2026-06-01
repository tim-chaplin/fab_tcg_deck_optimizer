package fabtype

import "testing"

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
