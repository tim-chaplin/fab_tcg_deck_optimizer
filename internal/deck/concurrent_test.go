package deck

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestEvaluateWith_ConcurrentNoMapPanic hammers deck.EvaluateWith from many goroutines with each
// holding its own hand.Evaluator. Shared state across goroutines includes the memo map and the
// card-meta lookup table; without proper synchronisation the Go runtime would panic with
// "concurrent map read and map write" or "concurrent map writes" and fail the test. This catches
// the most common category of race regressions without depending on -race (which requires cgo,
// and we don't ship a gcc toolchain on CI).
func TestEvaluateWith_ConcurrentNoMapPanic(t *testing.T) {
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 2 {
		t.Skip("need GOMAXPROCS >= 2 to exercise concurrent access")
	}
	const iterations = 25

	baseline := Random(hero.Viserai{}, 40, 2, rand.New(rand.NewSource(42)))

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ev := hand.NewEvaluator()
			rng := rand.New(rand.NewSource(int64(id)*7919 + 1))
			for i := 0; i < iterations; i++ {
				d := New(baseline.Hero, baseline.Weapons, baseline.Cards)
				// Small shuffle count per iteration keeps the test fast while still exercising
				// many Best calls per goroutine: 10 shuffles × handsPerCycle * 2 hands each =
				// ~200 Best invocations per goroutine, per iteration.
				stats := d.EvaluateWith(10, 0, rng, ev)
				if stats.Hands == 0 {
					t.Errorf("worker %d iter %d: Evaluate returned zero hands", id, i)
					return
				}
			}
		}(w)
	}
	wg.Wait()
}

// TestIterateParallel_RunsWithoutPanic is a smoke test for the parallel iterate entry point. It
// verifies the worker pool, cancellation signalling, and ready-channel coordination all complete
// without deadlock or panic on a realistic mutation list. Whether an improvement is found for
// this specific seed is not guaranteed, so the test only asserts invariants — not that the
// search succeeded.
func TestIterateParallel_RunsWithoutPanic(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	baseline := Random(hero.Viserai{}, 40, 2, rng)
	baseAvg := baseline.Evaluate(10, 0, rng).Avg()
	mutations := AllMutations(baseline, 2)
	// Cap mutations so the test stays under a second; full list is thousands of entries.
	if len(mutations) > 40 {
		mutations = mutations[:40]
	}

	d, avg, idx, found := IterateParallel(
		mutations, baseAvg, 10, 30, 0, 0,
		rng.Int63(), rng, nil,
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
