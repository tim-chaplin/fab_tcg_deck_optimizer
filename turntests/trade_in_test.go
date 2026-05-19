package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the discard cycle pops a hand card to graveyard and fires DrawOne.
func TestTradeIn_DiscardsHandAndDraws(t *testing.T) {
	for _, c := range []card.Card{cards.TradeInRed{}, cards.TradeInYellow{}, cards.TradeInBlue{}} {
		spare := testutils.GenericAttack(0, 0)
		drawTarget := testutils.GenericAttack(0, 0)
		d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{drawTarget})
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, spare})
		if got := summary.State.CardsDrawn(); got != 1 {
			t.Errorf("%s: CardsDrawn = %d, want 1 (cycle should fire one draw)", c.Name(), got)
		}
		// The discarded card lands in graveyard alongside Trade In itself.
		grav := summary.State.Graveyard()
		if len(grav) < 2 {
			t.Errorf("%s: Graveyard = %v, want at least 2 entries (Trade In + the discarded card)", c.Name(), grav)
		}
	}
}

// Tests that an empty hand suppresses the cycle: no draw fires.
func TestTradeIn_EmptyHandSuppressesCycle(t *testing.T) {
	for _, c := range []card.Card{cards.TradeInRed{}, cards.TradeInYellow{}, cards.TradeInBlue{}} {
		drawTarget := testutils.GenericAttack(0, 0)
		d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{drawTarget})
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c})
		if got := summary.State.CardsDrawn(); got != 0 {
			t.Errorf("%s: CardsDrawn = %d, want 0 (no cycle without a spare hand card)", c.Name(), got)
		}
	}
}
