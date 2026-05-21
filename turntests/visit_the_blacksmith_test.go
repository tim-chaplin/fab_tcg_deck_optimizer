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

// Tests that Visit the Blacksmith grants +1{p} to the next sword attack — a Reaping Blade
// swing (base 3) lands for 4.
func TestVisitTheBlacksmith_BuffsNextSwordAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, []deck.Weapon{weapons.ReapingBlade{}}, fillerDeck())
	hand := []card.Card{cards.VisitTheBlacksmithBlue{}, testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 4 {
		t.Fatalf("Value = %d, want 4 (Reaping Blade swings for 3, lifted to 4 by the +1{p} grant)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Visit the Blacksmith's grant fizzles with no sword attack to buff.
func TestVisitTheBlacksmith_NoSwordAttackFizzles(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{cards.VisitTheBlacksmithBlue{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 0 {
		t.Fatalf("Value = %d, want 0 (no sword attack to buff)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
