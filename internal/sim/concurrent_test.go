package sim

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Tests that concurrent Evaluate from many goroutines (each with its own Evaluator)
// doesn't panic on the shared card-meta lookup table.
func TestEvaluate_ConcurrentNoMapPanic(t *testing.T) {
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 2 {
		t.Skip("need GOMAXPROCS >= 2 to exercise concurrent access")
	}
	const iterations = 25

	baseline := deck.Random(heroes.Viserai, 40, 2, rand.New(rand.NewSource(42)), registry.Registry{})

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ev := NewEvaluator()
			rng := rand.New(rand.NewSource(int64(id)*7919 + 1))
			for i := 0; i < iterations; i++ {
				d := baseline.Copy()
				// Small shuffle count per iteration keeps the test fast while still exercising
				// many Best calls per goroutine: 10 shuffles × handsPerCycle * 2 hands each =
				// ~200 Best invocations per goroutine, per iteration.
				stats := ev.Evaluate(d, 10, Matchup{}, rng)
				if stats.Hands == 0 {
					t.Errorf("worker %d iter %d: Evaluate returned zero hands", id, i)
					return
				}
			}
		}(w)
	}
	wg.Wait()
}

// Tests RunMutationRound as a smoke test: completes without deadlock and reports invariant
// shapes regardless of whether an improvement is found.
func TestRunMutationRound_RunsWithoutPanic(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := deck.Random(heroes.Viserai, 40, 2, rng, registry.Registry{})
	baseAvg := NewEvaluator().Evaluate(baseline, 10, Matchup{}, rng).Mean()
	mutations := deck.AllMutations(baseline, 2, true, registry.Registry{})
	// Cap mutations so the test stays under a second; full list is thousands of entries.
	if len(mutations) > 40 {
		mutations = mutations[:40]
	}

	d, _, avg, idx, found := RunMutationRound(
		context.Background(), mutations, baseAvg, 0, 0,
		30, Matchup{}, 0, 0,
		rng.Int63(), nil, nil,
	)

	if found {
		if d == nil {
			t.Error("found=true but returned deck is nil")
		}
		if avg <= baseAvg {
			t.Errorf("found=true but avg %.3f <= baseAvg %.3f", avg, baseAvg)
		}
		if idx < 0 || idx >= len(mutations) {
			t.Errorf("found=true but idx %d outside [0, %d)", idx, len(mutations))
		}
	} else {
		if d != nil {
			t.Errorf("found=false but returned deck is non-nil")
		}
		if avg != baseAvg {
			t.Errorf("found=false but avg %.3f != baseAvg %.3f", avg, baseAvg)
		}
		if idx != -1 {
			t.Errorf("found=false but idx %d != -1", idx)
		}
	}
}

// Tests that RunMutationRound returns promptly with abort semantics on context cancellation.
// Pre-cancels for deterministic behaviour.
func TestRunMutationRound_AbortsOnContextCancel(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := deck.Random(heroes.Viserai, 40, 2, rng, registry.Registry{})
	baseAvg := NewEvaluator().Evaluate(baseline, 10, Matchup{}, rng).Mean()
	mutations := deck.AllMutations(baseline, 2, true, registry.Registry{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel so no mutation ever completes its shallow eval

	var tested atomic.Int64
	start := time.Now()
	d, _, avg, idx, found := RunMutationRound(
		ctx, mutations, baseAvg, 0, 0,
		1000, Matchup{}, 0, 0,
		rng.Int63(), &tested, nil,
	)
	elapsed := time.Since(start)

	if found {
		t.Errorf("found=true after cancel; want false (aborted)")
	}
	if d != nil {
		t.Error("deck is non-nil after cancel; want nil")
	}
	if avg != baseAvg {
		t.Errorf("avg=%.3f after cancel; want baseAvg=%.3f", avg, baseAvg)
	}
	if idx != -1 {
		t.Errorf("idx=%d after cancel; want -1", idx)
	}
	// A full-list shallow screen would take multiple seconds. Abort should land in well under a
	// second — the exact bound is loose to tolerate scheduler noise, but anything near the no-
	// cancellation runtime indicates we're not actually aborting.
	if elapsed > 3*time.Second {
		t.Errorf("RunMutationRound returned after %s; abort should be near-instant", elapsed)
	}
}

// Tests that RunMutationRound terminates promptly with no improvement when bestAvg is
// unreachable — workers drain the queue without a serial deep-confirm bottleneck.
func TestRunMutationRound_TerminatesWithNoImprovement(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := deck.Random(heroes.Viserai, 40, 2, rng, registry.Registry{})
	mutations := deck.AllMutations(baseline, 2, true, registry.Registry{})
	// Cap the mutation list so the test stays well under the hang-regression threshold even on
	// slower CI runners. Full mutation list is thousands of entries.
	if len(mutations) > 40 {
		mutations = mutations[:40]
	}

	start := time.Now()
	d, _, avg, idx, found := RunMutationRound(
		context.Background(), mutations, 1_000_000.0, 0, 0, // unreachable baseline, T=0
		100, Matchup{}, 0, 0,
		rng.Int63(), nil, nil,
	)
	elapsed := time.Since(start)

	if found {
		t.Errorf("found=true with unreachable bestAvg; want false")
	}
	if d != nil {
		t.Error("deck non-nil with found=false; want nil")
	}
	if avg != 1_000_000.0 {
		t.Errorf("avg=%f; want unchanged bestAvg 1_000_000", avg)
	}
	if idx != -1 {
		t.Errorf("idx=%d; want -1", idx)
	}
	// Threshold sized for the parallel 40-mutation pass.
	if elapsed > 300*time.Second {
		t.Errorf("RunMutationRound returned after %s for 40 mutations; want under 300s", elapsed)
	}
}
