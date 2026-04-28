package sim

// Iterate-mode round runner: IterateParallel walks the candidate Mutations in priority
// order, evaluating each one with a shuffle-parallel Evaluator (one shared cache + N
// worker goroutines fanning the shuffle loop), and applies the Metropolis acceptance rule
// so simulated-annealing temperature widens the acceptable set. Returns on the first
// mutation that passes the gate.

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync/atomic"

	"github.com/klauspost/cpuid/v2"
)

// IterateParallel walks `mutations` in order, evaluating each with a shuffle-parallel
// Evaluator (cache shared across all evals; ResetCache between mutations to bound memory)
// and returns on the first acceptance. Parallelism lives inside the per-mutation eval —
// numWorkers goroutines fan out the shuffle loop — rather than across mutations, so a
// single shared cache lets every shuffle of every mutation contribute to and read from
// the same memo.
//
// Annealing: at temperature == 0 only strict improvements clearing the minImprovement
// margin are accepted (classical hill climb with a noise floor). At temperature > 0 worse
// mutations are also accepted with probability exp((avg - baseline) / temperature) — a
// Metropolis-style SA gate that bypasses the minImprovement margin entirely (so the SA
// walk retains its escape-local-maxima behaviour even when the floor is non-zero).
//
// minImprovement is the noise floor on strict improvements: a mutation must lift avg by
// more than this amount above bestAvg to be accepted at T==0. Prevents infinite loops
// where repeated near-zero "wins" (within shuffle noise) keep accepting indefinitely.
// Pass 0 to disable the floor (any strictly-greater avg passes).
//
// bestAvg is the current deck's avg (SA "current state", not the all-time best). seed
// drives the RNG used both to shuffle decks and to flip Metropolis coins. completed is
// an optional live-progress counter (bumped per mutation evaluated). adaptive=true makes
// per-mutation evals stop early when the SE target is met (capped by the deck package's
// adaptiveShufflesCap, ignoring the shuffles arg).
//
// Returns (acceptedDeck, acceptedAvg, acceptedIndex, true) on first acceptance, or
// (nil, bestAvg, -1, false) if nothing cleared the gate or ctx was cancelled.
func IterateParallel(
	ctx context.Context,
	mutations []Mutation,
	bestAvg float64,
	temperature float64,
	minImprovement float64,
	shuffles, incoming, numWorkers int,
	seed int64,
	completed *atomic.Int64,
	adaptive bool,
) (*Deck, float64, int, bool) {
	if numWorkers <= 0 {
		numWorkers = defaultWorkers()
	}
	if len(mutations) == 0 {
		return nil, bestAvg, -1, false
	}

	ev := NewEvaluatorParallel(numWorkers)
	rng := rand.New(rand.NewSource(seed))
	for i := range mutations {
		if ctx.Err() != nil {
			return nil, bestAvg, -1, false
		}
		// Drop the cache between mutations. Different mutations evaluate different decks
		// whose hand multisets rarely overlap, so retaining entries from the previous
		// deck would just bloat memory without saving search work.
		ev.ResetCache()
		mut := mutations[i]
		d := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
		var avg float64
		if adaptive {
			avg = d.EvaluateAdaptiveWith(incoming, rng, ev).Mean()
		} else {
			avg = d.EvaluateWith(shuffles, incoming, rng, ev).Mean()
		}
		if completed != nil {
			completed.Add(1)
		}
		if acceptMutation(avg, bestAvg, temperature, minImprovement, rng) {
			return d, avg, i, true
		}
	}
	return nil, bestAvg, -1, false
}

// DefaultWorkers returns the recommended worker count when callers want the production
// default. The workload is purely CPU-bound, so SMT siblings fight for cache and execution
// units rather than adding throughput: capping at physical cores outperforms defaulting
// to GOMAXPROCS by ~20% on a typical consumer CPU. Still clamped by GOMAXPROCS so a lower
// user/cgroup override wins. Exported so cmd/fabsim's modes can size their parallel
// Evaluators consistently with iterate-mode's worker pool.
func DefaultWorkers() int { return defaultWorkers() }

// defaultWorkers is the package-internal helper iterate-mode and the public DefaultWorkers
// share. Kept private so callers go through DefaultWorkers and the sizing rule has one
// docstring.
func defaultWorkers() int {
	maxProcs := runtime.GOMAXPROCS(0)
	physical := cpuid.CPU.PhysicalCores
	if physical <= 0 || physical > maxProcs {
		return maxProcs
	}
	return physical
}

// acceptMutation implements the Metropolis acceptance rule with a noise-floor guard. Strict
// improvements that clear the minImprovement margin (deepAvg > bestAvg + minImprovement)
// always pass. Worse-or-marginal mutations pass with probability exp((deepAvg - bestAvg) /
// T) when T > 0; at T == 0 they're rejected, recovering the classical hill-climb behaviour.
//
// minImprovement guards against infinite loops where shuffle noise lets repeated near-zero
// "wins" keep accepting indefinitely. The probabilistic gate intentionally ignores it so SA
// can still walk through ties / shallow dips to escape local maxima — without that bypass,
// raising the floor would shrink the explorable region of the SA walk in proportion.
func acceptMutation(deepAvg, bestAvg, temperature, minImprovement float64, rng *rand.Rand) bool {
	if deepAvg > bestAvg+minImprovement {
		return true
	}
	if temperature <= 0 {
		return false
	}
	prob := math.Exp((deepAvg - bestAvg) / temperature)
	return rng.Float64() < prob
}
