package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that mode 0 grants +2{p} to the next cost-≤N attack action card.
func TestCaptainsCall_Mode0BuffsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(1, 4)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	self := &card.CardState{Card: cards.CaptainsCallRed{}, Mode: 0}
	ge.ResolveChainStep(ge.Logger(), self)
	if target.BonusAttack != 2 {
		t.Errorf("target.BonusAttack = %d, want 2 (mode 0 grants +2{p})", target.BonusAttack)
	}
	if target.GrantedGoAgain {
		t.Errorf("target.GrantedGoAgain = true; mode 0 should not grant go again")
	}
}

// Tests that mode 1 grants go again to the next cost-≤N attack action card.
func TestCaptainsCall_Mode1GrantsGoAgain(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(1, 4)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	self := &card.CardState{Card: cards.CaptainsCallRed{}, Mode: 1}
	ge.ResolveChainStep(ge.Logger(), self)
	if !target.GrantedGoAgain {
		t.Errorf("target.GrantedGoAgain = false; mode 1 should grant go again")
	}
	if target.BonusAttack != 0 {
		t.Errorf("target.BonusAttack = %d; mode 1 should not add to BonusAttack", target.BonusAttack)
	}
}

// Tests that the cost cap rejects too-expensive attack action cards.
func TestCaptainsCall_BlueRejectsCostAboveZero(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(1, 4)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	self := &card.CardState{Card: cards.CaptainsCallBlue{}, Mode: 0}
	ge.ResolveChainStep(ge.Logger(), self)
	if target.BonusAttack != 0 {
		t.Errorf("Blue (cost cap 0) buffed a cost-1 attack; got BonusAttack = %d, want 0",
			target.BonusAttack)
	}
}
