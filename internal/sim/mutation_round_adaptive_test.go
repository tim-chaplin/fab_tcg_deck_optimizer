package sim

import (
	"context"
	"math/rand"
	"sync/atomic"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// adaptiveFixture builds a deliberately weak random deck (so plenty of mutations improve it) plus
// its mutation list, the shared inputs for the adaptive-round tests.
func adaptiveFixture(t *testing.T, seed int64) (*deck.Deck, []deck.Mutation) {
	t.Helper()
	rng := rand.New(rand.NewSource(seed))
	d := deck.Random(heroes.Viserai, 40, 2, rng, registry.Registry{})
	muts := deck.AllMutations(d, 2, false, registry.Registry{})
	if len(muts) > 60 {
		muts = muts[:60] // keep the test quick; 60 candidates is plenty to find an improvement
	}
	return d, muts
}

// seqValue evaluates a deck sequentially at a fixed seed — the same path the adaptive confirm uses,
// so the test can reconstruct the confirm comparison exactly.
func seqValue(d *deck.Deck, shuffles int, mp Matchup, seed int64) float64 {
	return NewEvaluator().Evaluate(d, shuffles, mp, rand.New(rand.NewSource(seed))).Mean()
}

// TestRunMutationRoundAdaptive_ConfirmedImprovementClearsThreshold verifies the core contract: when
// the round reports an improvement, the accepted deck genuinely beats the incumbent by more than the
// threshold under the high-precision coupled confirm — i.e. a screen false-accept can't win.
func TestRunMutationRoundAdaptive_ConfirmedImprovementClearsThreshold(t *testing.T) {
	const (
		threshold     = 0.1
		statsShuffles = 400
		seed          = 4242
	)
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 1)

	// One mutation worker → deterministic FIFO winner, so the reconstruction below is exact.
	d, stats, avg, idx, found := RunMutationRoundAdaptive(
		context.Background(), muts, incumbent, threshold, mp, statsShuffles, 1, seed, nil, nil)
	if !found {
		t.Skip("weak deck produced no improvement at this seed; nothing to confirm")
	}
	if d == nil || idx < 0 || idx >= len(muts) {
		t.Fatalf("found=true but d=%v idx=%d", d, idx)
	}
	if stats.Mean() != avg {
		t.Errorf("returned avg %.4f != stats.Mean() %.4f", avg, stats.Mean())
	}
	// Reconstruct the confirm: incumbent and mutant at the confirm seed (perShuffleSeed(seed,0)).
	confirmSeed := perShuffleSeed(seed, 0)
	incAvg := seqValue(incumbent, statsShuffles, mp, confirmSeed)
	mutAvg := seqValue(d, statsShuffles, mp, confirmSeed)
	if mutAvg-incAvg <= threshold {
		t.Errorf("accepted mutation does not clear the confirm gate: mutAvg %.4f - incAvg %.4f = %.4f <= %.2f",
			mutAvg, incAvg, mutAvg-incAvg, threshold)
	}
	if mutAvg != avg {
		t.Errorf("reconstructed confirm avg %.4f != returned avg %.4f (coupling/seed mismatch)", mutAvg, avg)
	}
}

// TestRunMutationRoundAdaptive_UnreachableThresholdRejectsAll: with an impossibly high threshold no
// candidate clears the screen, so the round drains and reports no improvement.
func TestRunMutationRoundAdaptive_UnreachableThresholdRejectsAll(t *testing.T) {
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 2)

	var completed atomic.Int64
	d, _, avg, idx, found := RunMutationRoundAdaptive(
		context.Background(), muts, incumbent, 1_000_000.0, mp, 200, 0, 7, &completed, nil)
	if found {
		t.Errorf("found=true with an unreachable threshold; want false")
	}
	if d != nil || idx != -1 || avg != 0 {
		t.Errorf("no-improvement return shape: d=%v avg=%v idx=%d, want nil/0/-1", d, avg, idx)
	}
	if int(completed.Load()) != len(muts) {
		t.Errorf("completed=%d, want every candidate screened (%d)", completed.Load(), len(muts))
	}
}

// TestRunMutationRoundAdaptive_Deterministic: a single mutation worker with a fixed seed yields the
// same outcome across runs.
func TestRunMutationRoundAdaptive_Deterministic(t *testing.T) {
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 3)
	run := func() (int, bool, float64) {
		d, _, avg, idx, found := RunMutationRoundAdaptive(
			context.Background(), muts, incumbent, 0.1, mp, 300, 1, 99, nil, nil)
		_ = d
		return idx, found, avg
	}
	idx1, found1, avg1 := run()
	idx2, found2, avg2 := run()
	if idx1 != idx2 || found1 != found2 || avg1 != avg2 {
		t.Errorf("non-deterministic: (idx=%d found=%v avg=%.4f) vs (idx=%d found=%v avg=%.4f)",
			idx1, found1, avg1, idx2, found2, avg2)
	}
}
