package sim

import (
	"context"
	"math"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// fakeRecorder captures Promote calls for the recordImprovement test.
type fakeRecorder struct{ promoted []ids.CardID }

func (f *fakeRecorder) Promote(id ids.CardID) {
	f.promoted = append(f.promoted, id)
}

// TestRecordImprovement pins the ranking signal: a single-swap mutation promotes the card it added,
// and a non-single-swap (weapon / pair) mutation promotes nothing.
func TestRecordImprovement(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(heroes.Viserai, 40, 2, rng, registry.Registry{})
	var swap, weapon deck.Mutation
	for _, m := range deck.AllMutations(d, 2, false, registry.Registry{}) {
		if _, _, ok := m.Swap(); ok {
			swap = m
		} else {
			weapon = m
		}
	}
	added, _, _ := swap.Swap()

	cases := []struct {
		name string
		mut  deck.Mutation
		want []ids.CardID
	}{
		{"single swap promotes the added card", swap, []ids.CardID{added}},
		{"non-single-swap promotes nothing", weapon, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := &fakeRecorder{}
			recordImprovement(f, c.mut)
			if len(f.promoted) != len(c.want) {
				t.Fatalf("promoted %v, want %v", f.promoted, c.want)
			}
			for i := range c.want {
				if f.promoted[i] != c.want[i] {
					t.Errorf("promoted[%d] = %v, want %v", i, f.promoted[i], c.want[i])
				}
			}
		})
	}

	recordImprovement(nil, swap) // nil recorder must be a no-op (no panic)
}

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

// TestRunMutationRoundAdaptive_ConfirmedImprovementBeatsIncumbent verifies the core contract: when
// the round reports an improvement, the accepted deck genuinely beats the incumbent under the
// high-precision coupled confirm — i.e. a screen false-accept can't win.
func TestRunMutationRoundAdaptive_ConfirmedImprovementBeatsIncumbent(t *testing.T) {
	const (
		statsShuffles = 400
		seed          = 4242
	)
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 1)

	// One mutation worker → deterministic FIFO winner, so the reconstruction below is exact.
	d, stats, avg, idx, found := RunMutationRoundAdaptive(context.Background(), AdaptiveRoundConfig{
		Mutations: muts, Incumbent: incumbent, Matchup: mp,
		StatsShuffles: statsShuffles, MutationWorkers: 1, Seed: seed,
	})
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
	if mutAvg-incAvg <= 0 {
		t.Errorf("accepted mutation does not beat the incumbent under the confirm: mutAvg %.4f - incAvg %.4f = %.4f <= 0",
			mutAvg, incAvg, mutAvg-incAvg)
	}
	if mutAvg != avg {
		t.Errorf("reconstructed confirm avg %.4f != returned avg %.4f (coupling/seed mismatch)", mutAvg, avg)
	}
}

// TestRunMutationRoundAdaptive_MultiWorkerWinnerBeatsIncumbent runs with several mutation workers
// (so the winner is racy) and asserts that whichever candidate wins still beats the incumbent under
// the confirm — the safety property the confirm exists to guarantee, independent of which worker
// lands first.
func TestRunMutationRoundAdaptive_MultiWorkerWinnerBeatsIncumbent(t *testing.T) {
	const (
		statsShuffles = 400
		seed          = 8675309
		workers       = 4
	)
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 1)

	d, _, avg, _, found := RunMutationRoundAdaptive(context.Background(), AdaptiveRoundConfig{
		Mutations: muts, Incumbent: incumbent, Matchup: mp,
		StatsShuffles: statsShuffles, MutationWorkers: workers, Seed: seed,
	})
	if !found {
		t.Skip("weak deck produced no improvement at this seed")
	}
	confirmSeed := perShuffleSeed(seed, 0)
	incAvg := seqValue(incumbent, statsShuffles, mp, confirmSeed)
	mutAvg := seqValue(d, statsShuffles, mp, confirmSeed)
	if mutAvg-incAvg <= 0 {
		t.Errorf("multi-worker winner does not beat the incumbent under the confirm: %.4f - %.4f = %.4f <= 0",
			mutAvg, incAvg, mutAvg-incAvg)
	}
	if mutAvg != avg {
		t.Errorf("reconstructed confirm avg %.4f != returned avg %.4f", mutAvg, avg)
	}
}

// TestRunMutationRoundAdaptive_SAStepClearsItsTau verifies the T>0 path: the winner clears its own
// Metropolis threshold tau = temperature·ln(U) under the coupled confirm (the accepted deck may be
// worse than the incumbent — an SA step is allowed to move downhill).
func TestRunMutationRoundAdaptive_SAStepClearsItsTau(t *testing.T) {
	const (
		temperature   = 0.5 // warm enough that downhill moves can be accepted
		statsShuffles = 400
		seed          = 13
	)
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 1)

	d, _, avg, idx, found := RunMutationRoundAdaptive(context.Background(), AdaptiveRoundConfig{
		Mutations: muts, Incumbent: incumbent, Temperature: temperature,
		Matchup: mp, StatsShuffles: statsShuffles, MutationWorkers: 1, Seed: seed,
	})
	if !found {
		t.Skip("no candidate cleared its SA threshold at this seed")
	}
	// Reconstruct the winner's tau the same way the worker draws it, then check the confirm gate.
	u := rand.New(rand.NewSource(perShuffleSeed(seed, -(idx + 1)))).Float64()
	tau := temperature * math.Log(1-u)
	confirmSeed := perShuffleSeed(seed, 0)
	mutAvg := seqValue(d, statsShuffles, mp, confirmSeed)
	if mutAvg-seqValue(incumbent, statsShuffles, mp, confirmSeed) <= tau {
		t.Errorf("SA winner does not clear its tau %.4f (confirm delta %.4f)",
			tau, mutAvg-seqValue(incumbent, statsShuffles, mp, confirmSeed))
	}
	if mutAvg != avg {
		t.Errorf("reconstructed confirm avg %.4f != returned avg %.4f", mutAvg, avg)
	}
}

// TestRunMutationRoundAdaptive_UnreachableThresholdRejectsAll: with an impossibly wide accept
// threshold (a huge ConfirmSigmas inflates k·σ̂/√N out of reach) no candidate clears the screen, so
// the round drains and reports no improvement.
func TestRunMutationRoundAdaptive_UnreachableThresholdRejectsAll(t *testing.T) {
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 2)

	var progress AdaptiveProgress
	d, _, avg, idx, found := RunMutationRoundAdaptive(context.Background(), AdaptiveRoundConfig{
		Mutations: muts, Incumbent: incumbent, ConfirmSigmas: 1e9, Matchup: mp,
		StatsShuffles: 200, Seed: 7, Progress: &progress,
	})
	if found {
		t.Errorf("found=true with an unreachable threshold; want false")
	}
	if d != nil || idx != -1 || avg != 0 {
		t.Errorf("no-improvement return shape: d=%v avg=%v idx=%d, want nil/0/-1", d, avg, idx)
	}
	if got := progress.Screened.Load(); int(got) != len(muts) {
		t.Errorf("screened=%d, want every candidate (%d)", got, len(muts))
	}
}

// TestRunMutationRoundAdaptive_Deterministic: a single mutation worker with a fixed seed yields the
// same outcome across runs.
func TestRunMutationRoundAdaptive_Deterministic(t *testing.T) {
	mp := Matchup{IncomingPhysicalDamage: 5}
	incumbent, muts := adaptiveFixture(t, 3)
	run := func() (int, bool, float64) {
		d, _, avg, idx, found := RunMutationRoundAdaptive(context.Background(), AdaptiveRoundConfig{
			Mutations: muts, Incumbent: incumbent, Matchup: mp,
			StatsShuffles: 300, MutationWorkers: 1, Seed: 99,
		})
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
