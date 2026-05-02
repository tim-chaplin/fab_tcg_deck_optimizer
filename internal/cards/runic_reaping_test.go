package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Runic Reaping with no following attack-action target lands no riders.
func TestRunicReaping_NoNextAttackReturnsZero(t *testing.T) {
	var s sim.TurnState
	(RunicReapingRed{}).Play(&s, &sim.CardState{
		Card:          RunicReapingRed{},
		PitchedToPlay: []sim.Card{testutils.AttackWithPower{Power: 4}},
	})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false when no rider fires")
	}
}

// Tests that a Runeblade weapon as the next attack does not satisfy either rider.
func TestRunicReaping_WeaponNextDoesNotQualify(t *testing.T) {
	target := &sim.CardState{Card: testutils.RunebladeWeapon{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(RunicReapingRed{}).Play(&s, &sim.CardState{Card: RunicReapingRed{}})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if target.BonusAttack != 0 {
		t.Errorf("weapon target BonusAttack = %d, want 0 (weapons don't qualify)", target.BonusAttack)
	}
	if len(target.OnHit) != 0 {
		t.Errorf("target.OnHit = %d, want 0 (weapons don't qualify)", len(target.OnHit))
	}
}

// Tests that an attack-action target with attack-attributed funding gets +1{p} and registers the
// on-hit trigger.
func TestRunicReaping_RegistersTriggerAndGrantsPitchedAttackBonus(t *testing.T) {
	target := &sim.CardState{Card: testutils.AttackWithPower{Power: 3}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(RunicReapingRed{}).Play(&s, &sim.CardState{
		Card:          RunicReapingRed{},
		PitchedToPlay: []sim.Card{testutils.RunebladeAttack{}},
	})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0 (Runechant rider deferred to target's OnHit)", got)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1 (pitched-attack +1{p} rider)", target.BonusAttack)
	}
	if len(target.OnHit) != 1 {
		t.Fatalf("target.OnHit = %d, want 1 (on-hit Runechant rider deferred)", len(target.OnHit))
	}
}

// Tests that without an attack attributed, the +1{p} rider skips but the on-hit Runechant trigger
// still registers.
func TestRunicReaping_NoPitchedAttackSkipsBonusButRegistersTrigger(t *testing.T) {
	target := &sim.CardState{Card: testutils.AttackWithPower{Power: 4}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(RunicReapingRed{}).Play(&s, &sim.CardState{
		Card:          RunicReapingRed{},
		PitchedToPlay: []sim.Card{testutils.NonAttack{}},
	})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if target.BonusAttack != 0 {
		t.Errorf("target BonusAttack = %d, want 0 (no attack-typed card pitched)", target.BonusAttack)
	}
	if len(target.OnHit) != 1 {
		t.Fatalf("target.OnHit = %d, want 1", len(target.OnHit))
	}
}
