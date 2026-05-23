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

// Tests that each elixir in the cycle grants +3{p} to the next attack.
func TestElixirCycle_BuffsNextAttack(t *testing.T) {
	for _, elixir := range []card.Card{
		cards.ClearwaterElixirRed{},
		cards.RestvineElixirRed{},
		cards.SapwoodElixirRed{},
	} {
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{elixir, testutils.RunebladeAttack{}, testutils.BluePitch{}}

		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

		if got := summary.Value; got != 3 {
			t.Errorf("%s: Value = %d, want 3 (+3{p} turns the 0-power attack into a 3-damage hit)", elixir.Name(), got)
		}
	}
}
