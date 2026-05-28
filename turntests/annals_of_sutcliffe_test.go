package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// annalsDeck equips Annals of Sutcliffe and stocks a few filler cards so its activation
// has something to draw.
func annalsDeck() *deck.Deck {
	filler := testutils.FakeBlueResource()
	return deck.New(heroes.Viserai, []deck.Weapon{weapons.AnnalsOfSutcliffe{}},
		[]deck.Card{filler, filler, filler})
}

// Tests that activating Annals while pitching a single resource (not an attack action plus
// a non-attack action) draws a card but makes no Runechant. The blue must be pitched to
// fund Annals' 3-cost activation — a pitch that funds nothing is disqualified, so a Pitch
// role here proves Annals actually fired and still produced no Runechant.
func TestAnnalsOfSutcliffe_SingleResourcePitchNoRunechant(t *testing.T) {
	// One blue resource pitches for 3 — exactly Annals' activation cost.
	hand := []card.Card{testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(annalsDeck(), gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	pitched := false
	for _, a := range summary.BestLine {
		if a.Role == card.Pitch {
			pitched = true
		}
	}
	if !pitched {
		t.Fatalf("expected the resource pitched to fund Annals' activation\nBestLine: %s", formatBestLine(summary.BestLine))
	}
	if got := summary.State.RunechantCount(); got != 0 {
		t.Errorf("RunechantCount = %d, want 0 (a lone resource pitch isn't attack + non-attack)\nValue=%d BestLine: %s",
			got, summary.Value, formatBestLine(summary.BestLine))
	}
}

// Tests that activating Annals while pitching an attack action card AND a non-attack action
// card draws a card and creates a Runechant.
func TestAnnalsOfSutcliffe_AttackPlusNonAttackPitchMakesRunechant(t *testing.T) {
	// Red 0-power attack action pitches for 1, yellow 0-power non-attack action pitches for
	// 2 — together exactly Annals' 3-resource activation, and the attack + non-attack pair
	// the Runechant rider wants.
	hand := []card.Card{
		testutils.FakeRedAttack(),
		testutils.FakeYellowAction(),
	}
	summary := sim.EvalOneTurnForTesting(annalsDeck(), gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	if got := summary.State.RunechantCount(); got != 1 {
		t.Errorf("RunechantCount = %d, want 1 (attack + non-attack pitched into Annals)\nValue=%d BestLine: %s",
			got, summary.Value, formatBestLine(summary.BestLine))
	}
}
