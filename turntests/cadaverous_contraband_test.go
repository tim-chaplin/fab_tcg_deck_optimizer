package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Cadaverous Contraband's Play registers an OnHit handler.
func TestCadaverousContraband_RegistersOnHit(t *testing.T) {
	for _, c := range []card.Card{cards.CadaverousContrabandRed{}, cards.CadaverousContrabandYellow{}, cards.CadaverousContrabandBlue{}} {
		pc := &card.CardState{Card: c}
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), pc)
		if len(pc.OnHit) != 1 {
			t.Errorf("%s [%d{p}]: OnHit handlers = %d, want 1", c.Name(), c.Pitch(), len(pc.OnHit))
		}
	}
}

// Tests that the on-hit recycle moves a non-attack action card from graveyard to top of deck.
func TestCadaverousContraband_OnHitRecyclesNonAttackToTop(t *testing.T) {
	non := testutils.FakeRedAction()
	deck := []card.Card{testutils.FakeRedAttack()}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCards(deck).
		SetGraveyard([]card.Card{non}).
		Build()}
	pc := &card.CardState{Card: cards.CadaverousContrabandRed{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	pc.BonusAttack = 1
	testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (graveyard card moved onto the existing top)", got)
	}
	if top := ge.Deck().PeekTop(); top != card.Card(non) {
		t.Errorf("deck top after recycle = %v, want %v", top, non)
	}
	if len(ge.Graveyard()) != 0 {
		t.Errorf("graveyard after recycle = %v, want empty", ge.Graveyard())
	}
}

// Tests that with no non-attack action card in the graveyard, the on-hit recycle leaves the
// graveyard and deck untouched.
func TestCadaverousContraband_OnHitNoEligibleCardNoOp(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{testutils.FakeRedAttack()}).Build()}
	pc := &card.CardState{Card: cards.CadaverousContrabandRed{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
	if len(ge.Graveyard()) != 1 {
		t.Errorf("graveyard size = %d, want 1 (no eligible target, no recycle)", len(ge.Graveyard()))
	}
}
