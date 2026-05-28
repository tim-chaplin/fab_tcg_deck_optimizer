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

// Tests that Splintering Deadwood sacrifices a carried Arcane Cussing, cashing its
// 3-Runechant leave payoff on top of the Runechant the rider grants.
func TestSplinteringDeadwood_SacrificesArcaneCussingForRunechants(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingPhysicalDamage(0).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.SplinteringDeadwoodRed{}, testutils.FakeBlueResource()}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 4 {
		t.Fatalf("RunechantCount = %d, want 4 (Arcane Cussing's 3 + Splintering Deadwood's 1)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
	if got := summary.Value; got != 11 {
		t.Fatalf("Value = %d, want 11 (Splintering Deadwood 7 + 3 + 1 Runechants)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Splintering Deadwood with no aura to destroy grants no Runechant.
func TestSplinteringDeadwood_NoAuraGrantsNoRunechant(t *testing.T) {
	prior := gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.SplinteringDeadwoodRed{}, testutils.FakeBlueResource()}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 0 {
		t.Fatalf("RunechantCount = %d, want 0 (no aura to destroy, no Runechant)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
	if got := summary.Value; got != 7 {
		t.Fatalf("Value = %d, want 7 (Splintering Deadwood's attack only)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that, with two controlled auras, Splintering Deadwood cashes one on its attack leg
// and the second on its hit leg — both halves of "when this attacks or hits".
func TestSplinteringDeadwood_HitLegCashesSecondAura(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		CreateAuraFromCard(cards.ArcaneCussingYellow{}).
		SetIncomingPhysicalDamage(0).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.SplinteringDeadwoodRed{}, testutils.FakeBlueResource()}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 7 {
		t.Fatalf("RunechantCount = %d, want 7 (Arcane Cussing 3 + 2, plus 1 per leg)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
