package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that each printing's Play credits its prevention cap (3/2/1) when IncomingDamage
// has room.
func TestBrushOff_PreventsCap(t *testing.T) {
	cases := []struct {
		card card.Card
		want int
	}{
		{cards.BrushOffRed{}, 3},
		{cards.BrushOffYellow{}, 2},
		{cards.BrushOffBlue{}, 1},
	}
	for _, tc := range cases {
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{IncomingDamage: 5})
		self := &card.CardState{Card: tc.card}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if s.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), s.Value(), tc.want)
		}
		if s.IncomingDamage() != 5-tc.want {
			t.Errorf("%s: IncomingDamage = %d, want %d", tc.card.Name(), s.IncomingDamage(), 5-tc.want)
		}
	}
}

// Tests that prevention caps at IncomingDamage when incoming is less than Defense().
func TestBrushOff_CapsAtIncoming(t *testing.T) {
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{IncomingDamage: 1})
	self := &card.CardState{Card: cards.BrushOffRed{}}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if s.Value() != 1 {
		t.Errorf("Value = %d, want 1 (capped at IncomingDamage)", s.Value())
	}
}
