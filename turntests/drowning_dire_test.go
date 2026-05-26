package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Drowning Dire's Play does not flip GrantedDominate when no aura was played or
// created earlier this turn.
func TestDrowningDire_NoAuraNoDominate(t *testing.T) {
	for _, c := range []card.Card{cards.DrowningDireRed{}, cards.DrowningDireYellow{}, cards.DrowningDireBlue{}} {
		pc := &card.CardState{Card: c}
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), pc)
		if pc.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = true without prior aura, want false", c.Name(), c.Pitch())
		}
	}
}

// Tests that an aura played earlier this turn flips GrantedDominate via AuraCreated.
func TestDrowningDire_AuraGrantsDominate(t *testing.T) {
	for _, c := range []card.Card{cards.DrowningDireRed{}, cards.DrowningDireYellow{}, cards.DrowningDireBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetCardsPlayed([]card.Card{testutils.FakeRedAura()}).
			SetAuraCreated(true).
			Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if !pc.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = false after aura, want true", c.Name(), c.Pitch())
		}
	}
}

// Tests that the on-hit recycle moves a non-attack action card from graveyard to bottom of deck.
func TestDrowningDire_OnHitRecyclesNonAttackToBottom(t *testing.T) {
	non := testutils.FakeRedAction()
	deck := []card.Card{testutils.FakeRedAttack()}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCards(deck).
		SetGraveyard([]card.Card{non}).
		Build()}
	pc := &card.CardState{Card: cards.DrowningDireRed{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	pc.BonusAttack = 2
	testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("deck size after recycle = %d, want 2 (target appended to bottom)", got)
	}
	if top := ge.Deck().PeekTop(); top.DisplayName() != "FakeRedAttack" {
		t.Errorf("deck top after recycle = %v, want RedAttack still on top (target went to bottom)", top)
	}
}
