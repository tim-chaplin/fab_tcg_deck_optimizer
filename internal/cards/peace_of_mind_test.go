package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestPeaceOfMind_PreventsByPrinting pins each printing's prevention amount: the rider
// converts the printed Defense via DealEffectiveDefense, so playing into 10 incoming
// damage credits 4/3/2 prevented for Red/Yellow/Blue.
func TestPeaceOfMind_PreventsByPrinting(t *testing.T) {
	cases := []struct {
		card sim.Card
		want int
	}{
		{PeaceOfMindRed{}, 4},
		{PeaceOfMindYellow{}, 3},
		{PeaceOfMindBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{IncomingDamage: 10}
		self := &sim.CardState{Card: tc.card}
		tc.card.Play(&s, self)
		if s.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), s.Value, tc.want)
		}
	}
}

// TestPeaceOfMind_CreatesPonder pins the on-play Ponder rider: every printing creates
// exactly one Ponder token regardless of the prevention amount, and the AuraCreated
// flag flips so same-turn "if you've created an aura" effects see it.
func TestPeaceOfMind_CreatesPonder(t *testing.T) {
	for _, c := range []sim.Card{PeaceOfMindRed{}, PeaceOfMindYellow{}, PeaceOfMindBlue{}} {
		s := sim.TurnState{IncomingDamage: 10}
		self := &sim.CardState{Card: c}
		c.Play(&s, self)
		if got := s.Ponders(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated = false, want true after creating a Ponder", c.Name())
		}
	}
}
