package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestPerformanceBonus_LikelyHitCreditsToken: Blue (1) is the only variant whose printed attack
// lands in the likely-to-hit set; the Gold-token rider credits +1. (Red 3 and Yellow 2 are
// both blockable.)
func TestPerformanceBonus_LikelyHitCreditsToken(t *testing.T) {
	var s card.TurnState
	if got := (PerformanceBonusBlue{}).Play(&s); got != 1+1 {
		t.Errorf("Blue: Play() = %d, want 2 (1 likely to hit + 1 Gold)", got)
	}
}

// TestPerformanceBonus_BlockableSuppressesToken: Red (3) and Yellow (2) are blockable; the
// rider doesn't fire.
func TestPerformanceBonus_BlockableSuppressesToken(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{PerformanceBonusRed{}, 3},
		{PerformanceBonusYellow{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no token)", tc.c.Name(), got, tc.want)
		}
	}
}
