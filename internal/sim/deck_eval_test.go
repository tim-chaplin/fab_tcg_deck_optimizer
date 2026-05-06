package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"math/rand"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests that BestTurn.StartingAuras reflects pre-hand carryover (empty on the first hand).
func TestEvaluate_BestTurnStartingAurasIsPreHandCarryover(t *testing.T) {
	// 4-card deck + Viserai Intel=4 gives exactly one hand per run.
	read := GetCard(ids.ReadTheRunesRed)
	d := New(heroes.Viserai{}, nil, []Card{read, read, read, read})

	d.Evaluate(1, Matchup{}, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	// Sanity: the hand creates runechants — without that, the assertion below couldn't
	// distinguish "no carryover" from "no token tracking happening at all".
	if d.Stats.Best.Summary.Value == 0 {
		t.Fatalf("expected nonzero Value from a hand of Read the Runes; got 0")
	}
	if len(d.Stats.Best.StartingAuras) != 0 {
		t.Errorf("StartingAuras = %v, want empty (first hand of the run has no previous-turn carryover)",
			d.Stats.Best.StartingAuras)
	}
}

// Tests that BestTurn deep-copies the winning turn's full CarryState (Hand, Deck, Graveyard,
// Arsenal, Log, etc.) — including mid-chain draws.
func TestEvaluate_BestTurnSnapshotsState(t *testing.T) {
	snatch := GetCard(ids.SnatchRed)
	d := New(heroes.Viserai{}, nil, []Card{snatch, snatch, snatch, snatch, snatch, snatch, snatch, snatch})
	d.Evaluate(1, Matchup{}, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	state := d.Stats.Best.Summary.State
	if len(state.Graveyard) == 0 {
		t.Errorf("State.Graveyard is empty; want the played Snatch in graveyard")
	}
	surfaceCount := len(state.Hand) + len(state.Graveyard)
	if state.Arsenal != nil {
		surfaceCount++
	}
	const handSize = 4 // Viserai's Intelligence
	if surfaceCount <= handSize {
		t.Errorf("surface count = %d, want >%d (Hand=%d Arsenal=%v Graveyard=%d). The mid-turn-drawn Snatch should have surfaced — without the State snapshot the carry would lose it.",
			surfaceCount, handSize, len(state.Hand), state.Arsenal, len(state.Graveyard))
	}
	// State.Log carries the per-event chain trace; recordBestTurn → CarryState.Clone must copy
	// it through so fabsim eval's "Best turn played" printout has the chain attribution lines.
	// A Snatch chain has at least one ATTACK entry — if Log is empty the deck-level snapshot
	// dropped the Log field on the floor.
	if len(state.Log) == 0 {
		t.Errorf("State.Log is empty after the deck-level snapshot; FormatBestTurn will render no chain lines for the saved Best")
	}
}

// Tests the marginal-stats invariant: PresentHands + AbsentHands == Stats.Hands per card.
func TestEvaluate_PerCardMarginalCoversEveryHand(t *testing.T) {
	read := GetCard(ids.ReadTheRunesRed)
	snatch := GetCard(ids.SnatchRed)
	// 4 of each so Snatch isn't pinned to a single hand and the absent bucket gets exercised.
	deckCards := []Card{read, read, read, read, snatch, snatch, snatch, snatch}
	d := New(heroes.Viserai{}, nil, deckCards)
	d.Evaluate(20, Matchup{}, rand.New(rand.NewSource(1)))

	if d.Stats.PerCardMarginal == nil {
		t.Fatalf("PerCardMarginal should be initialised after Evaluate")
	}
	for _, id := range []ids.CardID{ids.ReadTheRunesRed, ids.SnatchRed} {
		m, ok := d.Stats.PerCardMarginal[id]
		if !ok {
			t.Errorf("PerCardMarginal missing entry for %s", GetCard(id).Name())
			continue
		}
		if got := m.PresentHands + m.AbsentHands; got != d.Stats.Hands {
			t.Errorf("%s: PresentHands+AbsentHands = %d, want Stats.Hands = %d (every hand must end up in exactly one bucket)",
				GetCard(id).Name(), got, d.Stats.Hands)
		}
		if m.PresentHands == 0 {
			t.Errorf("%s: PresentHands = 0 — this card should have been in at least one dealt hand across 20 shuffles",
				GetCard(id).Name())
		}
	}
}

// Tests the singleton-deck case: AbsentHands == 0, Marginal() == 0 (no comparison possible).
func TestEvaluate_PerCardMarginalAlwaysPresent(t *testing.T) {
	read := GetCard(ids.ReadTheRunesRed)
	d := New(heroes.Viserai{}, nil, []Card{read, read, read, read, read, read, read, read})
	d.Evaluate(5, Matchup{}, rand.New(rand.NewSource(1)))

	m := d.Stats.PerCardMarginal[ids.ReadTheRunesRed]
	if m.AbsentHands != 0 {
		t.Errorf("AbsentHands = %d, want 0 (single-card deck means card is always present)", m.AbsentHands)
	}
	if m.PresentHands != d.Stats.Hands {
		t.Errorf("PresentHands = %d, want %d (every hand contains the only card in the deck)",
			m.PresentHands, d.Stats.Hands)
	}
	if m.Marginal() != 0 {
		t.Errorf("Marginal() = %v, want 0 (no absent comparison possible)", m.Marginal())
	}
}

// Tests that a Held DR carries into the next turn and reduces the next draw count, halting
// the loop when no further plays are possible.
func TestEvaluate_HeldCardDefersDrawToNextTurn(t *testing.T) {
	// 40 copies so a missing held-carryover regression would loop or run far longer than 2 hands.
	deckCards := make([]Card, 40)
	for i := range deckCards {
		deckCards[i] = cards.ToughenUpBlue{}
	}
	d := New(testutils.Hero{Intel: 1}, nil, deckCards)
	d.Evaluate(1, Matchup{}, rand.New(rand.NewSource(1)))

	if d.Stats.Hands != 2 {
		t.Errorf("Stats.Hands = %d, want 2 (turn 1 arsenals the card, turn 2 holds its successor, turn 3 can't draw)", d.Stats.Hands)
	}
	// Best captures turn 1 (first hand with a recorded play). That hand's single card got
	// promoted from Held to Arsenal by the post-hoc upgrade.
	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after at least one hand")
	}
	if d.Stats.Best.Summary.BestLine[0].Role != Arsenal {
		t.Errorf("Best.Play.Roles[0] = %s, want ARSENAL (empty slot on turn 1 → Held promoted)", d.Stats.Best.Summary.BestLine[0].Role)
	}
}

// Tests that a card promoted to Arsenal on one turn becomes arsenalCardIn on the next.
func TestEvaluate_ArsenalPersistsAcrossTurns(t *testing.T) {
	d := New(testutils.Hero{Intel: 1}, nil, []Card{cards.ToughenUpBlue{}, cards.ToughenUpBlue{}})
	d.Evaluate(1, Matchup{IncomingDamage: 4}, rand.New(rand.NewSource(1)))

	// Best captures turn 2 — only turn with Value > 0 (arsenal DR fires).
	if d.Stats.Best.Summary.Value != 4 {
		t.Errorf("Best.Play.Value = %d, want 4 (turn 2 plays arsenal DR, pitches hand DR to pay; prevents 4)", d.Stats.Best.Summary.Value)
	}
	// Turn 1: arsenal the drawn card. Turn 2: play arsenal DR (paid by pitching drawn card).
	// Turn 3: draw the recycled pitched card, arsenal it (deck is then empty). Loop ends.
	if d.Stats.Hands != 3 {
		t.Errorf("Stats.Hands = %d, want 3", d.Stats.Hands)
	}
}

// Tests Evaluate's infinite-loop guard: a steady-state pitched-pitch cycle halts at
// 2 × handsPerCycle.
func TestEvaluate_TerminatesAfterTwoCycles(t *testing.T) {
	deckCards := make([]Card, 40)
	for i := range deckCards {
		deckCards[i] = cards.ToughenUpBlue{}
	}
	d := New(heroes.Viserai{}, []Weapon{weapons.ReapingBlade{}}, deckCards)
	done := make(chan struct{})
	go func() {
		d.Evaluate(1, Matchup{}, rand.New(rand.NewSource(1)))
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
	if d.Stats.Hands != maxHands {
		t.Errorf("Stats.Hands = %d, want exactly %d (steady-state pitched-pitch loop hits the cap)",
			d.Stats.Hands, maxHands)
	}
}

// Tests that EvaluateAdaptive stops on a multiple of AdaptiveCheckInterval, below the cap,
// when the SE target is met.
func TestEvaluateAdaptive_StopsBeforeMaxRunsWhenSEMet(t *testing.T) {
	deckCards := append([]Card{},
		GetCard(ids.ReadTheRunesRed), GetCard(ids.ReadTheRunesRed),
		GetCard(ids.ReadTheRunesYellow), GetCard(ids.ReadTheRunesYellow),
	)
	for len(deckCards) < 40 {
		deckCards = append(deckCards, GetCard(ids.ReadTheRunesBlue))
	}
	d := New(heroes.Viserai{}, []Weapon{weapons.ReapingBlade{}}, deckCards)
	stats := d.EvaluateAdaptive(Matchup{}, rand.New(rand.NewSource(42)))
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
	deckCards := append([]Card{},
		GetCard(ids.ReadTheRunesRed), GetCard(ids.ReadTheRunesRed),
		GetCard(ids.ReadTheRunesYellow), GetCard(ids.ReadTheRunesYellow),
	)
	for len(deckCards) < 40 {
		deckCards = append(deckCards, GetCard(ids.ReadTheRunesBlue))
	}
	d := New(heroes.Viserai{}, []Weapon{weapons.ReapingBlade{}}, deckCards)
	// Negative targetSE is structurally unreachable — MeanStandardError is always >= 0, so
	// the `<= targetSE` predicate never fires. Loop should exhaust at maxRuns regardless of
	// the deck's actual variance (which can be zero for trivially-identical-card decks).
	stats := d.EvaluateImplForTest(1000, Matchup{}, rand.New(rand.NewSource(42)), nil, MakeAdaptiveStop(-1))
	if stats.Runs != 1000 {
		t.Errorf("Runs = %d; expected exactly maxRuns=1000 when SE target is unreachable", stats.Runs)
	}
}

// TestMeanStandardError_FromHistogram verifies the histogram-based SE computation against a
// hand-computed expected value.
func TestMeanStandardError_FromHistogram(t *testing.T) {
	// Three turns: values 10, 12, 14. Mean = 12, sample variance = ((10-12)^2 + (12-12)^2 +
	// (14-12)^2) / (3-1) = 8/2 = 4. SE = sqrt(4/3) ≈ 1.1547.
	stats := &Stats{
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
