package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the hand cycle pops the spare hand card to the deck bottom.
func TestScourTheBattlescape_CyclesHandToDeckBottom(t *testing.T) {
	for _, c := range []card.Card{cards.ScourTheBattlescapeRed{}, cards.ScourTheBattlescapeYellow{}, cards.ScourTheBattlescapeBlue{}} {
		// BluePitch's ID sorts after Scour's, so the chain runner's PopHandAt(0) inside
		// Scour's Play pops BluePitch. Five-card filler deck keeps BluePitch at the
		// bottom past Scour's mid-turn DrawOne (1) plus end-of-turn refill (4 more).
		spare := testutils.BluePitch{}
		d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{
			testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
			testutils.RedAttack{}, testutils.RedAttack{},
		})
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, spare})
		if got := summary.State.Deck().NameCounts()[spare.DisplayName()]; got != 1 {
			t.Errorf("%s: deck contains %d copies of the cycled %s, want 1", c.Name(), got, spare.DisplayName())
		}
	}
}
