package sim

// Iterate-mode round runner: IterateParallel walks candidate Mutations and applies the
// Metropolis acceptance rule on each, returning the first acceptance. Two independent
// parallelism knobs let callers dial each dimension to the workload:
//
//   - mutationWorkers fans the mutation queue across N worker goroutines that pull from a
//     shared FIFO channel. Each worker evaluates one mutation at a time end-to-end. The
//     round's Cache is shared across every worker so lookup work pools across mutations.
//   - shuffleWorkers fans every per-mutation Evaluate call's shuffle loop across W worker
//     goroutines (per the parallel path in deck_eval.go's evaluateParallelImpl).
//
// (mutationWorkers, shuffleWorkers) shapes:
//   - (1, W): the iterate.go-equivalent of a single-deck eval — every mutation runs
//     sequentially, but each mutation's shuffle loop fans across W goroutines. Right shape
//     for fabsim eval (one deck → no mutation parallelism available).
//   - (M, 1): every mutation runs in parallel, each one shuffle-single-threaded. Right
//     shape for anneal when M decks worth of independent work fits the core count.
//   - (M, W): both layers active. The product M×W can exceed the core count — useful for
//     experimenting with oversubscription, where shuffle workers stalled on a barrier free
//     up cores for sibling mutation workers.

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/klauspost/cpuid/v2"
)

// iterateImprovement is the per-acceptance message a worker sends to the coordinator: the
// mutation index that won, its evaluated average, and the deck-after-mutation that
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
	arcaneIncoming int
	shuffleWorkers int
	seed           int64
	completed      *atomic.Int64
	cache          *Cache
	adaptive       bool
}

// IterateParallel runs one iterate-mode round. mutationWorkers goroutines pull mutation
// indices from a shared queue, evaluate each one with a per-worker Evaluator that points
// at the round's shared Cache, and apply the Metropolis acceptance gate. The first worker
// to land an acceptable mutation wins; the others are cancelled.
//
// shuffleWorkers controls the per-mutation Evaluate's parallelism: 0 or 1 runs the
// shuffle loop single-threaded; >1 fans the shuffle loop across that many goroutines per
// mutation.
//
// Defaults (passed as 0): mutationWorkers=1, shuffleWorkers=DefaultWorkers(). The
// (1, DefaultWorkers()) shape wins the worker_sweep benchmark on adaptive Viserai by
// ~20% over (DefaultWorkers(), 1) — sequential mutations let the cache fill with one
// deck's hand multisets at a time (~70% hit rate within a mutation), and the per-shuffle
// barrier balances variance better than the per-mutation queue does. Pass
// mutationWorkers=N explicitly to override (e.g. for experiments or workloads where the
// cache is disabled).
//
// Cache: every worker constructs its Evaluator via NewEvaluatorWithCache pointing at one
// Cache built locally for this round. Workers' lookups and stores all hit the same memo
// so cross-mutation hand multisets share work; the cache's RWMutex serialises stores but
// lookups remain parallel. The round's cache lives only for IterateParallel's lifetime,
// so memory growth from accumulating one round's worth of distinct hand multisets is
// capped at the round boundary — the function returns, the cache pointer drops, the
// next round starts fresh.
//
// Annealing: at temperature == 0 only strict improvements clearing the minImprovement
// margin are accepted (classical hill climb with a noise floor). At temperature > 0 worse
// mutations are also accepted with probability exp((avg - baseline) / temperature) — a
// Metropolis-style SA gate that bypasses the minImprovement margin entirely (so the SA
// walk retains its escape-local-maxima behaviour even when the floor is non-zero).
//
// minImprovement is the noise floor on strict improvements. Pass 0 to disable.
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
	shuffles, incoming, arcaneIncoming int,
	mutationWorkers, shuffleWorkers int,
	seed int64,
	completed *atomic.Int64,
	adaptive bool,
) (*Deck, float64, int, bool) {
	if mutationWorkers <= 0 {
		// 1 mutation worker is the empirical default — see the BenchmarkAnnealWorkerSweep
		// table on the IterateParallel docstring for the rationale.
		mutationWorkers = 1
	}
	if shuffleWorkers <= 0 {
		shuffleWorkers = defaultWorkers()
	}
	if len(mutations) == 0 {
		return nil, bestAvg, -1, false
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Buffer sized to mutationWorkers so every worker can land an acceptance without
	// blocking even if all of them race past cancellation simultaneously — only the
	// first one drained by the coordinator's select wins, the rest are GC'd with the
	// channel.
	improvementCh := make(chan iterateImprovement, mutationWorkers)

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
		arcaneIncoming: arcaneIncoming,
		shuffleWorkers: shuffleWorkers,
		seed:           seed,
		completed:      completed,
		cache:          NewCache(),
		adaptive:       adaptive,
	}

	var wg sync.WaitGroup
	for w := 0; w < mutationWorkers; w++ {
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
// own per-call scratch (Evaluator with private attackBufs) but points to the round's
// shared Cache so lookup work pools across all workers' mutation evals. The Evaluator's
// numWorkers is set to cfg.shuffleWorkers so the per-mutation eval may layer shuffle-
// level fan-out on top. Returns when jobs is drained or the context is cancelled.
func runIterateWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	workerIdx int,
	cfg iterateWorkerConfig,
	jobs <-chan int,
	improvementCh chan<- iterateImprovement,
) {
	ev := NewEvaluatorWithCache(cfg.cache)
	if cfg.shuffleWorkers > 1 {
		ev.numWorkers = cfg.shuffleWorkers
	}
	rng := rand.New(rand.NewSource(cfg.seed ^ (int64(workerIdx)+1)*int64(0x9e3779b9)))
	for i := range jobs {
		if ctx.Err() != nil {
			return
		}
		mut := cfg.mutations[i]
		d := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
		var avg float64
		if cfg.adaptive {
			avg = d.EvaluateAdaptiveWith(cfg.incoming, cfg.arcaneIncoming, rng, ev).Mean()
		} else {
			avg = d.EvaluateWith(cfg.shuffles, cfg.incoming, cfg.arcaneIncoming, rng, ev).Mean()
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
			// Buffer is sized to mutationWorkers, so this default fires only if every
			// peer already filled the channel — coordinator drains exactly one anyway.
		}
		cancel()
		return
	}
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
