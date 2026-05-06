package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Relentless Pursuit always marks the opposing hero on Play.
func TestRelentlessPursuit_MarksOpponent(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	(RelentlessPursuitBlue{}).Play(s, &sim.CardState{Card: RelentlessPursuitBlue{}})
	if !s.OpponentMarked {
		t.Errorf("OpponentMarked = false after Play, want true")
	}
}

// Tests that Play without a prior attack leaves the recycle clause off — the card resolves
// normally and lands in the graveyard via the chain dispatcher.
func TestRelentlessPursuit_NoRecycleWithoutPriorAttack(t *testing.T) {
	self := &sim.CardState{Card: RelentlessPursuitBlue{}}
	s := sim.NewTurnState(nil, nil)
	(RelentlessPursuitBlue{}).Play(s, self)
	if self.SkipGraveyard {
		t.Error("SkipGraveyard = true without a prior attack, want false")
	}
}

// Tests that Play after an attack-typed CardsPlayed entry recycles to the bottom of the
// deck and flips SkipGraveyard so the dispatcher will skip the graveyard append.
func TestRelentlessPursuit_RecyclesAfterPriorAttack(t *testing.T) {
	self := &sim.CardState{Card: RelentlessPursuitBlue{}}
	s := sim.NewTurnState(nil, nil)
	s.CardsPlayed = []sim.Card{testutils.GenericAttack(0, 3)}
	(RelentlessPursuitBlue{}).Play(s, self)
	if !self.SkipGraveyard {
		t.Fatal("SkipGraveyard = false after prior attack, want true")
	}
	deck := s.Deck()
	if len(deck) != 1 || deck[0] != (RelentlessPursuitBlue{}) {
		t.Errorf("deck after recycle = %v, want [RelentlessPursuitBlue{}]", deck)
	}
}
