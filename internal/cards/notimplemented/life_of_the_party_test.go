package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestLifeOfTheParty_BaseGoAgainFalse pins GoAgain() = false: go again is one of three random
// Crazy Brew modes the sim doesn't model, so it must not be granted unconditionally.
func TestLifeOfTheParty_BaseGoAgainFalse(t *testing.T) {
	for _, c := range []sim.Card{LifeOfThePartyRed{}, LifeOfThePartyYellow{}, LifeOfThePartyBlue{}} {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (random-mode selection not modelled)", c.Name())
		}
	}
}
