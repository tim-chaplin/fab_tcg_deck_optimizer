package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that each printing's Play credits its prevention cap (3/2/1) when IncomingDamage
// has room.
func TestBrushOff_PreventsCap(t *testing.T) {
	cases := []struct {
		card sim.Card
		want int
	}{
		{BrushOffRed{}, 3},
		{BrushOffYellow{}, 2},
		{BrushOffBlue{}, 1},
	}
	for _, tc := range cases {
		s := sim.TurnState{IncomingDamage: 5}
		self := &sim.CardState{Card: tc.card}
		tc.card.Play(&s, s.Logger(), self)
		if s.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), s.Value, tc.want)
		}
		if s.IncomingDamage != 5-tc.want {
			t.Errorf("%s: IncomingDamage = %d, want %d", tc.card.Name(), s.IncomingDamage, 5-tc.want)
		}
	}
}

// Tests that prevention caps at IncomingDamage when incoming is less than Defense().
func TestBrushOff_CapsAtIncoming(t *testing.T) {
	s := sim.TurnState{IncomingDamage: 1}
	self := &sim.CardState{Card: BrushOffRed{}}
	BrushOffRed{}.Play(&s, s.Logger(), self)
	if s.Value != 1 {
		t.Errorf("Value = %d, want 1 (capped at IncomingDamage)", s.Value)
	}
}
