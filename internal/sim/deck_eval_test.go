package sim

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/weapon/weapons"
)

// Tests the marginal-stats invariant: PresentHands + AbsentHands == Stats.Hands per card.
func TestEvaluate_PerCardMarginalCoversEveryHand(t *testing.T) {
	read := registry.GetCard(ids.ReadTheRunesRed)
	snatch := registry.GetCard(ids.SnatchRed)
	// 4 of each so Snatch isn't pinned to a single hand and the absent bucket gets exercised.
	deckCards := []deck.Card{read, read, read, read, snatch, snatch, snatch, snatch}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	stats := NewEvaluator().Evaluate(d, 20, Matchup{}, rand.New(rand.NewSource(1)))

	if stats.PerCardMarginal == nil {
		t.Fatalf("PerCardMarginal should be initialised after Evaluate")
	}
	for _, id := range []ids.CardID{ids.ReadTheRunesRed, ids.SnatchRed} {
		m, ok := stats.PerCardMarginal[id]
		if !ok {
			t.Errorf("PerCardMarginal missing entry for %s", registry.GetCard(id).Name())
			continue
		}
		if got := m.PresentHands + m.AbsentHands; got != stats.Hands {
			t.Errorf("%s: PresentHands+AbsentHands = %d, want Stats.Hands = %d (every hand must end up in exactly one bucket)",
				registry.GetCard(id).Name(), got, stats.Hands)
		}
		if m.PresentHands == 0 {
			t.Errorf("%s: PresentHands = 0 — this card should have been in at least one dealt hand across 20 shuffles",
				registry.GetCard(id).Name())
		}
	}
}

// Tests the singleton-deck case: AbsentHands == 0, Marginal() == 0 (no comparison possible).
func TestEvaluate_PerCardMarginalAlwaysPresent(t *testing.T) {
	read := registry.GetCard(ids.ReadTheRunesRed)
	d := deck.New(heroes.Viserai{}, nil, []deck.Card{read, read, read, read, read, read, read, read})
	stats := NewEvaluator().Evaluate(d, 5, Matchup{}, rand.New(rand.NewSource(1)))

	m := stats.PerCardMarginal[ids.ReadTheRunesRed]
	if m.AbsentHands != 0 {
		t.Errorf("AbsentHands = %d, want 0 (single-card deck means card is always present)", m.AbsentHands)
	}
	if m.PresentHands != stats.Hands {
		t.Errorf("PresentHands = %d, want %d (every hand contains the only card in the deck)",
			m.PresentHands, stats.Hands)
	}
	if m.Marginal() != 0 {
		t.Errorf("Marginal() = %v, want 0 (no absent comparison possible)", m.Marginal())
	}
}

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
	if stats.Best.BestLine[0].Role != deck.Arsenal {
		t.Errorf("Best.Play.Roles[0] = %s, want ARSENAL (empty slot on turn 1 → Held promoted)", stats.Best.BestLine[0].Role)
	}
}

// Tests that a card promoted to Arsenal on one turn becomes arsenalCardIn on the next.
func TestEvaluate_ArsenalPersistsAcrossTurns(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 1}, nil, []deck.Card{cards.ToughenUpBlue{}, cards.ToughenUpBlue{}})
	stats := NewEvaluator().Evaluate(d, 1, Matchup{IncomingDamage: 4}, rand.New(rand.NewSource(1)))

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
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
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
	handsPerCycle := len(deckCards) / heroes.Viserai{}.Intelligence()
	maxHands := 2 * handsPerCycle
	if stats.Hands != maxHands {
		t.Errorf("Stats.Hands = %d, want exactly %d (steady-state pitched-pitch loop hits the cap)",
			stats.Hands, maxHands)
	}
}

// Tests that EvaluateAdaptive stops on a multiple of AdaptiveCheckInterval, below the cap,
// when the SE target is met.
func TestEvaluateAdaptive_StopsBeforeMaxRunsWhenSEMet(t *testing.T) {
	deckCards := append([]deck.Card{},
		registry.GetCard(ids.ReadTheRunesRed), registry.GetCard(ids.ReadTheRunesRed),
		registry.GetCard(ids.ReadTheRunesYellow), registry.GetCard(ids.ReadTheRunesYellow),
	)
	for len(deckCards) < 40 {
		deckCards = append(deckCards, registry.GetCard(ids.ReadTheRunesBlue))
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
	stats := NewEvaluator().EvaluateAdaptive(d, 0.1, Matchup{}, rand.New(rand.NewSource(42)))
	if stats.Runs >= AdaptiveShufflesCap {
		t.Errorf("Runs = %d; expected adaptive stop well before cap=%d", stats.Runs, AdaptiveShufflesCap)
	}
	if stats.Runs%AdaptiveCheckInterval != 0 {
		t.Errorf("Runs = %d; expected a multiple of AdaptiveCheckInterval=%d", stats.Runs, AdaptiveCheckInterval)
	}
}

// Tests that EvaluateAdaptive runs to the cap when the SE target is structurally unreachable.
// Drives evaluateImpl directly so the cap can be small.
func TestEvaluateAdaptive_RespectsMaxRunsCapWhenSEUnreachable(t *testing.T) {
	deckCards := append([]deck.Card{},
		registry.GetCard(ids.ReadTheRunesRed), registry.GetCard(ids.ReadTheRunesRed),
		registry.GetCard(ids.ReadTheRunesYellow), registry.GetCard(ids.ReadTheRunesYellow),
	)
	for len(deckCards) < 40 {
		deckCards = append(deckCards, registry.GetCard(ids.ReadTheRunesBlue))
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
	// Negative targetSE is structurally unreachable — MeanStandardError is always >= 0, so
	// the `<= targetSE` predicate never fires. Loop should exhaust at maxRuns regardless of
	// the deck's actual variance (which can be zero for trivially-identical-card decks).
	stats := NewEvaluator().EvaluateImplForTest(d, 1000, Matchup{}, rand.New(rand.NewSource(42)), MakeAdaptiveStop(-1))
	if stats.Runs != 1000 {
		t.Errorf("Runs = %d; expected exactly maxRuns=1000 when SE target is unreachable", stats.Runs)
	}
}

// TestMeanStandardError_FromHistogram verifies the histogram-based SE computation against a
// hand-computed expected value.
func TestMeanStandardError_FromHistogram(t *testing.T) {
	// Three turns: values 10, 12, 14. Mean = 12, sample variance = ((10-12)^2 + (12-12)^2 +
	// (14-12)^2) / (3-1) = 8/2 = 4. SE = sqrt(4/3) ≈ 1.1547.
	stats := &deck.Stats{
		Hands:      3,
		TotalValue: 36,
		Histogram:  map[int]int{10: 1, 12: 1, 14: 1},
	}
	got := MeanStandardError(stats)
	const want = 1.1547005383792515
	if abs := got - want; abs < -1e-9 || abs > 1e-9 {
		t.Errorf("MeanStandardError = %v, want ~%v", got, want)
	}
}
