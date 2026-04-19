package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestStrikeGold_LikelyHitCreditsToken: Red (4) is the only variant whose printed attack lands
// in the likely-to-hit set; the Gold-token rider credits card.GoldTokenValue.
func TestStrikeGold_LikelyHitCreditsToken(t *testing.T) {
	var s card.TurnState
	if got := (StrikeGoldRed{}).Play(&s); got != 4+card.GoldTokenValue {
		t.Errorf("Red: Play() = %d, want %d (4 likely to hit + GoldTokenValue)", got, 4+card.GoldTokenValue)
	}
}

// TestStrikeGold_BlockableSuppressesToken: Yellow (3) and Blue (2) are blockable; the rider
// doesn't fire.
func TestStrikeGold_BlockableSuppressesToken(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{StrikeGoldYellow{}, 3},
		{StrikeGoldBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no token)", tc.c.Name(), got, tc.want)
		}
	}
}
