package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the hand cycle fires one DrawOne when there's a spare card to send to the
// deck bottom.
func TestEmissaryOfMoon_CyclesHandAndDraws(t *testing.T) {
	spare := testutils.GenericAttack(0, 0)
	drawTarget := testutils.GenericAttack(0, 0)
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{drawTarget})
	summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{cards.EmissaryOfMoonRed{}, spare})
	if got := summary.State.CardsDrawn(); got != 1 {
		t.Errorf("CardsDrawn = %d, want 1 (cycle should fire one draw)", got)
	}
}

// Tests that an empty hand suppresses the cycle: no DrawOne fires.
func TestEmissaryOfMoon_EmptyHandSuppressesCycle(t *testing.T) {
	drawTarget := testutils.GenericAttack(0, 0)
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{drawTarget})
	summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{cards.EmissaryOfMoonRed{}})
	if got := summary.State.CardsDrawn(); got != 0 {
		t.Errorf("CardsDrawn = %d, want 0 (no cycle without a spare hand card)", got)
	}
}
