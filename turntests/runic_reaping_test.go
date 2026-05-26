package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Runic Reaping with no following attack-action target lands no riders.
func TestRunicReaping_NoNextAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	(cards.RunicReapingRed{}).Play(ge, ge.Logger(), &card.CardState{
		Card:      cards.RunicReapingRed{},
		Ephemeral: card.Ephemeral{PitchedToPlay: []card.Card{testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}},
	})
	if got := ge.Value(); got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if ge.AuraCreated() {
		t.Fatalf("AuraCreated should stay false when no rider fires")
	}
}

// Tests that a Runeblade weapon as the next attack does not satisfy either rider.
func TestRunicReaping_WeaponNextDoesNotQualify(t *testing.T) {
	target := &card.CardState{Card: testutils.FakeWeaponSwing().WithTypes(card.TypeRuneblade)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.RunicReapingRed{}})
	if got := ge.Value(); got != 0 {
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
	target := &card.CardState{Card: testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	(cards.RunicReapingRed{}).Play(ge, ge.Logger(), &card.CardState{
		Card:      cards.RunicReapingRed{},
		Ephemeral: card.Ephemeral{PitchedToPlay: []card.Card{testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}},
	})
	if got := ge.Value(); got != 0 {
		t.Fatalf("Play() = %d, want 0 (Runechant rider deferred to target'ge OnHit)", got)
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
	target := &card.CardState{Card: testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	(cards.RunicReapingRed{}).Play(ge, ge.Logger(), &card.CardState{
		Card:      cards.RunicReapingRed{},
		Ephemeral: card.Ephemeral{PitchedToPlay: []card.Card{testutils.FakeRedAction()}},
	})
	if got := ge.Value(); got != 0 {
		t.Fatalf("Play() = %d, want 0", got)
	}
	if target.BonusAttack != 0 {
		t.Errorf("target BonusAttack = %d, want 0 (no attack-typed card pitched)", target.BonusAttack)
	}
	if len(target.OnHit) != 1 {
		t.Fatalf("target.OnHit = %d, want 1", len(target.OnHit))
	}
}
