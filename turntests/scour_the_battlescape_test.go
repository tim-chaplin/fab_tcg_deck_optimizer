package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Scour plays for its printed attack and fires the hand cycle (DrawOne)
// when there's a spare hand card to send to the deck bottom.
func TestScourTheBattlescape_CyclesHandAndDraws(t *testing.T) {
	for _, c := range []card.Card{cards.ScourTheBattlescapeRed{}, cards.ScourTheBattlescapeYellow{}, cards.ScourTheBattlescapeBlue{}} {
		cycled := testutils.GenericAttack(0, 0)
		drawTarget := testutils.GenericAttack(0, 0)
		d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{drawTarget})
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, cycled})
		if summary.Value != c.Attack() {
			t.Errorf("%s: Value = %d, want %d (printed attack)", c.Name(), summary.Value, c.Attack())
		}
		if got := summary.State.CardsDrawn(); got != 1 {
			t.Errorf("%s: CardsDrawn = %d, want 1 (cycle should fire one draw)", c.Name(), got)
		}
	}
}

// Tests that an empty hand suppresses the cycle: no DrawOne fires.
func TestScourTheBattlescape_EmptyHandSuppressesCycle(t *testing.T) {
	for _, c := range []card.Card{cards.ScourTheBattlescapeRed{}, cards.ScourTheBattlescapeYellow{}, cards.ScourTheBattlescapeBlue{}} {
		drawTarget := testutils.GenericAttack(0, 0)
		d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{drawTarget})
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c})
		if summary.Value != c.Attack() {
			t.Errorf("%s: Value = %d, want %d (printed attack)", c.Name(), summary.Value, c.Attack())
		}
		if got := summary.State.CardsDrawn(); got != 0 {
			t.Errorf("%s: CardsDrawn = %d, want 0 (no cycle without a spare hand card)", c.Name(), got)
		}
	}
}
