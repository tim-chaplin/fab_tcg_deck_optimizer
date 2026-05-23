package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that Bloodspill Invocation, popped by an attack action card hitting, creates its N
// Runechants.
func TestBloodspillInvocation_AttackActionHitCreatesRunechants(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.BloodspillInvocationRed{}, 3},
		{cards.BloodspillInvocationYellow{}, 2},
		{cards.BloodspillInvocationBlue{}, 1},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			CreateAuraFromCard(tc.c).
			Build()}
		ge.FireTriggers(triggertype.Hit, testutils.FakeRedAttack().WithTypes(card.TypeRuneblade))
		if got := ge.RunechantCount(); got != tc.n {
			t.Errorf("%s: RunechantCount = %d, want %d (popped by an attack action hit)", tc.c.Name(), got, tc.n)
		}
	}
}

// Tests that a weapon swing hitting does not pop Bloodspill Invocation — its trigger is
// gated on attack action cards, so the aura survives a weapon hit.
func TestBloodspillInvocation_WeaponHitDoesNotPop(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.BloodspillInvocationRed{}).
		Build()}
	ge.FireTriggers(triggertype.Hit, testutils.FakeWeaponSwing().WithTypes(card.TypeRuneblade))
	if got := ge.RunechantCount(); got != 0 {
		t.Errorf("RunechantCount = %d, want 0 (a weapon hit doesn't trigger Bloodspill)", got)
	}
	if got := len(ge.Auras()); got != 1 {
		t.Errorf("Auras = %d, want 1 (Bloodspill survives a weapon hit)", got)
	}
}

// Tests that Bloodspill Invocation popped by DamageTaken is destroyed without creating
// Runechants — its create clause is bound to the attack-action-hit destroy only.
func TestBloodspillInvocation_DamageTakenDestroysWithoutRunechants(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.BloodspillInvocationRed{}).
		Build()}
	ge.FireTriggers(triggertype.DamageTaken, nil)
	if got := ge.RunechantCount(); got != 0 {
		t.Errorf("RunechantCount = %d, want 0 (DamageTaken destroy creates nothing)", got)
	}
	if got := len(ge.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Bloodspill destroyed by DamageTaken)", got)
	}
}
