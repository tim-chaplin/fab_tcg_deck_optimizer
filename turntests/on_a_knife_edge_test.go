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

// Tests that On a Knife Edge grants go again to the next sword attack, funding a second swing.
func TestOnAKnifeEdge_GrantsGoAgainToSwordAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, []deck.Weapon{weapons.ReapingBlade{}}, nil)
	hand := []card.Card{cards.OnAKnifeEdgeYellow{}, testutils.FakeNoGoAgainAttack{}, testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 4 {
		t.Fatalf("Value = %d, want 4 (Reaping Blade 3 + a second attack 1 — the granted go again funds it)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that On a Knife Edge's grant fizzles with no sword attack to receive it.
func TestOnAKnifeEdge_NoSwordAttackFizzles(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.OnAKnifeEdgeYellow{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 0 {
		t.Fatalf("Value = %d, want 0 (On a Knife Edge alone scores nothing)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
