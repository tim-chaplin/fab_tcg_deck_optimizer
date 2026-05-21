package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the start-of-action-phase handler destroys Sigil of Cycles and fires the
// leave-arena discard-then-draw rider against the dealt next-turn hand.
func TestSigilOfCycles_DestroysAtStartOfNextTurnAndCyclesHand(t *testing.T) {
	sigil := cards.SigilOfCyclesBlue{}
	// Hero Intel=1 + single-card deck makes turn 2's dealt hand exactly [RedAttack], so
	// the aura's Discard target is deterministic and the resulting graveyard membership
	// distinguishes "rider fired" from "rider was a no-op".
	d := deck.New(testutils.Hero{Intel: 1}, nil, []deck.Card{testutils.RedAttack{}})
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{sigil})

	if !bestLineHasRole(turn1.BestLine, ids.SigilOfCyclesBlue, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play Sigil of Cycles as Role=Attack: %+v", turn1.BestLine)
	}
	if !graveyardContains(turn2.State.Graveyard(), ids.SigilOfCyclesBlue) {
		t.Errorf("turn 2 graveyard = %v, want it to contain Sigil of Cycles (start-of-action-phase handler destroys it)",
			turn2.State.Graveyard())
	}
	if !graveyardContains(turn2.State.Graveyard(), testutils.FakeRedAttack) {
		t.Errorf("turn 2 graveyard = %v, want it to contain the RedAttack the leave-arena rider discarded",
			turn2.State.Graveyard())
	}
}
