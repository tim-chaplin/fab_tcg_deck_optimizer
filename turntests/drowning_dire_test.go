package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Drowning Dire's Play does not flip GrantedDominate when no aura was played or
// created earlier this turn.
func TestDrowningDire_NoAuraNoDominate(t *testing.T) {
	for _, c := range []card.Card{cards.DrowningDireRed{}, cards.DrowningDireYellow{}, cards.DrowningDireBlue{}} {
		self := &card.CardState{Card: c}
		s := sim.NewTurnStateFromCards(nil, nil)
		sim.ResolveChainStep(s, s.Logger(), self)
		if self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = true without prior aura, want false", c.Name(), c.Pitch())
		}
	}
}

// Tests that an aura played earlier this turn flips GrantedDominate via AuraCreated.
func TestDrowningDire_AuraGrantsDominate(t *testing.T) {
	for _, c := range []card.Card{cards.DrowningDireRed{}, cards.DrowningDireYellow{}, cards.DrowningDireBlue{}} {
		spec := sim.TurnStateSpec{CardsPlayed: []card.Card{testutils.Aura{}}, AuraCreated: true}
		s := sim.NewTurnStateFromSpec(spec)
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if !self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = false after aura, want true", c.Name(), c.Pitch())
		}
	}
}

// Tests that the on-hit recycle moves a non-attack action card from graveyard to bottom of deck.
func TestDrowningDire_OnHitRecyclesNonAttackToBottom(t *testing.T) {
	non := testutils.GenericAction()
	deck := []card.Card{testutils.RedAttack{}}
	s := sim.NewTurnStateFromCards(deck, []card.Card{non})
	self := &card.CardState{Card: cards.DrowningDireRed{}}
	sim.ResolveChainStep(s, s.Logger(), self)
	self.BonusAttack = 2
	testutils.FireOnHitIfLikely(s, s.Logger(), self)
	if got := s.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (target appended to bottom)", got)
	}
	if top := s.Deck().PeekTop(); top != (testutils.RedAttack{}) {
		t.Errorf("deck top after recycle = %v, want RedAttack still on top (target went to bottom)", top)
	}
}
