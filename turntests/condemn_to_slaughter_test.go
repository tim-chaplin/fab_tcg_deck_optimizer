package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Condemn to Slaughter sacrifices a carried Arcane Cussing, cashing its
// 3-Runechant leave payoff during our turn.
func TestCondemnToSlaughter_SacrificesArcaneCussing(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingDamage(0).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{cards.CondemnToSlaughterRed{}, testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 3 {
		t.Fatalf("RunechantCount = %d, want 3 (Arcane Cussing cashed by the aura sacrifice)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Condemn to Slaughter buffs the next Runeblade attack: a 0-power Runeblade
// attack lands for 1 only because of the Blue variant's +1{p}.
func TestCondemnToSlaughter_BuffsNextRunebladeAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{cards.CondemnToSlaughterBlue{}, testutils.RunebladeAttack{}, testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)

	if got := summary.Value; got != 1 {
		t.Fatalf("Value = %d, want 1 (Condemn's +1{p} turns the 0-power Runeblade attack into a 1-damage hit)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
