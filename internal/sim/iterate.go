package sim

// Parallel iterate-mode round runner: IterateParallel fans out candidate Mutations across
// workers, two-phase evaluates each (shallow screen then deep confirm), and applies the
// Metropolis acceptance rule so simulated-annealing temperature widens the acceptable set.

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/klauspost/cpuid/v2"
)

// iterateImprovement is the per-acceptance message sent from a worker to the coordinator:
// the mutation index that won, its deep-confirm average, and the deck-after-mutation that
// produced it.
type iterateImprovement struct {
	idx  int
	avg  float64
	deck *Deck
}

// iterateWorkerConfig bundles every read-only parameter a worker shares with its peers so
// the goroutine body can take a single struct instead of a long argument list.
type iterateWorkerConfig struct {
	mutations      []Mutation
	bestAvg        float64
	temperature    float64
	minImprovement float64
	shuffles       int
	incoming       int
	seed           int64
	completed      *atomic.Int64
	// adaptive=true swaps the eval for EvaluateAdaptiveWith (early-stop on SE target);
	// false uses EvaluateWith with shuffles as a fixed count. Callers pick the mode based
	// on whether the user pinned -shuffles for repro / apples-to-apples flows.
	adaptive bool
}

// IterateParallel runs one iterate-mode round. Workers share a queue; each goroutine
// evaluates one mutation at a time and, on a passing result, sends an iterateImprovement
// and cancels the shared context. The first worker to land an acceptable mutation wins.
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
// Mutations are pulled FIFO so the earliest-position-wins heuristic generally holds, but
// a worker locked on an eval at position 20 doesn't block position 25 — a later-position
// mutation can occasionally win if its eval finishes first.
//
// bestAvg is the current deck's avg (SA "current state", not the all-time best). seed is
// a base; worker w uses a derived stream. completed is an optional live-progress counter.
// adaptive=true makes per-mutation evals stop early when the SE target is met (capped by
// the deck package's adaptiveShufflesCap, ignoring the shuffles arg).
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

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	improvementCh := make(chan iterateImprovement, numWorkers)

	jobs := make(chan int, len(mutations))
	for i := range mutations {
		jobs <- i
	}
	close(jobs)

	cfg := iterateWorkerConfig{
		mutations:      mutations,
		bestAvg:        bestAvg,
		temperature:    temperature,
		minImprovement: minImprovement,
		shuffles:       shuffles,
		incoming:       incoming,
		seed:           seed,
		completed:      completed,
		adaptive:       adaptive,
	}

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerIdx int) {
			defer wg.Done()
			runIterateWorker(innerCtx, cancel, workerIdx, cfg, jobs, improvementCh)
		}(w)
	}

	workersDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(workersDone)
	}()

	select {
	case imp := <-improvementCh:
		<-workersDone
		return imp.deck, imp.avg, imp.idx, true
	case <-workersDone:
		// A last-moment acceptance may have landed just before all senders returned.
		select {
		case imp := <-improvementCh:
			return imp.deck, imp.avg, imp.idx, true
		default:
		}
		return nil, bestAvg, -1, false
	}
}

// runIterateWorker pulls mutation indices from jobs, evaluates each, and on a passing
// result sends an iterateImprovement and cancels the shared context. Each worker owns its
// own evaluator and rng so the inner loop is allocation-free and free of cross-worker
// coupling. Returns when jobs is drained or when the context is cancelled (either by
// another winner or by the caller).
func runIterateWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	workerIdx int,
	cfg iterateWorkerConfig,
	jobs <-chan int,
	improvementCh chan<- iterateImprovement,
) {
	ev := NewEvaluator()
	rng := rand.New(rand.NewSource(cfg.seed ^ (int64(workerIdx)+1)*int64(0x9e3779b9)))
	for i := range jobs {
		if ctx.Err() != nil {
			return
		}
		mut := cfg.mutations[i]
		d := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
		var avg float64
		if cfg.adaptive {
			avg = d.EvaluateAdaptiveWith(cfg.incoming, rng, ev).Mean()
		} else {
			avg = d.EvaluateWith(cfg.shuffles, cfg.incoming, rng, ev).Mean()
		}
		if cfg.completed != nil {
			cfg.completed.Add(1)
		}
		if !acceptMutation(avg, cfg.bestAvg, cfg.temperature, cfg.minImprovement, rng) {
			continue
		}
		select {
		case improvementCh <- iterateImprovement{idx: i, avg: avg, deck: d}:
		default:
			// Another worker already filled the buffer; drop silently.
		}
		cancel()
		return
	}
}

// defaultWorkers returns the worker count when callers pass numWorkers<=0. The workload is
// purely CPU-bound, so SMT siblings fight for cache and execution units rather than adding
// throughput: capping at physical cores outperforms defaulting to GOMAXPROCS by ~20% on a
// typical consumer CPU. Still clamped by GOMAXPROCS so a lower user/cgroup override wins.
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
