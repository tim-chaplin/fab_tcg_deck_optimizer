package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Memorial Ground with no eligible card in the graveyard leaves the deck empty.
func TestMemorialGround_NoEligibleNoOp(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.MemorialGroundRed{}})
	if s.Deck().Size() != 0 {
		t.Errorf("deck size = %d, want 0 (no eligible recycle target)", s.Deck().Size())
	}
}

// Tests that Memorial Ground recycles a cost-2-or-less attack action card from the
// graveyard to the top of the deck.
func TestMemorialGround_RecyclesEligibleAttackActionToTop(t *testing.T) {
	target := testutils.GenericAttack(2, 4)
	deck := []card.Card{testutils.BlueAttack{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCards(deck).
		SetGraveyard([]card.Card{target}).
		Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.MemorialGroundRed{}})
	if got := s.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (target moved onto the existing top)", got)
	}
	if top := s.Deck().PeekTop(); top != card.Card(target) {
		t.Errorf("deck top after recycle = %v, want %v", top, target)
	}
	if len(s.Graveyard()) != 0 {
		t.Errorf("graveyard size = %d, want 0 (target recycled out)", len(s.Graveyard()))
	}
}

// Tests that a graveyard with only an over-cost or non-attack-action card leaves Memorial
// Ground unable to recycle.
func TestMemorialGround_IgnoresIneligibleCards(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{testutils.GenericAttack(3, 5), testutils.GenericAction()}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.MemorialGroundRed{}})
	if s.Deck().Size() != 0 {
		t.Errorf("deck size = %d, want 0 (no eligible target)", s.Deck().Size())
	}
	if len(s.Graveyard()) != 2 {
		t.Errorf("graveyard size = %d, want 2 (no banish)", len(s.Graveyard()))
	}
}
