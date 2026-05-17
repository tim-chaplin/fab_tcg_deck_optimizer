package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(5).Build()}
		self := &card.CardState{Card: tc.card}
		ge.ResolveChainStep(ge.Logger(), self)
		if ge.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), ge.Value(), tc.want)
		}
		if ge.RemainingUnblockedDamage() != 5-tc.want {
			t.Errorf("%s: RemainingUnblockedDamage = %d, want %d", tc.card.Name(), ge.RemainingUnblockedDamage(), 5-tc.want)
		}
	}
}

// Tests that prevention caps at IncomingDamage when incoming is less than Defense().
func TestBrushOff_CapsAtIncoming(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(1).Build()}
	self := &card.CardState{Card: cards.BrushOffRed{}}
	ge.ResolveChainStep(ge.Logger(), self)
	if ge.Value() != 1 {
		t.Errorf("Value = %d, want 1 (capped at IncomingDamage)", ge.Value())
	}
}
