package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// swordWeapon is a minimal Sword-typed weapon stub for the lookahead.
type swordWeapon struct{}

func (swordWeapon) ID() ids.CardID          { return ids.InvalidCard }
func (swordWeapon) Name() string            { return "SwordWeapon" }
func (swordWeapon) Cost(*sim.TurnState) int { return 0 }
func (swordWeapon) Pitch() int              { return 0 }
func (swordWeapon) Attack() int             { return 0 }
func (swordWeapon) Defense() int            { return 0 }
func (swordWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword)
}
func (swordWeapon) GoAgain() bool                       { return false }
func (swordWeapon) Play(*sim.TurnState, *sim.CardState) {}

// swordAttack is a minimal Sword-typed attack action card stub.
type swordAttack struct{}

func (swordAttack) ID() ids.CardID          { return ids.InvalidCard }
func (swordAttack) Name() string            { return "SwordAttack" }
func (swordAttack) Cost(*sim.TurnState) int { return 0 }
func (swordAttack) Pitch() int              { return 0 }
func (swordAttack) Attack() int             { return 0 }
func (swordAttack) Defense() int            { return 0 }
func (swordAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack, card.TypeSword)
}
func (swordAttack) GoAgain() bool                       { return false }
func (swordAttack) Play(*sim.TurnState, *sim.CardState) {}

// Tests that Visit the Blacksmith with no following attack lands no buff.
func TestVisitTheBlacksmith_NoNextAttack(t *testing.T) {
	var s sim.TurnState
	(VisitTheBlacksmithBlue{}).Play(&s, &sim.CardState{Card: VisitTheBlacksmithBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// Tests that Visit the Blacksmith grants +1{p} to the next sword-typed attack-action card.
func TestVisitTheBlacksmith_GrantsToSwordAttackAction(t *testing.T) {
	target := &sim.CardState{Card: swordAttack{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(VisitTheBlacksmithBlue{}).Play(&s, &sim.CardState{Card: VisitTheBlacksmithBlue{}})
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1", target.BonusAttack)
	}
}

// Tests that Visit the Blacksmith grants +1{p} to the next sword-typed weapon swing.
func TestVisitTheBlacksmith_GrantsToSwordWeapon(t *testing.T) {
	target := &sim.CardState{Card: swordWeapon{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(VisitTheBlacksmithBlue{}).Play(&s, &sim.CardState{Card: VisitTheBlacksmithBlue{}})
	if target.BonusAttack != 1 {
		t.Errorf("weapon target BonusAttack = %d, want 1", target.BonusAttack)
	}
}

// Tests that Visit the Blacksmith does not grant to a non-sword attack-action card.
func TestVisitTheBlacksmith_SkipsNonSwordAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(VisitTheBlacksmithBlue{}).Play(&s, &sim.CardState{Card: VisitTheBlacksmithBlue{}})
	if target.BonusAttack != 0 {
		t.Errorf("non-sword target BonusAttack = %d, want 0", target.BonusAttack)
	}
}
