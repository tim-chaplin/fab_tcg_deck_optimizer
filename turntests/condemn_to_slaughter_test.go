package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestCondemnToSlaughter_NoNextAttackReturnsZero(t *testing.T) {
	// No Runeblade attack follows → rider doesn't fire, Play returns 0.
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.CondemnToSlaughterRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 when CardsRemaining is empty", got)
	}
}

func TestCondemnToSlaughter_NextAttackActionTriggers(t *testing.T) {
	// A Runeblade attack action card in CardsRemaining picks up +N{p} on its BonusAttack;
	// Play returns 0 (the +N attributes to the buffed attack, not Condemn).
	cases := []struct {
		c card.Card
		n int
	}{
		{notimplemented.CondemnToSlaughterRed{}, 3},
		{notimplemented.CondemnToSlaughterYellow{}, 2},
		{notimplemented.CondemnToSlaughterBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.RunebladeAttack{}}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.n {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.n)
		}
	}
}

func TestCondemnToSlaughter_WeaponCountsAsNextAttack(t *testing.T) {
	// Unlike Runic Reaping, Condemn'ge rider accepts weapon swings as the "next attack."
	target := &card.CardState{Card: testutils.RunebladeWeapon{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.CondemnToSlaughterRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3 (weapon should qualify)", target.BonusAttack)
	}
}

func TestCondemnToSlaughter_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	// A Generic attack-action card later in the chain doesn't satisfy the Runeblade-only rider.
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.NonRunebladeAttack{}}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.CondemnToSlaughterRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't qualify)", got)
	}
}
