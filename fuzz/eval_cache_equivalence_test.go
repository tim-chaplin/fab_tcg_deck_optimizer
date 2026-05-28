// Package fuzz holds slow / randomised correctness sweeps and the pinned regression
// seeds they've surfaced. Tests here run on demand (in CI nightly, locally when
// investigating sim divergences) — not part of the default `go test ./...` hot path.
// External package (peer of turntests/), uses only the sim public API.
package fuzz

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Sweeps random setup seeds within a time budget, failing the first seed whose cached
// vs uncached Evaluate diverge on Hands or TotalValue. Set EVAL_FUZZ_SECONDS=N to extend.
func TestEvalCache_EquivalenceWithUncached_FuzzAutomatic(t *testing.T) {
	const (
		deckSize      = 40
		maxCopies     = 2
		incoming      = 7
		shuffles      = 20
		defaultRunFor = 10 * time.Second
	)
	runFor := defaultRunFor
	if s := os.Getenv("EVAL_FUZZ_SECONDS"); s != "" {
		if n, err := time.ParseDuration(s + "s"); err == nil {
			runFor = n
		}
	}
	// EVAL_FUZZ_SETUP_SEED pins setupSeed to a specific value and runs exactly one
	// iteration — used to repro a divergence the random sweep surfaced. The reproduction
	// recipe printed in reportFuzzDivergence sets this variable to the failing seed.
	var pinnedSeed int64
	pinned := false
	if s := os.Getenv("EVAL_FUZZ_SETUP_SEED"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			pinnedSeed = n
			pinned = true
		}
	}
	deadline := time.Now().Add(runFor)
	seedGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	tested := 0
	for time.Now().Before(deadline) {
		var setupSeed int64
		if pinned {
			setupSeed = pinnedSeed
		} else {
			setupSeed = seedGen.Int63()
		}
		setupRNG := rand.New(rand.NewSource(setupSeed))
		baseline := deck.Random(heroes.Viserai, deckSize, maxCopies, setupRNG, registry.Registry{})

		cachedStats := sim.NewEvaluator().Evaluate(baseline.Copy(), shuffles, sim.Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))
		uncachedStats := sim.NewEvaluatorWithoutCache().Evaluate(baseline.Copy(), shuffles, sim.Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))
		tested++

		if cachedStats.Hands != uncachedStats.Hands || cachedStats.TotalValue != uncachedStats.TotalValue {
			reportFuzzDivergence(t, setupSeed, tested, cachedStats, uncachedStats,
				baseline, deckSize, maxCopies, incoming, shuffles)
		}
		if pinned {
			break
		}
	}
	t.Logf("swept %d random setup seeds in %v with no divergence", tested, runFor)
}

// reportFuzzDivergence formats a divergence as a t.Fatalf with the reproduction recipe
// surfaced up top. The leading banner exists because randomized fuzz failures are easy
// to mis-read as transient when a follow-up run passes — the divergence is real and
// pinned to the seed printed here, not to the runtime conditions of either run.
func reportFuzzDivergence(t *testing.T, setupSeed int64, tested int, cached, uncached deck.Stats,
	baseline *deck.Deck, deckSize, maxCopies, incoming, shuffles int) {
	t.Helper()
	divergeAt, cAtK, uAtK := bisectDivergentShuffle(baseline, incoming, shuffles)
	var narrowing strings.Builder
	if divergeAt <= 0 {
		fmt.Fprintf(&narrowing, "Bisection failed: re-run did not reproduce divergence at any shuffle count up to %d.\n", shuffles)
	} else {
		fmt.Fprintf(&narrowing, "Bisection: divergence first appears at shuffles=%d (cached.TotalValue=%.0f uncached.TotalValue=%.0f delta=%.0f).\n",
			divergeAt, cAtK.TotalValue, uAtK.TotalValue, cAtK.TotalValue-uAtK.TotalValue)
		if divergeAt == 1 {
			fmt.Fprintf(&narrowing, "The bug surfaces in the very first shuffle; investigate by tracing per-Best behaviour through shuffle 1 on a fresh evaluator.\n")
		} else {
			fmt.Fprintf(&narrowing, "Shuffles 1..%d are clean; the bug surfaces during shuffle %d's turns. Investigate by tracing per-Best behaviour through that shuffle on a fresh evaluator.\n", divergeAt-1, divergeAt)
		}
	}

	t.Fatalf(`THIS IS NOT A TRANSIENT ERROR, even if rerunning this test succeeds; this intentionally tests different sets of inputs on every run.

To repeat this failure, run:

  EVAL_FUZZ_SETUP_SEED=%d go test -count=1 -run TestEvalCache_EquivalenceWithUncached_FuzzAutomatic ./fuzz/

Divergence at setupSeed=%d (after %d seeds):
  Hands:      cached=%d  uncached=%d
  TotalValue: cached=%.0f uncached=%.0f

%s`,
		setupSeed, setupSeed, tested,
		cached.Hands, uncached.Hands,
		cached.TotalValue, uncached.TotalValue,
		narrowing.String(),
	)
}

// bisectDivergentShuffle finds the smallest shuffle count K (1..maxShuffles) at which
// cached vs uncached Evaluate disagrees on Hands or TotalValue. Returns (K, statsAtK
// for cached, statsAtK for uncached) when a divergent boundary exists, or (0, _, _)
// when the bisection didn't reproduce. O(log maxShuffles) Evaluate-pair calls.
func bisectDivergentShuffle(baseline *deck.Deck, incoming, maxShuffles int) (int, deck.Stats, deck.Stats) {
	diverges := func(n int) (deck.Stats, deck.Stats, bool) {
		c := sim.NewEvaluator().Evaluate(baseline.Copy(), n, sim.Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))
		u := sim.NewEvaluatorWithoutCache().Evaluate(baseline.Copy(), n, sim.Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))
		return c, u, c.Hands != u.Hands || c.TotalValue != u.TotalValue
	}
	if c, u, ok := diverges(maxShuffles); !ok {
		return 0, c, u
	}
	lo, hi := 1, maxShuffles
	var cAtK, uAtK deck.Stats
	for lo < hi {
		mid := (lo + hi) / 2
		c, u, ok := diverges(mid)
		if ok {
			hi = mid
			cAtK, uAtK = c, u
		} else {
			lo = mid + 1
		}
	}
	if cAtK.Runs == 0 {
		cAtK, uAtK, _ = diverges(lo)
	}
	return lo, cAtK, uAtK
}
