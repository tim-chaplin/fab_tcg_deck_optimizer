package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Drowning Dire's Play does not flip GrantedDominate when no aura was played or
// created earlier this turn.
func TestDrowningDire_NoAuraNoDominate(t *testing.T) {
	for _, c := range []sim.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}} {
		self := &sim.CardState{Card: c}
		c.Play(sim.NewTurnStateFromCards(nil, nil), self)
		if self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = true without prior aura, want false", c.Name(), c.Pitch())
		}
	}
}

// Tests that an aura played earlier this turn flips GrantedDominate via HasPlayedOrCreatedAura.
func TestDrowningDire_AuraGrantsDominate(t *testing.T) {
	for _, c := range []sim.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}} {
		s := sim.NewTurnStateFromCards(nil, nil)
		s.CardsPlayed = []sim.Card{testutils.Aura{}}
		self := &sim.CardState{Card: c}
		c.Play(s, self)
		if !self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = false after aura, want true", c.Name(), c.Pitch())
		}
	}
}

// Tests that the on-hit recycle moves a non-attack action card from graveyard to bottom of deck.
func TestDrowningDire_OnHitRecyclesNonAttackToBottom(t *testing.T) {
	non := testutils.GenericAction()
	deck := []sim.Card{testutils.RedAttack{}}
	s := sim.NewTurnStateFromCards(deck, []sim.Card{non})
	self := &sim.CardState{Card: DrowningDireRed{}}
	(DrowningDireRed{}).Play(s, self)
	self.BonusAttack = 2
	testutils.FireOnHitIfLikely(s, self)
	gotDeck := s.Deck().PeekTopN(s.Deck().Size())
	if len(gotDeck) != 2 || gotDeck[1].(sim.Card) != sim.Card(non) {
		t.Errorf("deck after recycle = %v, want non-attack at bottom under RedAttack", gotDeck)
	}
}
