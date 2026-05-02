package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that mode 0 grants +2{p} to the next cost-≤N attack action card.
func TestCaptainsCall_Mode0BuffsBonusAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(1, 4)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	self := &sim.CardState{Card: CaptainsCallRed{}, Mode: 0}
	(CaptainsCallRed{}).Play(&s, self)
	if target.BonusAttack != 2 {
		t.Errorf("target.BonusAttack = %d, want 2 (mode 0 grants +2{p})", target.BonusAttack)
	}
	if target.GrantedGoAgain {
		t.Errorf("target.GrantedGoAgain = true; mode 0 should not grant go again")
	}
}

// Tests that mode 1 grants go again to the next cost-≤N attack action card.
func TestCaptainsCall_Mode1GrantsGoAgain(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(1, 4)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	self := &sim.CardState{Card: CaptainsCallRed{}, Mode: 1}
	(CaptainsCallRed{}).Play(&s, self)
	if !target.GrantedGoAgain {
		t.Errorf("target.GrantedGoAgain = false; mode 1 should grant go again")
	}
	if target.BonusAttack != 0 {
		t.Errorf("target.BonusAttack = %d; mode 1 should not add to BonusAttack", target.BonusAttack)
	}
}

// Tests that the cost cap rejects too-expensive attack action cards.
func TestCaptainsCall_BlueRejectsCostAboveZero(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(1, 4)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	self := &sim.CardState{Card: CaptainsCallBlue{}, Mode: 0}
	(CaptainsCallBlue{}).Play(&s, self)
	if target.BonusAttack != 0 {
		t.Errorf("Blue (cost cap 0) buffed a cost-1 attack; got BonusAttack = %d, want 0",
			target.BonusAttack)
	}
}
