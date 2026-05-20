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

// Tests that a carried Arcane Cussing destroyed by unblocked incoming damage pays out
// nothing — it leaves the arena during the defense phase, not our turn.
func TestArcaneCussing_PoppedUnblockedCreatesNoRunechants(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingDamage(3).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 0 {
		t.Fatalf("RunechantCount = %d, want 0 (Arcane Cussing popped by DamageTaken, not our turn)", got)
	}
	if got := len(summary.State.Auras()); got != 0 {
		t.Fatalf("Auras = %d, want 0 (Arcane Cussing destroyed taking damage)", got)
	}
}

// Tests that a carried Arcane Cussing survives a fully-blocked turn with no attack —
// nothing destroys it, so it neither pops nor pays out.
func TestArcaneCussing_FullyBlockedDoesNotPop(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingDamage(3).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.PitchOneDR{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 0 {
		t.Fatalf("RunechantCount = %d, want 0 (Arcane Cussing did not pop)", got)
	}
	if got := len(summary.State.Auras()); got != 1 {
		t.Fatalf("Auras = %d, want 1 (Arcane Cussing survives a fully-blocked turn)", got)
	}
}

// Tests that, after a fully-blocked turn, a return attack that hits pops Arcane Cussing
// during our turn — paying out its 3 Runechants.
func TestArcaneCussing_PoppedByOwnAttackCreatesRunechants(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingDamage(3).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.PitchOneDR{}, testutils.AttackWithPower{Power: 4}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 3 {
		t.Fatalf("RunechantCount = %d, want 3 (our attack hit and popped Arcane Cussing)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
