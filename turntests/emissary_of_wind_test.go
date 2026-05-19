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

// Tests that the hand cycle fires and the granted go-again chains an arsenal play.
func TestEmissaryOfWind_CyclesHandAndChainsArsenal(t *testing.T) {
	spare := testutils.GenericAttack(0, 0)
	arsenalAttack := testutils.GenericAttack(0, 3)
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	state := gameengine.GameStateBuilder().SetArsenal(arsenalAttack).Build()
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{cards.EmissaryOfWindRed{}, spare})
	if summary.Value != 7 {
		t.Errorf("Value = %d, want 7 (Emissary 4 + arsenal 3 via go-again from cycle)", summary.Value)
	}
}

// Tests that an empty hand suppresses the granted go-again.
func TestEmissaryOfWind_EmptyHandSuppressesGoAgain(t *testing.T) {
	arsenalAttack := testutils.GenericAttack(0, 3)
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	state := gameengine.GameStateBuilder().SetArsenal(arsenalAttack).Build()
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{cards.EmissaryOfWindRed{}})
	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (Emissary alone; no cycle → no go-again → arsenal can't chain)", summary.Value)
	}
}
