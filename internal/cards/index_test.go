package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestAuraTriggerCreatorsOptInToAddsFutureValue: a card whose Play registers an AuraTrigger
// hides its real payoff on a future turn — without AddsFutureValue the beatsBest tiebreaker
// would pick a Held → arsenal promotion at equal current-turn Value over actually playing
// the card. Probes Play with a fresh TurnState and checks every registrant carries the
// marker.
func TestAuraTriggerCreatorsOptInToAddsFutureValue(t *testing.T) {
	for _, id := range All() {
		c := Get(id)
		var s card.TurnState
		// self carries the card so Plays that consult self.EffectiveGoAgain /
		// self.EffectiveDominate (reading Card.GoAgain / the Dominator marker) don't
		// nil-dereference.
		c.Play(&s, &card.CardState{Card: c})
		if len(s.AuraTriggers) == 0 {
			continue
		}
		if _, addsFuture := c.(card.AddsFutureValue); !addsFuture {
			t.Errorf("%s registers an AuraTrigger but doesn't implement AddsFutureValue — beatsBest tiebreaker won't favour playing it",
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

func TestDisplayNamesAreUnique(t *testing.T) {
	// card.DisplayName(c) is used as the reverse-lookup key, so every registered card must
	// have a distinct display name. A collision would silently overwrite the earlier entry
	// in byName. (Bare Name() collides intentionally — pitch variants share it.)
	seen := map[string]ID{}
	for _, id := range All() {
		name := card.DisplayName(Get(id))
		if prev, dup := seen[name]; dup {
			t.Errorf("duplicate DisplayName %q for IDs %d and %d", name, prev, id)
		}
		seen[name] = id
	}
}

func TestByNameRoundTrip(t *testing.T) {
	for _, id := range All() {
		name := card.DisplayName(Get(id))
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
