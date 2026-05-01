package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestPerformanceBonus_LikelyHitCreditsToken: Blue (1) is the only variant whose printed attack
// lands in the likely-to-hit set; the Gold-token rider credits sim.GoldTokenValue. (Red 3 and
// Yellow 2 are both blockable.)
func TestPerformanceBonus_LikelyHitCreditsToken(t *testing.T) {
	var s sim.TurnState
	c := PerformanceBonusBlue{}
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 1+sim.GoldTokenValue {
		t.Errorf("Blue: Play() = %d, want %d (1 likely to hit + GoldTokenValue)", got, 1+sim.GoldTokenValue)
	}
}

// TestPerformanceBonus_BlockableSuppressesToken: Red (3) and Yellow (2) are blockable; the
// rider doesn't fire.
func TestPerformanceBonus_BlockableSuppressesToken(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{PerformanceBonusRed{}, 3},
		{PerformanceBonusYellow{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no token)", tc.c.Name(), got, tc.want)
		}
	}
}
