package registry

import "testing"

func TestAllIDsResolve(t *testing.T) {
	// Every CardID from AllCards() must map to a non-nil card; catches nil holes in cardsByID.
	for _, id := range AllCards() {
		if GetCard(id) == nil {
			t.Errorf("CardID %d resolves to nil", id)
		}
	}
}

func TestDisplayNamesAreUnique(t *testing.T) {
	// DisplayName is the reverse-lookup key; collisions would silently overwrite the earlier
	// entry in cardsByName. (Bare Name collides intentionally — pitch variants share it.)
	seen := map[string]CardID{}
	for _, id := range AllCards() {
		name := GetCard(id).DisplayName()
		if prev, dup := seen[name]; dup {
			t.Errorf("duplicate DisplayName %q for IDs %d and %d", name, prev, id)
		}
		seen[name] = id
	}
}

func TestByNameRoundTrip(t *testing.T) {
	for _, id := range AllCards() {
		name := GetCard(id).DisplayName()
		got, ok := CardByName(name)
		if !ok || got != id {
			t.Errorf("CardByName(%q) = (%d, %v), want (%d, true)", name, got, ok, id)
		}
	}
}

func TestByNameUnknown(t *testing.T) {
	if _, ok := CardByName("Not A Real Card"); ok {
		t.Error("CardByName of unknown card should return ok=false")
	}
}
