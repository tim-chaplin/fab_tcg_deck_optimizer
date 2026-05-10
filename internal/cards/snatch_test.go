package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the on-hit DrawOne fires on a likely-hit attack, popping the deck top into
// hand.
func TestSnatch_LikelyHitFiresDrawOne(t *testing.T) {
	top := testutils.GenericAttack(0, 3)
	s := sim.NewTurnStateFromCards([]sim.Card{top}, nil)
	c := SnatchRed{}
	cs := &sim.CardState{Card: c}
	c.Play(s, s.Logger(), cs)
	testutils.FireOnHitIfLikely(s, s.Logger(), cs)
	if got := s.Value; got != 4 {
		t.Errorf("Red: Play() = %d, want 4", got)
	}
	if h := s.Hand(); len(h) != 1 || h[0] != top {
		t.Errorf("Hand = %v, want [top-of-deck]", h)
	}
	if d := s.Deck(); d.Size() != 0 {
		t.Errorf("Deck size = %d, want 0 (top consumed)", d.Size())
	}
}

// Tests that the on-hit DrawOne doesn't fire on blockable variants.
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
		s := sim.NewTurnStateFromCards([]sim.Card{top}, nil)
		tc.c.Play(s, s.Logger(), &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no draw)", tc.c.Name(), got, tc.want)
		}
		if h := s.Hand(); len(h) != 0 {
			t.Errorf("%s: Hand = %v, want empty (no draw fired)", tc.c.Name(), h)
		}
		if d := s.Deck(); d.Size() != 1 {
			t.Errorf("%s: Deck size = %d, want 1 (top preserved)", tc.c.Name(), d.Size())
		}
	}
}
