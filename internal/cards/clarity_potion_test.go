package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Clarity Potion's Play credits Opt 2 (treating the activated ability as
// always activated on the same turn).
func TestClarityPotion_PlayCreditsOpt2(t *testing.T) {
	c := ClarityPotionBlue{}
	var s sim.TurnState
	c.Play(&s, &sim.CardState{Card: c})
	if want := 2 * sim.OptValue; s.Value != want {
		t.Errorf("Play() Value = %d, want %d", s.Value, want)
	}
}
