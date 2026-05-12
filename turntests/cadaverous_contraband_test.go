package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Cadaverous Contraband's Play registers an OnHit handler.
func TestCadaverousContraband_RegistersOnHit(t *testing.T) {
	for _, c := range []card.Card{cards.CadaverousContrabandRed{}, cards.CadaverousContrabandYellow{}, cards.CadaverousContrabandBlue{}} {
		self := &card.CardState{Card: c}
		s := gameengine.NewFromCards(nil, nil)
		s.ResolveChainStep(s.Logger(), self)
		if len(self.OnHit) != 1 {
			t.Errorf("%s [%d{p}]: OnHit handlers = %d, want 1", c.Name(), c.Pitch(), len(self.OnHit))
		}
	}
}

// Tests that the on-hit recycle moves a non-attack action card from graveyard to top of deck.
func TestCadaverousContraband_OnHitRecyclesNonAttackToTop(t *testing.T) {
	non := testutils.GenericAction()
	deck := []card.Card{testutils.RedAttack{}}
	s := gameengine.NewFromCards(deck, []card.Card{non})
	self := &card.CardState{Card: cards.CadaverousContrabandRed{}}
	s.ResolveChainStep(s.Logger(), self)
	self.BonusAttack = 1
	testutils.FireOnHitIfLikely(s, s.Logger(), self)
	if got := s.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (graveyard card moved onto the existing top)", got)
	}
	if top := s.Deck().PeekTop(); top != card.Card(non) {
		t.Errorf("deck top after recycle = %v, want %v", top, non)
	}
	if len(s.Graveyard()) != 0 {
		t.Errorf("graveyard after recycle = %v, want empty", s.Graveyard())
	}
}

// Tests that with no non-attack action card in the graveyard, the on-hit recycle leaves the
// graveyard and deck untouched.
func TestCadaverousContraband_OnHitNoEligibleCardNoOp(t *testing.T) {
	s := gameengine.NewFromCards(nil, []card.Card{testutils.RedAttack{}})
	self := &card.CardState{Card: cards.CadaverousContrabandRed{}}
	s.ResolveChainStep(s.Logger(), self)
	testutils.FireOnHitIfLikely(s, s.Logger(), self)
	if len(s.Graveyard()) != 1 {
		t.Errorf("graveyard size = %d, want 1 (no eligible target, no recycle)", len(s.Graveyard()))
	}
}
