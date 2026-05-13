package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Relentless Pursuit always marks the opposing hero on Play.
func TestRelentlessPursuit_MarksOpponent(t *testing.T) {
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.RelentlessPursuitBlue{}})
	if !s.OpponentMarked() {
		t.Errorf("OpponentMarked = false after Play, want true")
	}
}

// Tests that Play without a prior attack leaves the recycle clause off — the card
// resolves normally and the deck stays empty.
func TestRelentlessPursuit_NoRecycleWithoutPriorAttack(t *testing.T) {
	self := &card.CardState{Card: cards.RelentlessPursuitBlue{}}
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), self)
	if got := s.Deck().Size(); got != 0 {
		t.Errorf("deck size = %d, want 0 (no recycle without prior attack)", got)
	}
}

// Tests that Play after an attack-typed CardsPlayed entry recycles to the bottom of the
// deck via RecycleToDeckBottom.
func TestRelentlessPursuit_RecyclesAfterPriorAttack(t *testing.T) {
	self := &card.CardState{Card: cards.RelentlessPursuitBlue{}}
	s := gameengine.New()
	s.SetCardsPlayed([]card.Card{testutils.GenericAttack(0, 3)})
	s.ResolveChainStep(s.Logger(), self)
	if got := s.Deck().Size(); got != 1 {
		t.Errorf("deck size after recycle = %d, want 1 (Relentless Pursuit went onto an empty deck)", got)
	}
	if top := s.Deck().PeekTop(); top != (cards.RelentlessPursuitBlue{}) {
		t.Errorf("deck top after recycle = %v, want RelentlessPursuitBlue{}", top)
	}
}
