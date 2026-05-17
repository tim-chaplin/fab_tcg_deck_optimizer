package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Tests that Strategic Planning queues a TriggerEndOfTurn keyed to itself, not a Ponder.
func TestStrategicPlanning_QueuesEndOfTurnTrigger(t *testing.T) {
	for _, c := range []card.Card{cards.StrategicPlanningRed{}, cards.StrategicPlanningYellow{}, cards.StrategicPlanningBlue{}} {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		matching := 0
		for _, tr := range ge.Triggers() {
			if tr.TriggerType() == triggertype.EndOfTurn && tr.CardName() == c.DisplayName() {
				matching++
			}
		}
		if matching != 1 {
			t.Errorf("%s [%d{p}]: end-of-turn triggers keyed to self = %d, want 1", c.Name(), c.Pitch(), matching)
		}
		if ge.PonderCount() != 0 {
			t.Errorf("%s [%d{p}]: Ponders = %d, want 0 (trigger is standalone, not a Ponder token)", c.Name(), c.Pitch(), ge.PonderCount())
		}
	}
}

// Tests that an eligible graveyard action is recycled to the bottom of the deck.
func TestStrategicPlanning_RecyclesEligibleActionToBottom(t *testing.T) {
	target := testutils.GenericAction()
	deck := []card.Card{testutils.BlueAttack{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCards(deck).
		SetGraveyard([]card.Card{target}).
		Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.StrategicPlanningRed{}})
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (target appended to bottom)", got)
	}
	if top := ge.Deck().PeekTop(); top != (testutils.BlueAttack{}) {
		t.Errorf("deck top after recycle = %v, want BlueAttack still on top (target went to bottom)", top)
	}
}
