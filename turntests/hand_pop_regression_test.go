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

// Tests that a card moved out of hand by an alt-cost Play (Seek Horizon's hand-on-top
// rider) can't be replayed later in the same chain.
func TestChainRunner_AltCostPoppedCardCannotPhantomPlay(t *testing.T) {
	bigAttack := testutils.FakeRedAttack()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(
		d,
		gameengine.GameStateBuilder().SetIncomingDamage(0).Build(),
		[]card.Card{cards.SeekHorizonRed{}, bigAttack},
	)
	if summary.Value > 5 {
		t.Errorf("Value = %d, want ≤ 5 (best legal line is GenericAttack alone for 5)", summary.Value)
	}
}
