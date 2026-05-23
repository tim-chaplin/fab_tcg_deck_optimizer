package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that each printing's defensive-instant prevention credits its cap (3/2/1) when
// IncomingDamage has room.
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
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{tc.card}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(5).Build(), hand)
		if summary.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), summary.Value, tc.want)
		}
	}
}

// Tests that prevention caps at IncomingDamage when incoming is less than Defense().
func TestBrushOff_CapsAtIncoming(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.BrushOffRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(1).Build(), hand)
	if summary.Value != 1 {
		t.Errorf("Value = %d, want 1 (capped at IncomingDamage)", summary.Value)
	}
}
