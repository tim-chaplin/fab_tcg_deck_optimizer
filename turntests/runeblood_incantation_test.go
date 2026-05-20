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

// Tests that a carried Runeblood Incantation spends a verse counter at the start of the
// turn to create a Runechant, surviving with verses to spare.
func TestRunebloodIncantation_StartOfTurnFireCreatesRunechant(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.RunebloodIncantationRed{}).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 1 {
		t.Fatalf("RunechantCount = %d, want 1 (one verse counter spent on a Runechant)", got)
	}
	if got := len(summary.State.Auras()); got != 2 {
		t.Fatalf("Auras = %d, want 2 (Runeblood Incantation survives alongside its Runechant)", got)
	}
}

// Tests that Runeblood Incantation on its last verse counter creates a final Runechant and
// is then destroyed.
func TestRunebloodIncantation_LastVerseCreatesRunechantAndDestroys(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.RunebloodIncantationBlue{}).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 1 {
		t.Fatalf("RunechantCount = %d, want 1 (last verse counter spent on a Runechant)", got)
	}
	if got := len(summary.State.Auras()); got != 1 {
		t.Fatalf("Auras = %d, want 1 (Runeblood Incantation destroyed, only its Runechant remains)", got)
	}
}
