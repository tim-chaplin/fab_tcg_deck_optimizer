package registry

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestAuraTriggerCreatorsOptInToAddsFutureValue: a card whose Play registers a card-style
// Aura hides its real payoff on a future turn — without AddsFutureValue the beatsBest
// tiebreaker would pick a Held → arsenal promotion at equal current-turn Value over
// actually playing the card. Token auras (Runechants, …) are skipped: they get credited
// at creation time and fire on the next attack this turn, so there's no deferred payoff
// the marker would protect against.
func TestAuraTriggerCreatorsOptInToAddsFutureValue(t *testing.T) {
	for _, id := range AllCards() {
		c := GetCard(id)
		var s sim.TurnState
		// self carries the card so Plays that consult self.EffectiveGoAgain /
		// self.EffectiveDominate (reading Card.GoAgain / the Dominator marker) don't
		// nil-dereference.
		c.Play(&s, &sim.CardState{Card: c})
		hasCardAura := false
		for _, a := range s.Auras {
			if a.Self.Card != nil {
				hasCardAura = true
				break
			}
		}
		if !hasCardAura {
			continue
		}
		if _, addsFuture := c.(sim.AddsFutureValue); !addsFuture {
			t.Errorf("%s registers a card-style Aura but doesn't implement AddsFutureValue — beatsBest tiebreaker won't favour playing it",
				c.Name())
		}
	}
}

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
	// sim.DisplayName(c) is used as the reverse-lookup key, so every registered card must
	// have a distinct display name. A collision would silently overwrite the earlier entry
	// in byName. (Bare Name() collides intentionally — pitch variants share it.)
	seen := map[string]CardID{}
	for _, id := range AllCards() {
		name := sim.DisplayName(GetCard(id))
		if prev, dup := seen[name]; dup {
			t.Errorf("duplicate DisplayName %q for IDs %d and %d", name, prev, id)
		}
		seen[name] = id
	}
}

func TestByNameRoundTrip(t *testing.T) {
	for _, id := range AllCards() {
		name := sim.DisplayName(GetCard(id))
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
