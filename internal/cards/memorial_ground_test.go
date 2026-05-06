package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Memorial Ground with no eligible card in the graveyard leaves the deck empty.
func TestMemorialGround_NoEligibleNoOp(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	(MemorialGroundRed{}).Play(s, &sim.CardState{Card: MemorialGroundRed{}})
	if len(s.Deck()) != 0 {
		t.Errorf("deck size = %d, want 0 (no eligible recycle target)", len(s.Deck()))
	}
}

// Tests that Memorial Ground recycles a cost-2-or-less attack action card from the
// graveyard to the top of the deck.
func TestMemorialGround_RecyclesEligibleAttackActionToTop(t *testing.T) {
	target := testutils.GenericAttack(2, 4)
	deck := []sim.Card{testutils.BlueAttack{}}
	s := sim.NewTurnState(deck, []sim.Card{target})
	(MemorialGroundRed{}).Play(s, &sim.CardState{Card: MemorialGroundRed{}})
	gotDeck := s.Deck()
	if len(gotDeck) != 2 || gotDeck[0] != sim.Card(target) {
		t.Errorf("deck after recycle = %v, want target on top of BlueAttack", gotDeck)
	}
	if len(s.Graveyard()) != 0 {
		t.Errorf("graveyard size = %d, want 0 (target recycled out)", len(s.Graveyard()))
	}
}

// Tests that a graveyard with only an over-cost or non-attack-action card leaves Memorial
// Ground unable to recycle.
func TestMemorialGround_IgnoresIneligibleCards(t *testing.T) {
	s := sim.NewTurnState(nil, []sim.Card{testutils.GenericAttack(3, 5), testutils.GenericAction()})
	(MemorialGroundRed{}).Play(s, &sim.CardState{Card: MemorialGroundRed{}})
	if len(s.Deck()) != 0 {
		t.Errorf("deck size = %d, want 0 (no eligible target)", len(s.Deck()))
	}
	if len(s.Graveyard()) != 2 {
		t.Errorf("graveyard size = %d, want 2 (no banish)", len(s.Graveyard()))
	}
}
