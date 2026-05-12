package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Relentless Pursuit always marks the opposing hero on Play.
func TestRelentlessPursuit_MarksOpponent(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: cards.RelentlessPursuitBlue{}})
	if !s.OpponentMarked() {
		t.Errorf("OpponentMarked = false after Play, want true")
	}
}

// Tests that Play without a prior attack leaves the recycle clause off — the card resolves
// normally and lands in the graveyard via the chain dispatcher.
func TestRelentlessPursuit_NoRecycleWithoutPriorAttack(t *testing.T) {
	self := &card.CardState{Card: cards.RelentlessPursuitBlue{}}
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), self)
	if self.SkipGraveyard {
		t.Error("SkipGraveyard = true without a prior attack, want false")
	}
}

// Tests that Play after an attack-typed CardsPlayed entry recycles to the bottom of the
// deck and flips SkipGraveyard so the dispatcher will skip the graveyard append.
func TestRelentlessPursuit_RecyclesAfterPriorAttack(t *testing.T) {
	self := &card.CardState{Card: cards.RelentlessPursuitBlue{}}
	s := sim.NewTurnStateFromCards(nil, nil)
	s.SetCardsPlayed([]card.Card{testutils.GenericAttack(0, 3)})
	sim.ResolveChainStep(s, s.Logger(), self)
	if !self.SkipGraveyard {
		t.Fatal("SkipGraveyard = false after prior attack, want true")
	}
	if got := s.Deck().Size(); got != 1 {
		t.Errorf("deck size after recycle = %d, want 1 (Relentless Pursuit went onto an empty deck)", got)
	}
	if top := s.Deck().PeekTop(); top != (cards.RelentlessPursuitBlue{}) {
		t.Errorf("deck top after recycle = %v, want RelentlessPursuitBlue{}", top)
	}
}
