package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSnatch_LikelyHitCreditsDraw: Red (4) is the only variant whose printed attack lands in
// the likely-to-hit set; the draw rider credits +3.
func TestSnatch_LikelyHitCreditsDraw(t *testing.T) {
	var s card.TurnState
	if got := (SnatchRed{}).Play(&s); got != 4+3 {
		t.Errorf("Red: Play() = %d, want 7 (4 likely to hit + 3 draw)", got)
	}
}

// TestSnatch_BlockableSuppressesDraw: Yellow (3) and Blue (2) are blockable totals the opponent
// won't let through, so the draw rider doesn't fire.
func TestSnatch_BlockableSuppressesDraw(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{SnatchYellow{}, 3},
		{SnatchBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no draw)", tc.c.Name(), got, tc.want)
		}
	}
}
