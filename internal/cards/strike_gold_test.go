package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestStrikeGold_LikelyHitCreditsToken: Red (4) is the only variant whose printed attack lands
// in the likely-to-hit set; the Gold-token rider credits sim.GoldTokenValue.
func TestStrikeGold_LikelyHitCreditsToken(t *testing.T) {
	var s sim.TurnState
	c := StrikeGoldRed{}
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 4+sim.GoldTokenValue {
		t.Errorf("Red: Play() = %d, want %d (4 likely to hit + GoldTokenValue)", got, 4+sim.GoldTokenValue)
	}
}

// TestStrikeGold_BlockableSuppressesToken: Yellow (3) and Blue (2) are blockable; the rider
// doesn't fire.
func TestStrikeGold_BlockableSuppressesToken(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{StrikeGoldYellow{}, 3},
		{StrikeGoldBlue{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no token)", tc.c.Name(), got, tc.want)
		}
	}
}
