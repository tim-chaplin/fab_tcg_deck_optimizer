package sim

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that a Held DR carries into the next turn and reduces the next draw count, halting
// the loop when no further plays are possible.
func TestEvaluate_HeldCardDefersDrawToNextTurn(t *testing.T) {
	// 40 copies so a missing held-carryover regression would loop or run far longer than 2 hands.
	deckCards := make([]deck.Card, 40)
	for i := range deckCards {
		deckCards[i] = cards.ToughenUpBlue{}
	}
	d := deck.New(testutils.Hero{Intel: 1}, nil, deckCards)
	stats := NewEvaluator().Evaluate(d, 1, Matchup{}, rand.New(rand.NewSource(1)))

	if stats.Hands != 2 {
		t.Errorf("Stats.Hands = %d, want 2 (turn 1 arsenals the card, turn 2 holds its successor, turn 3 can't draw)", stats.Hands)
	}
	// Best captures turn 1 (first hand with a recorded play). That hand's single card got
	// promoted from Held to Arsenal by the post-hoc upgrade.
	if len(stats.Best.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after at least one hand")
	}
	if stats.Best.BestLine[0].Role != card.Arsenal {
		t.Errorf("Best.Play.Roles[0] = %s, want ARSENAL (empty slot on turn 1 → Held promoted)", stats.Best.BestLine[0].Role)
	}
}

// Tests that a card promoted to Arsenal on one turn becomes arsenalCardIn on the next.
func TestEvaluate_ArsenalPersistsAcrossTurns(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 1}, nil, []deck.Card{cards.ToughenUpBlue{}, cards.ToughenUpBlue{}})
	stats := NewEvaluator().Evaluate(d, 1, Matchup{IncomingPhysicalDamage: 4}, rand.New(rand.NewSource(1)))

	// Best captures turn 2 — only turn with Value > 0 (arsenal DR fires).
	if stats.Best.Value != 4 {
		t.Errorf("Best.Play.Value = %d, want 4 (turn 2 plays arsenal DR, pitches hand DR to pay; prevents 4)", stats.Best.Value)
	}
	// Turn 1: arsenal the drawn card. Turn 2: play arsenal DR (paid by pitching drawn card).
	// Turn 3: draw the recycled pitched card, arsenal it (deck is then empty). Loop ends.
	if stats.Hands != 3 {
		t.Errorf("Stats.Hands = %d, want 3", stats.Hands)
	}
}

// Tests Evaluate's infinite-loop guard: a steady-state pitched-pitch cycle halts at
// 2 × handsPerCycle.
func TestEvaluate_TerminatesAfterTwoCycles(t *testing.T) {
	deckCards := make([]deck.Card, 40)
	for i := range deckCards {
		deckCards[i] = cards.ToughenUpBlue{}
	}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
	done := make(chan struct{})
	var stats deck.Stats
	go func() {
		stats = NewEvaluator().Evaluate(d, 1, Matchup{}, rand.New(rand.NewSource(1)))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("Evaluate did not terminate within 2 seconds — infinite loop regression")
	}
	// Two cycles of a 40-card / 4-hand-size deck is exactly 20 hands.
	handsPerCycle := len(deckCards) / heroes.Viserai.Intelligence()
	maxHands := 2 * handsPerCycle
	if stats.Hands != maxHands {
		t.Errorf("Stats.Hands = %d, want exactly %d (steady-state pitched-pitch loop hits the cap)",
			stats.Hands, maxHands)
	}
}
