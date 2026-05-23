package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that Regain Composure buffs the next attack: a 0-power attack lands for 1 only
// because of the +1{p} grant.
func TestRegainComposure_BuffsNextAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.RegainComposureBlue{}, testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 1 {
		t.Fatalf("Value = %d, want 1 (Regain Composure's +1{p} turns the 0-power attack into a 1-damage hit)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Regain Composure's grant fizzles when no attack follows it.
func TestRegainComposure_NoAttackGrantsNothing(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.RegainComposureBlue{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 0 {
		t.Fatalf("Value = %d, want 0 (no attack to buff)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Regain Composure's +1{p} reaches a weapon attack: a Reaping Blade swing
// (base 3) misses on its own but lands for 4 once buffed.
func TestRegainComposure_BuffsWeaponAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, []deck.Weapon{weapons.ReapingBlade{}}, nil)
	hand := []card.Card{cards.RegainComposureBlue{}, testutils.FakeBlueResource()}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 4 {
		t.Fatalf("Value = %d, want 4 (Reaping Blade base 3 misses; Regain Composure's +1{p} makes it a 4-damage hit)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
