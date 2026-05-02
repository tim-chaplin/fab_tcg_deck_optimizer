package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Thrust's +3{p} buff lands on a sword attack action card.
func TestThrust_BuffsSwordAttackAction(t *testing.T) {
	swordAction := testutils.NewStubCard("SwordAction").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack, card.TypeSword))
	target := &sim.CardState{Card: swordAction}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(ThrustRed{}).Play(&s, &sim.CardState{Card: ThrustRed{}})
	if target.BonusAttack != 3 {
		t.Errorf("sword-action target BonusAttack = %d, want 3", target.BonusAttack)
	}
}

// Tests that Thrust's predicate accepts sword weapons (per the printed "sword attack"
// wording — no "action card" qualifier).
func TestThrust_AcceptsSwordWeapon(t *testing.T) {
	swordWeapon := testutils.NewStubCard("SwordWeapon").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword))
	if !(ThrustRed{}).ARTargetAllowed(swordWeapon) {
		t.Error("sword weapon should be a legal target")
	}
}

// Tests that a non-sword attack (e.g. plain Generic Action - Attack) is rejected — the
// Sword subtype is required.
func TestThrust_RejectsNonSwordAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(ThrustRed{}).Play(&s, &sim.CardState{Card: ThrustRed{}})
	if target.BonusAttack != 0 {
		t.Errorf("non-sword target BonusAttack = %d, want 0", target.BonusAttack)
	}
}
