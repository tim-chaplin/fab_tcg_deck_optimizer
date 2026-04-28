package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestRunicReaping_NoNextAttackReturnsZero pins the no-target case: with no Runeblade attack
// action card following Runic Reaping, neither rider lands — no BonusAttack, no trigger
// registered, AuraCreated stays false.
func TestRunicReaping_NoNextAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{Pitched: []sim.Card{testutils.AttackWithPower{Power: 4}}}
	(RunicReapingRed{}).Play(&s, &sim.CardState{Card: RunicReapingRed{}})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if len(s.EphemeralAttackTriggers) != 0 {
		t.Fatalf("EphemeralAttackTriggers = %d, want 0 when no target matches", len(s.EphemeralAttackTriggers))
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false when no rider fires")
	}
}

// TestRunicReaping_WeaponNextDoesNotQualify: a Runeblade weapon swing later in the turn is
// not an attack action card, so neither rider applies.
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
	if len(s.EphemeralAttackTriggers) != 0 {
		t.Errorf("EphemeralAttackTriggers = %d, want 0 (weapons don't qualify)", len(s.EphemeralAttackTriggers))
	}
}

// TestRunicReaping_RegistersTriggerAndGrantsPitchedAttackBonus: with a Runeblade attack
// action target queued and an attack-typed card pitched, Runic Reaping (a) sets +1
// BonusAttack on the target so EffectiveAttack folds it into hit-likelihood, and (b)
// registers an EphemeralAttackTrigger that defers the on-hit Runechant rider until the
// target's full resolution. Play returns 0 — the trigger handler's damage routes back via
// SourceIndex when the target lands.
func TestRunicReaping_RegistersTriggerAndGrantsPitchedAttackBonus(t *testing.T) {
	target := &sim.CardState{Card: testutils.AttackWithPower{Power: 3}}
	s := sim.TurnState{
		CardsRemaining: []*sim.CardState{target},
		Pitched:        []sim.Card{testutils.RunebladeAttack{}},
	}
	(RunicReapingRed{}).Play(&s, &sim.CardState{Card: RunicReapingRed{}})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0 (rider fires through ephemeral trigger after target's resolution)", got)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1 (pitched-attack +1{p} rider)", target.BonusAttack)
	}
	if len(s.EphemeralAttackTriggers) != 1 {
		t.Fatalf("EphemeralAttackTriggers = %d, want 1 (on-hit Runechant rider deferred)", len(s.EphemeralAttackTriggers))
	}
}

// TestRunicReaping_NoPitchedAttackSkipsBonusButRegistersTrigger: without an attack-typed
// card in Pitched the +1{p} rider doesn't fire, but the on-hit Runechant rider still
// registers — the two riders are independent.
func TestRunicReaping_NoPitchedAttackSkipsBonusButRegistersTrigger(t *testing.T) {
	target := &sim.CardState{Card: testutils.AttackWithPower{Power: 4}}
	s := sim.TurnState{
		CardsRemaining: []*sim.CardState{target},
		Pitched:        []sim.Card{testutils.NonAttack{}},
	}
	(RunicReapingRed{}).Play(&s, &sim.CardState{Card: RunicReapingRed{}})
	if got := s.Value; got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if target.BonusAttack != 0 {
		t.Errorf("target BonusAttack = %d, want 0 (no attack-typed card pitched)", target.BonusAttack)
	}
	if len(s.EphemeralAttackTriggers) != 1 {
		t.Fatalf("EphemeralAttackTriggers = %d, want 1", len(s.EphemeralAttackTriggers))
	}
}
