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

// Tests that a defending Rally card discards a Held card for the printed +3{d} — blocking
// all 5 incoming — and blocks only its base 2 when the hand holds nothing to discard.
func TestRallyCycle_BlockDiscardsHeldCardForBonusDefense(t *testing.T) {
	for _, rally := range []card.Card{cards.RallyTheCoastGuardRed{}, cards.RallyTheRearguardRed{}} {
		d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
		junk := testutils.GenericAttackPitch(0, 0, 0)

		withSpare := sim.EvalOneTurnForTesting(d,
			gameengine.GameStateBuilder().SetIncomingDamage(5).Build(),
			[]card.Card{rally, junk, junk})
		if got := withSpare.Value; got != 5 {
			t.Errorf("%s with a spare card: Value = %d, want 5 (blocks 2, discards a Held card for +3)\nBestLine: %s",
				rally.Name(), got, formatBestLine(withSpare.BestLine))
		}

		alone := sim.EvalOneTurnForTesting(d,
			gameengine.GameStateBuilder().SetIncomingDamage(5).Build(),
			[]card.Card{rally})
		if got := alone.Value; got != 2 {
			t.Errorf("%s alone: Value = %d, want 2 (base block, nothing to discard)\nBestLine: %s",
				rally.Name(), got, formatBestLine(alone.BestLine))
		}
	}
}
