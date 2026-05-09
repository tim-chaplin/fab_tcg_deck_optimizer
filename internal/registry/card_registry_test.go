package registry

import (
	"testing"
)

func TestAllIDsResolve(t *testing.T) {
	// Every CardID returned by AllCards() must map to a non-nil card. Catches gaps in the byID
	// slice (an undeclared const would leave a nil hole).
	for _, id := range AllCards() {
		if GetCard(id) == nil {
			t.Errorf("CardID %d resolves to nil", id)
		}
	}
}

func TestDisplayNamesAreUnique(t *testing.T) {
	// c.DisplayName() is used as the reverse-lookup key, so every registered card must
	// have a distinct display name. A collision would silently overwrite the earlier entry
	// in byName. (Bare Name() collides intentionally — pitch variants share it.)
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
		t.Error("ByName of unknown card should return ok=false")
	}
}
