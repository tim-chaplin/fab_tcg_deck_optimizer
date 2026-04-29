package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"math/rand"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// TestEvaluate_BestTurnStartingRunechantsIsPreHandCarryover pins the contract of
// BestTurn.StartingRunechants: it's the Runechant count carried in from the previous turn when
// the hand was played, so for the first hand of a run it's always 0 — even if the hand itself
// creates runechants that carry out into the next turn.
func TestEvaluate_BestTurnStartingRunechantsIsPreHandCarryover(t *testing.T) {
	// Viserai has Intelligence 4. A 4-card deck gives exactly one hand per run, so the Best
	// record always reflects that first hand — no previous turn ever existed.
	read := GetCard(ids.ReadTheRunesRed)
	d := New(heroes.Viserai{}, nil, []Card{read, read, read, read})

	// Seed doesn't matter (all cards identical), but fix it for determinism.
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	// Sanity: the hand should have left runechants on the table (otherwise the bug couldn't
	// manifest — pre-hand and post-hand counts would both be 0).
	if d.Stats.Best.Summary.Value == 0 {
		t.Fatalf("expected nonzero Value from a hand of Read the Runes; got 0")
	}
	if d.Stats.Best.StartingRunechants != 0 {
		t.Errorf("StartingRunechants = %d, want 0 (first hand of the run has no previous-turn carryover)",
			d.Stats.Best.StartingRunechants)
	}
}

// TestEvaluate_BestTurnSnapshotsState pins the BestTurn snapshot's completeness: the winning
// turn's CarryState (Hand, Deck, Graveyard, Arsenal, Runechants, etc.) must be deep-copied
// into Stats.Best.Summary.State. A Snatch-heavy Viserai hand attacks with at least one Snatch
// (its draw-rider fires DrawOne, pulling another Snatch off the deck), so at least one drawn
// card surfaces in State.Hand or State.Arsenal alongside the played Snatch in State.Graveyard.
// The total card count across the three surfaces must exceed handSize (4) — proof the snapshot
// carried the mid-chain draw rather than just the partition's static slice.
func TestEvaluate_BestTurnSnapshotsState(t *testing.T) {
	snatch := GetCard(ids.SnatchRed)
	d := New(heroes.Viserai{}, nil, []Card{snatch, snatch, snatch, snatch, snatch, snatch, snatch, snatch})
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

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

// TestEvaluate_PerCardMarginalCoversEveryHand pins the marginal-stats invariant: for every
// unique ids.CardID in the deck, PresentHands + AbsentHands equals Stats.Hands. The bucket sums
// are also non-negative and reflect the per-turn hand-value tally so a regression that
// double-counts (or skips) a hand surfaces immediately. A multi-card deck mixes "present
// every turn" vs "present some turns" so both buckets are exercised.
func TestEvaluate_PerCardMarginalCoversEveryHand(t *testing.T) {
	read := GetCard(ids.ReadTheRunesRed)
	snatch := GetCard(ids.SnatchRed)
	// 4 of each so Snatch isn't pinned to a single hand and the absent bucket gets exercised.
	deckCards := []Card{read, read, read, read, snatch, snatch, snatch, snatch}
	d := New(heroes.Viserai{}, nil, deckCards)
	d.Evaluate(20, 0, rand.New(rand.NewSource(1)))

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

// TestEvaluate_PerCardMarginalAlwaysPresent pins the singleton case: a deck containing only
// one unique card never registers an absent bucket. PresentMean equals the deck's overall
// Mean, and Marginal is 0 (no comparison possible). Single-card decks are a degenerate but
// realistic test fixture.
func TestEvaluate_PerCardMarginalAlwaysPresent(t *testing.T) {
	read := GetCard(ids.ReadTheRunesRed)
	d := New(heroes.Viserai{}, nil, []Card{read, read, read, read, read, read, read, read})
	d.Evaluate(5, 0, rand.New(rand.NewSource(1)))

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

// TestEvaluate_HeldCardDefersDrawToNextTurn pins the "up to Intelligence" draw rule plus arsenal
// carryover. Intelligence-1 hero, deck of Toughen Up Blue (DR, cost 2, defense 4): the lone
// card has no legal play (can't pay its 2-cost, can't pitch with nothing on the stack, DRs
// can't Attack). Turn 1 holds then promotes it to Arsenal (empty slot). Turn 2 draws a new DR;
// the arsenal card stays on tie, the new card goes Held, so drawCount = 0 next turn and the
// loop halts at Stats.Hands = 2.
func TestEvaluate_HeldCardDefersDrawToNextTurn(t *testing.T) {
	// 40 copies of the DR so we have enough deck to fill many hands if held carryover weren't
	// wired up — the assertion would fail catastrophically (loop or much larger Hands count).
	deckCards := make([]Card, 40)
	for i := range deckCards {
		deckCards[i] = cards.ToughenUpBlue{}
	}
	d := New(Int1StubHero, nil, deckCards)
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

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

// TestEvaluate_ArsenalPersistsAcrossTurns confirms the arsenal slot threads through Evaluate's
// per-turn loop: a card promoted to Arsenal on one turn becomes arsenalCardIn on the next.
// Intelligence-1 hero, 2-card deck of two Toughen Up Blue. Turn 1 arsenals the drawn TU.
// Turn 2 draws the second TU and against incoming 4 plays the arsenal-in DR, pitching the
// drawn card to fund its 2-cost — Value = 4 (prevents the full attack). Turn 3 re-draws the
// pitched card (returned to deck bottom) and arsenals it again. Loop stops when the deck's
// empty and nothing new can be drawn.
func TestEvaluate_ArsenalPersistsAcrossTurns(t *testing.T) {
	d := New(Int1StubHero, nil, []Card{cards.ToughenUpBlue{}, cards.ToughenUpBlue{}})
	d.Evaluate(1, 4, rand.New(rand.NewSource(1)))

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

// TestEvaluate_TerminatesAfterTwoCycles pins the infinite-loop guard on Evaluate's per-run loop.
// 40 Toughen Up Blue DRs with Reaping Blade equipped, incoming=0, reaches a steady state after
// turn 1 (pitch one TU, swing Reaping Blade for +3, hold the other 3). From then on every turn
// draws and pitches one card — net deck change zero, Best returns the same TurnSummary.
// Without the cap the loop would spin forever; with it, Stats.Hands halts at 2 × handsPerCycle.
func TestEvaluate_TerminatesAfterTwoCycles(t *testing.T) {
	deckCards := make([]Card, 40)
	for i := range deckCards {
		deckCards[i] = cards.ToughenUpBlue{}
	}
	d := New(heroes.Viserai{}, []Weapon{weapons.ReapingBlade{}}, deckCards)
	done := make(chan struct{})
	go func() {
		d.Evaluate(1, 0, rand.New(rand.NewSource(1)))
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

// TestEvaluateAdaptive_StopsBeforeMaxRunsWhenSEMet runs a deck whose per-turn variance is
// low enough that the SE target is met well before the cap. The actual run count should
// land on a multiple of AdaptiveCheckInterval and be strictly less than AdaptiveShufflesCap.
func TestEvaluateAdaptive_StopsBeforeMaxRunsWhenSEMet(t *testing.T) {
	deckCards := append([]Card{},
		GetCard(ids.ReadTheRunesRed), GetCard(ids.ReadTheRunesRed),
		GetCard(ids.ReadTheRunesYellow), GetCard(ids.ReadTheRunesYellow),
	)
	for len(deckCards) < 40 {
		deckCards = append(deckCards, GetCard(ids.ReadTheRunesBlue))
	}
	d := New(heroes.Viserai{}, []Weapon{weapons.ReapingBlade{}}, deckCards)
	stats := d.EvaluateAdaptive(0, rand.New(rand.NewSource(42)))
	if stats.Runs >= AdaptiveShufflesCap {
		t.Errorf("Runs = %d; expected adaptive stop well before cap=%d", stats.Runs, AdaptiveShufflesCap)
	}
	if stats.Runs%AdaptiveCheckInterval != 0 {
		t.Errorf("Runs = %d; expected a multiple of AdaptiveCheckInterval=%d", stats.Runs, AdaptiveCheckInterval)
	}
}

// TestEvaluateAdaptive_RespectsMaxRunsCapWhenSEUnreachable pins the cap behavior: when the
// SE target is structurally unreachable, the loop runs to the cap rather than overshooting.
// Goes through evaluateImpl directly so the test can use a small cap (production uses
// AdaptiveShufflesCap=50000 which is too slow for a unit test).
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
	stats := d.EvaluateImplForTest(1000, 0, rand.New(rand.NewSource(42)), nil, MakeAdaptiveStop(-1))
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
