package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
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

// Tests that Runeblood Incantation spends its last verse counter for a Runechant one turn,
// then is destroyed the next turn, when the fire finds no counter left to remove.
func TestRunebloodIncantation_LastVerseSpentThenDestroyedNextTurn(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.RunebloodIncantationBlue{}).
		SetIncomingDamage(0).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.BluePitch{}}

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, prior, hand)

	if got := turn1.State.RunechantCount(); got != 1 {
		t.Fatalf("turn 1 RunechantCount = %d, want 1 (last verse counter spent on a Runechant)", got)
	}
	if graveyardContains(turn1.State.Graveyard(), ids.RunebloodIncantationBlue) {
		t.Fatalf("turn 1 graveyard has Runeblood Incantation — it should survive the turn its last counter is removed")
	}
	if got := turn2.State.RunechantCount(); got != 1 {
		t.Fatalf("turn 2 RunechantCount = %d, want 1 (no verse counter left — no new Runechant)", got)
	}
	if !graveyardContains(turn2.State.Graveyard(), ids.RunebloodIncantationBlue) {
		t.Fatalf("turn 2 graveyard = %v, want Runeblood Incantation destroyed (no verse counter to remove)",
			turn2.State.Graveyard())
	}
}
