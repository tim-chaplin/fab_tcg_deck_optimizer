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

// Tests that Reinforce the Line prevents N damage while a defending attack action card is
// present: a 3-defense attack action blocks 3 and Reinforce the Line prevents the other 4.
func TestReinforceTheLine_PreventsDamageWithAttackActionDefender(t *testing.T) {
	prior := gameengine.GameStateBuilder().SetIncomingPhysicalDamage(7).Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.ReinforceTheLineRed{}, testutils.FakeBlueAttack().WithDefense(3)}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 7 {
		t.Fatalf("Value = %d, want 7 (BlueAttack blocks 3, Reinforce the Line prevents 4)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that Reinforce the Line prevents nothing when no defending attack action card is
// present — only the defense reaction's own block scores.
func TestReinforceTheLine_NoAttackActionDefenderPreventsNothing(t *testing.T) {
	prior := gameengine.GameStateBuilder().SetIncomingPhysicalDamage(3).Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.ReinforceTheLineRed{}, testutils.FakeRedDR().WithDefense(3)}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.Value; got != 3 {
		t.Fatalf("Value = %d, want 3 (PitchOneDR blocks 3; Reinforce the Line prevents nothing)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
