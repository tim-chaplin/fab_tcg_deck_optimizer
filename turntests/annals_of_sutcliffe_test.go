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

// annalsDeckFiller is the distinctly-named deck card, so a test can assert the card Annals
// drew (then arsenaled at end of turn) came from the deck rather than a leftover hand card.
// A non-attack action (not a resource) so it's arsenal-eligible once drawn.
var annalsDeckFiller = testutils.FakeRedAction().WithName("AnnalsDeckFiller")

// annalsDeck equips Annals of Sutcliffe and stocks a few filler cards so its activation
// has something to draw.
func annalsDeck() *deck.Deck {
	return deck.New(heroes.Viserai, []deck.Weapon{weapons.AnnalsOfSutcliffe{}},
		[]deck.Card{annalsDeckFiller, annalsDeckFiller, annalsDeckFiller})
}

// assertDrewIntoArsenal checks that Annals' draw landed the distinctly-named deck filler in
// the arsenal slot — proof the activation's "Draw a card" clause actually ran.
func assertDrewIntoArsenal(t *testing.T, summary sim.TurnSummary) {
	t.Helper()
	ars := summary.State.Arsenal()
	if ars == nil {
		t.Fatalf("Arsenal empty; expected the drawn %q\nBestLine: %s", annalsDeckFiller.DisplayName(), formatBestLine(summary.BestLine))
	}
	if ars.DisplayName() != annalsDeckFiller.DisplayName() {
		t.Fatalf("Arsenal = %q, want the drawn %q\nBestLine: %s", ars.DisplayName(), annalsDeckFiller.DisplayName(), formatBestLine(summary.BestLine))
	}
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
	assertDrewIntoArsenal(t, summary)
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

	assertDrewIntoArsenal(t, summary)
	if got := summary.State.RunechantCount(); got != 1 {
		t.Errorf("RunechantCount = %d, want 1 (attack + non-attack pitched into Annals)\nValue=%d BestLine: %s",
			got, summary.Value, formatBestLine(summary.BestLine))
	}
}
