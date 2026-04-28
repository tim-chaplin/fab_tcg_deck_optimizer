package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestSnatch_LikelyHitFiresDrawOne: Red (attack 4) is the only variant whose printed attack
// lands in the likely-to-hit set. The rider fires state.DrawOne, popping the deck top and
// appending it to s.Hand. Play returns just the attack.
func TestSnatch_LikelyHitFiresDrawOne(t *testing.T) {
	top := testutils.GenericAttack(0, 3)
	s := sim.NewTurnState([]sim.Card{top}, nil)
	c := SnatchRed{}
	c.Play(s, &sim.CardState{Card: c})
	if got := s.Value; got != 4 {
		t.Errorf("Red: Play() = %d, want 4", got)
	}
	if len(s.Hand) != 1 || s.Hand[0] != top {
		t.Errorf("Hand = %v, want [top-of-deck]", s.Hand)
	}
	if d := s.Deck(); len(d) != 0 {
		t.Errorf("Deck len = %d, want 0 (top consumed)", len(d))
	}
}

// TestSnatch_BlockableSuppressesDraw: Yellow (3) and Blue (2) are blockable; the on-hit rider
// doesn't fire, so the deck is untouched and Hand stays empty.
func TestSnatch_BlockableSuppressesDraw(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{SnatchYellow{}, 3},
		{SnatchBlue{}, 2},
	}
	for _, tc := range cases {
		top := testutils.GenericAttack(0, 3)
		s := sim.NewTurnState([]sim.Card{top}, nil)
		tc.c.Play(s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no draw)", tc.c.Name(), got, tc.want)
		}
		if len(s.Hand) != 0 {
			t.Errorf("%s: Hand = %v, want empty (no draw fired)", tc.c.Name(), s.Hand)
		}
		if d := s.Deck(); len(d) != 1 {
			t.Errorf("%s: Deck len = %d, want 1 (top preserved)", tc.c.Name(), len(d))
		}
	}
}
