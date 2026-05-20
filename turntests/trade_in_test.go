package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the discard cycle pops the spare hand card to the graveyard and that the
// drawn card lands in the arsenal slot via post-hoc promotion.
func TestTradeIn_DiscardsHandToGraveyardAndDrawsToArsenal(t *testing.T) {
	for _, c := range []card.Card{cards.TradeInRed{}, cards.TradeInYellow{}, cards.TradeInBlue{}} {
		// RedPitch sits at the deck top so Trade In's DrawOne pulls it; the four RedAttack
		// fillers cover end-of-turn refill so the arsenal slot ends up holding the
		// mid-turn-drawn RedPitch rather than a refill card.
		spare := testutils.BluePitch{}
		d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{
			testutils.RedPitch{},
			testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		})
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, spare})

		if !graveyardContains(summary.State.Graveyard(), testutils.FakeBluePitch) {
			t.Errorf("%s: graveyard = %v, want it to contain the discarded %s",
				c.Name(), summary.State.Graveyard(), spare.DisplayName())
		}
		ars := summary.State.Arsenal()
		if ars == nil || ars.ID() != testutils.FakeRedPitch {
			t.Errorf("%s: Arsenal() = %v, want the drawn RedPitch (deck top) promoted into the slot", c.Name(), ars)
		}
	}
}
