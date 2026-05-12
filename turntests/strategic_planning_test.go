package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Strategic Planning queues a TriggerEndOfTurn keyed to itself, not a Ponder.
func TestStrategicPlanning_QueuesEndOfTurnTrigger(t *testing.T) {
	for _, c := range []card.Card{cards.StrategicPlanningRed{}, cards.StrategicPlanningYellow{}, cards.StrategicPlanningBlue{}} {
		s := sim.NewTurnStateFromCards(nil, nil)
		sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: c})
		matching := 0
		for _, tr := range s.Triggers() {
			if tr.TriggerType == sim.TriggerEndOfTurn && tr.Source == c {
				matching++
			}
		}
		if matching != 1 {
			t.Errorf("%s [%d{p}]: end-of-turn triggers keyed to self = %d, want 1", c.Name(), c.Pitch(), matching)
		}
		if s.Ponders() != 0 {
			t.Errorf("%s [%d{p}]: Ponders = %d, want 0 (trigger is standalone, not a Ponder token)", c.Name(), c.Pitch(), s.Ponders())
		}
	}
}

// Tests that an eligible graveyard action is recycled to the bottom of the deck.
func TestStrategicPlanning_RecyclesEligibleActionToBottom(t *testing.T) {
	target := testutils.GenericAction()
	deck := []card.Card{testutils.BlueAttack{}}
	s := sim.NewTurnStateFromCards(deck, []card.Card{target})
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: cards.StrategicPlanningRed{}})
	if got := s.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (target appended to bottom)", got)
	}
	if top := s.Deck().PeekTop(); top != (testutils.BlueAttack{}) {
		t.Errorf("deck top after recycle = %v, want BlueAttack still on top (target went to bottom)", top)
	}
}
