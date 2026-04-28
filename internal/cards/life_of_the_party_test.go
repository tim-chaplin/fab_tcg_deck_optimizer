package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestLifeOfTheParty_BaseGoAgainFalse pins the simplification: go again is one of three random
// modes gated on paying an alternate Crazy Brew cost we don't model, so GoAgain() must return
// false. Returning true would let Life of the Party always chain, over-crediting sequences vs.
// the baseline where the random mode rolled +2{p} or on-hit life instead of go again.
func TestLifeOfTheParty_BaseGoAgainFalse(t *testing.T) {
	for _, c := range []sim.Card{LifeOfThePartyRed{}, LifeOfThePartyYellow{}, LifeOfThePartyBlue{}} {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (random-mode selection not modelled)", c.Name())
		}
	}
}
