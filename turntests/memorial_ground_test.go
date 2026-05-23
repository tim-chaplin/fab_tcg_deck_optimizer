package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Memorial Ground with no eligible card in the graveyard leaves the deck empty.
func TestMemorialGround_NoEligibleNoOp(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.MemorialGroundRed{}})
	if ge.Deck().Size() != 0 {
		t.Errorf("deck size = %d, want 0 (no eligible recycle target)", ge.Deck().Size())
	}
}

// Tests that Memorial Ground recycles a cost-2-or-less attack action card from the
// graveyard to the top of the deck.
func TestMemorialGround_RecyclesEligibleAttackActionToTop(t *testing.T) {
	target := testutils.FakeRedAttack().WithCost(2)
	deck := []card.Card{testutils.FakeBlueAttack().WithCost(1)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCards(deck).
		SetGraveyard([]card.Card{target}).
		Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.MemorialGroundRed{}})
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (target moved onto the existing top)", got)
	}
	if top := ge.Deck().PeekTop(); top != card.Card(target) {
		t.Errorf("deck top after recycle = %v, want %v", top, target)
	}
	if len(ge.Graveyard()) != 0 {
		t.Errorf("graveyard size = %d, want 0 (target recycled out)", len(ge.Graveyard()))
	}
}

// Tests that a graveyard with only an over-cost or non-attack-action card leaves Memorial
// Ground unable to recycle.
func TestMemorialGround_IgnoresIneligibleCards(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{testutils.FakeRedAttack().WithCost(3), testutils.FakeRedAction()}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.MemorialGroundRed{}})
	if ge.Deck().Size() != 0 {
		t.Errorf("deck size = %d, want 0 (no eligible target)", ge.Deck().Size())
	}
	if len(ge.Graveyard()) != 2 {
		t.Errorf("graveyard size = %d, want 2 (no banish)", len(ge.Graveyard()))
	}
}
