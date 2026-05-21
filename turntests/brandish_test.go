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

// Tests that Brandish on hit grants +1{p} to the next weapon attack.
func TestBrandish_OnHitBuffsNextWeaponAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, []deck.Weapon{weapons.ReapingBlade{}}, fillerDeck())
	hand := []card.Card{cards.BrandishBlue{}, testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 5 {
		t.Fatalf("Value = %d, want 5 (Brandish [B] hits for 1; its rider lifts the Reaping Blade swing 3 -> 4)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Brandish's rider fizzles when no weapon swing follows.
func TestBrandish_NoWeaponRiderFizzles(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{cards.BrandishBlue{}, testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 1 {
		t.Fatalf("Value = %d, want 1 (Brandish [B] hits for 1; no weapon swing for the rider to buff)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
