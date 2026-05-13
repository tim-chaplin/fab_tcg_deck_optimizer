package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(5).Build()}
		self := &card.CardState{Card: tc.card}
		s.ResolveChainStep(s.Logger(), self)
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
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(1).Build()}
	self := &card.CardState{Card: cards.BrushOffRed{}}
	s.ResolveChainStep(s.Logger(), self)
	if s.Value() != 1 {
		t.Errorf("Value = %d, want 1 (capped at IncomingDamage)", s.Value())
	}
}
