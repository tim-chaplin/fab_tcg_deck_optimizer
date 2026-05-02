package sim_test

import (
	"context"
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
)

// TestEvaluateWith_ConcurrentNoMapPanic hammers EvaluateWith from many goroutines with
// each holding its own Evaluator. Shared state includes the card-meta lookup table;
// unsynchronised access would panic with "concurrent map read and map write" or "concurrent
// map writes" and fail the test. Catches common race regressions without depending on -race
// (which requires cgo).
func TestEvaluateWith_ConcurrentNoMapPanic(t *testing.T) {
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 2 {
		t.Skip("need GOMAXPROCS >= 2 to exercise concurrent access")
	}
	const iterations = 25

	baseline := Random(heroes.Viserai{}, 40, 2, rand.New(rand.NewSource(42)), nil)

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ev := NewEvaluator()
			rng := rand.New(rand.NewSource(int64(id)*7919 + 1))
			for i := 0; i < iterations; i++ {
				d := New(baseline.Hero, baseline.Weapons, baseline.Cards)
				// Small shuffle count per iteration keeps the test fast while still exercising
				// many Best calls per goroutine: 10 shuffles × handsPerCycle * 2 hands each =
				// ~200 Best invocations per goroutine, per iteration.
				stats := d.EvaluateWith(10, Matchup{}, rng, ev)
				if stats.Hands == 0 {
					t.Errorf("worker %d iter %d: Evaluate returned zero hands", id, i)
					return
				}
			}
		}(w)
	}
	wg.Wait()
}

// TestIterateParallel_RunsWithoutPanic is a smoke test for the parallel iterate entry
// point: verifies the worker pool, cancellation signalling, and ready-channel coordination
// complete without deadlock or panic on a realistic mutation list. Whether an improvement is
// found is seed-dependent, so the test only asserts invariants.
func TestIterateParallel_RunsWithoutPanic(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, 40, 2, rng, nil)
	baseAvg := baseline.Evaluate(10, Matchup{}, rng).Mean()
	mutations := AllMutations(baseline, 2, nil)
	// Cap mutations so the test stays under a second; full list is thousands of entries.
	if len(mutations) > 40 {
		mutations = mutations[:40]
	}

	d, avg, idx, found := IterateParallel(
		context.Background(), mutations, baseAvg, 0, 0,
		30, Matchup{}, 0, 0,
		rng.Int63(), nil, false,
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

// TestIterateParallel_AbortsOnContextCancel pins the abort path that iterate's stdin-listener
// depends on: a context cancellation must unblock both the worker pool and the main-goroutine
// select on ready[i], and IterateParallel must return promptly with found=false and ctx.Err()
// set.
//
// Pre-cancels the context so the outcome is deterministic regardless of which worker happens to
// deep-confirm first — the interesting assertion is "returns promptly with abort semantics,"
// not "cancel races vs shallow completion."
func TestIterateParallel_AbortsOnContextCancel(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, 40, 2, rng, nil)
	baseAvg := baseline.Evaluate(10, Matchup{}, rng).Mean()
	mutations := AllMutations(baseline, 2, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel so no mutation ever completes its shallow eval

	var tested atomic.Int64
	start := time.Now()
	d, avg, idx, found := IterateParallel(
		ctx, mutations, baseAvg, 0, 0,
		1000, Matchup{}, 0, 0,
		rng.Int63(), &tested, false,
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
		t.Errorf("IterateParallel returned after %s; abort should be near-instant", elapsed)
	}
}

// TestIterateParallel_TerminatesWithNoImprovement pins prompt return when no mutation
// confirms: workers drain the shared queue with no serial deep-confirm bottleneck. Uses an
// artificially high bestAvg so every mutation fails the shallow screen cleanly AND any
// noise-driven shallow passer fails deep confirmation too — reliably hits the
// "drain-queue-no-improvement-found" path.
func TestIterateParallel_TerminatesWithNoImprovement(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, 40, 2, rng, nil)
	mutations := AllMutations(baseline, 2, nil)
	// Cap the mutation list so the test stays well under the hang-regression threshold even on
	// slower CI runners. Full mutation list is thousands of entries.
	if len(mutations) > 40 {
		mutations = mutations[:40]
	}

	start := time.Now()
	d, avg, idx, found := IterateParallel(
		context.Background(), mutations, 1_000_000.0, 0, 0, // unreachable baseline, T=0
		100, Matchup{}, 0, 0,
		rng.Int63(), nil, false,
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
	if elapsed > 30*time.Second {
		t.Errorf("IterateParallel returned after %s for 40 mutations; want under 30s", elapsed)
	}
}
