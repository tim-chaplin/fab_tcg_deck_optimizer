package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that the start-of-action-phase handler destroys Sigil of Cycles into the graveyard.
func TestSigilOfCycles_DestroysAtStartOfNextTurn(t *testing.T) {
	sigil := cards.SigilOfCyclesBlue{}
	d := deck.New(heroes.Viserai{}, nil, nil)
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{sigil})

	if !bestLineHasRole(turn1.BestLine, ids.SigilOfCyclesBlue, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play Sigil of Cycles as Role=Attack: %+v", turn1.BestLine)
	}
	if !graveyardContains(turn2.State.Graveyard(), ids.SigilOfCyclesBlue) {
		t.Errorf("turn 2 graveyard = %v, want it to contain Sigil of Cycles (start-of-action-phase handler destroys it)",
			turn2.State.Graveyard())
	}
}
