package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Potion of Seeing's Play credits nothing (the activated reveal is dropped as
// pure information).
func TestPotionOfSeeing_PlayCreditsNothing(t *testing.T) {
	c := PotionOfSeeingBlue{}
	var s sim.TurnState
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 0 {
		t.Errorf("Play() Value = %d, want 0", got)
	}
}
