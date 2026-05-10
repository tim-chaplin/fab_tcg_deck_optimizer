package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Cadaverous Contraband's Play registers an OnHit handler.
func TestCadaverousContraband_RegistersOnHit(t *testing.T) {
	for _, c := range []sim.Card{CadaverousContrabandRed{}, CadaverousContrabandYellow{}, CadaverousContrabandBlue{}} {
		self := &sim.CardState{Card: c}
		c.Play(sim.NewTurnStateFromCards(nil, nil), self)
		if len(self.OnHit) != 1 {
			t.Errorf("%s [%d{p}]: OnHit handlers = %d, want 1", c.Name(), c.Pitch(), len(self.OnHit))
		}
	}
}

// Tests that the on-hit recycle moves a non-attack action card from graveyard to top of deck.
func TestCadaverousContraband_OnHitRecyclesNonAttackToTop(t *testing.T) {
	non := testutils.GenericAction()
	deck := []sim.Card{testutils.RedAttack{}}
	s := sim.NewTurnStateFromCards(deck, []sim.Card{non})
	self := &sim.CardState{Card: CadaverousContrabandRed{}}
	(CadaverousContrabandRed{}).Play(s, self)
	self.BonusAttack = 1
	testutils.FireOnHitIfLikely(s, self)
	gotDeck := s.Deck().PeekTopN(s.Deck().Size())
	if len(gotDeck) != 2 || gotDeck[0].(sim.Card) != sim.Card(non) {
		t.Errorf("deck after recycle = %v, want non-attack on top of RedAttack", gotDeck)
	}
	if len(s.Graveyard()) != 0 {
		t.Errorf("graveyard after recycle = %v, want empty", s.Graveyard())
	}
}

// Tests that with no non-attack action card in the graveyard, the on-hit recycle leaves the
// graveyard and deck untouched.
func TestCadaverousContraband_OnHitNoEligibleCardNoOp(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, []sim.Card{testutils.RedAttack{}})
	self := &sim.CardState{Card: CadaverousContrabandRed{}}
	(CadaverousContrabandRed{}).Play(s, self)
	testutils.FireOnHitIfLikely(s, self)
	if len(s.Graveyard()) != 1 {
		t.Errorf("graveyard size = %d, want 1 (no eligible target, no recycle)", len(s.Graveyard()))
	}
}
