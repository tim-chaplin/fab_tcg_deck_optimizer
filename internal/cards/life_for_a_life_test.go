package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestLifeForALife_LikelyHitCreditsHeal: Red (4) is the only variant whose printed attack lands
// in the likely-to-hit set; the 1{h} gain credits +1.
func TestLifeForALife_LikelyHitCreditsHeal(t *testing.T) {
	var s sim.TurnState
	c := LifeForALifeRed{}
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 4+1 {
		t.Errorf("Red: Play() = %d, want 5 (4 likely to hit + 1 heal)", got)
	}
}

// TestLifeForALife_BlockableSuppressesHeal: Yellow (3) and Blue (2) are blockable; the heal
// rider doesn't fire.
func TestLifeForALife_BlockableSuppressesHeal(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{LifeForALifeYellow{}, 3},
		{LifeForALifeBlue{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no heal)", tc.c.Name(), got, tc.want)
		}
	}
}
