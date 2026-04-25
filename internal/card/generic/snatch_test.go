package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSnatch_LikelyHitFiresDrawOne: Red (attack 4) is the only variant whose printed attack
// lands in the likely-to-hit set. The rider fires state.DrawOne, advancing the deck and
// recording the drawn card in state.Drawn. Play returns just the attack.
func TestSnatch_LikelyHitFiresDrawOne(t *testing.T) {
	top := stubGenericAttack(0, 3)
	s := card.TurnState{Deck: []card.Card{top}}
	c := SnatchRed{}
	if got := c.Play(&s, &card.CardState{Card: c}); got != 4 {
		t.Errorf("Red: Play() = %d, want 4", got)
	}
	if len(s.Drawn) != 1 || s.Drawn[0] != top {
		t.Errorf("Drawn = %v, want [top-of-deck]", s.Drawn)
	}
	if len(s.Deck) != 0 {
		t.Errorf("Deck len = %d, want 0 (top consumed)", len(s.Deck))
	}
}

// TestSnatch_BlockableSuppressesDraw: Yellow (3) and Blue (2) are blockable; the on-hit rider
// doesn't fire, so the deck is untouched and Drawn stays empty.
func TestSnatch_BlockableSuppressesDraw(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{SnatchYellow{}, 3},
		{SnatchBlue{}, 2},
	}
	for _, tc := range cases {
		top := stubGenericAttack(0, 3)
		s := card.TurnState{Deck: []card.Card{top}}
		if got := tc.c.Play(&s, &card.CardState{Card: tc.c}); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no draw)", tc.c.Name(), got, tc.want)
		}
		if len(s.Drawn) != 0 {
			t.Errorf("%s: Drawn = %v, want empty (no draw fired)", tc.c.Name(), s.Drawn)
		}
		if len(s.Deck) != 1 {
			t.Errorf("%s: Deck len = %d, want 1 (top preserved)", tc.c.Name(), len(s.Deck))
		}
	}
}
