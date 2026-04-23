package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestDelayedPlayImpliesAddsFutureValue: every card with a next-turn callback should also opt
// into the AddsFutureValue beatsBest tiebreaker — its damage / token / reveal payoff doesn't
// land in this turn's Value, so without the marker the solver would pick a Held → arsenal
// promotion at equal Value over actually playing the card.
func TestDelayedPlayImpliesAddsFutureValue(t *testing.T) {
	for _, id := range All() {
		c := Get(id)
		if _, isDelayed := c.(card.DelayedPlay); !isDelayed {
			continue
		}
		if _, addsFuture := c.(card.AddsFutureValue); !addsFuture {
			t.Errorf("%s implements DelayedPlay but not AddsFutureValue — beatsBest tiebreaker won't favour playing it",
				c.Name())
		}
	}
}

func TestAllIDsResolve(t *testing.T) {
	// Every ID returned by All() must map to a non-nil card. Catches gaps in the byID slice (an
	// undeclared const would leave a nil hole).
	for _, id := range All() {
		if Get(id) == nil {
			t.Errorf("ID %d resolves to nil", id)
		}
	}
}

func TestNamesAreUnique(t *testing.T) {
	// Card.Name() is used as the reverse-lookup key, so every registered card must have a distinct
	// name. A collision would silently overwrite the earlier entry in byName.
	seen := map[string]ID{}
	for _, id := range All() {
		name := Get(id).Name()
		if prev, dup := seen[name]; dup {
			t.Errorf("duplicate Name() %q for IDs %d and %d", name, prev, id)
		}
		seen[name] = id
	}
}

func TestByNameRoundTrip(t *testing.T) {
	for _, id := range All() {
		name := Get(id).Name()
		got, ok := ByName(name)
		if !ok || got != id {
			t.Errorf("ByName(%q) = (%d, %v), want (%d, true)", name, got, ok, id)
		}
	}
}

func TestByNameUnknown(t *testing.T) {
	if _, ok := ByName("Not A Real Card"); ok {
		t.Error("ByName of unknown card should return ok=false")
	}
}
