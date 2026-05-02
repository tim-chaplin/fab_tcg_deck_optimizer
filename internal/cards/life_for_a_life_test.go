package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the on-hit 1{h} gain credits +1 on a likely-hit attack.
func TestLifeForALife_LikelyHitCreditsHeal(t *testing.T) {
	var s sim.TurnState
	c := LifeForALifeRed{}
	cs := &sim.CardState{Card: c}
	c.Play(&s, cs)
	testutils.FireOnHitIfLikely(&s, cs)
	if got := s.Value; got != 4+1 {
		t.Errorf("Red: Play() = %d, want 5 (4 likely to hit + 1 heal)", got)
	}
}

// Tests that the heal rider doesn't fire on blockable variants.
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
