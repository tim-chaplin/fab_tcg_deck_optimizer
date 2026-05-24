package sim

// Mutation-round runner: RunMutationRound walks a list of candidate mutations and applies
// the Metropolis acceptance rule on each, returning the first acceptance. The function is
// generic — any "evaluate this batch of decks under SA-style acceptance" workload fits.
//
// Two independent parallelism knobs:
//
//   - mutationWorkers fans the mutation queue across N goroutines pulling from a shared
//     FIFO channel; each worker evaluates one mutation end-to-end. The round's Cache is
//     shared so lookup work pools across mutations.
//   - shuffleWorkers fans every per-mutation Evaluate's shuffle loop across W goroutines.
//
// (mutationWorkers, shuffleWorkers) shapes:
//   - (1, W): mutations run sequentially, shuffles fan across W. Single-deck workloads.
//   - (M, 1): mutations run in parallel, each shuffle-single-threaded. Anneal default
//     when M independent evals fit the core count.
//   - (M, W): both layers active. M×W can exceed core count — oversubscription useful
//     when shuffle workers stall on a barrier and free up cores for sibling mutations.

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/klauspost/cpuid/v2"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// mutationImprovement is the per-acceptance message sent to the coordinator: the winning
// mutation index, its evaluated average, the resulting deck, and its full deck.Stats.
type mutationImprovement struct {
	idx       int
	avg       float64
	candidate *deck.Deck
	stats     deck.Stats
}

// deckEvalConfig bundles the read-only worker parameters into one struct.
type deckEvalConfig struct {
	mutations      []deck.Mutation
	bestAvg        float64
	temperature    float64
	minImprovement float64
	shuffles       int
	matchup        Matchup
	shuffleWorkers int
	seed           int64
	completed      *atomic.Int64
	cache          *Cache
	adaptive       bool
	precision      float64
}

// RunMutationRound runs one mutation-round pass. mutationWorkers goroutines pull mutation
// indices from a shared FIFO queue, evaluate each with a per-worker Evaluator pointing at
// the round's shared Cache, and apply the Metropolis acceptance gate. The first worker to
// land an acceptable mutation wins; the others are cancelled.
//
// Defaults (passed as 0): mutationWorkers=1, shuffleWorkers=DefaultWorkers(). The
// (1, DefaultWorkers()) shape wins the worker_sweep benchmark on adaptive Viserai by
// ~20% over (DefaultWorkers(), 1) — sequential mutations let the cache fill with one
// deck's hand multisets at a time (~70% hit rate within a mutation), and the per-shuffle
// barrier balances variance better than the per-mutation queue.
//
// Cache: every worker shares one round-scoped Cache, so cross-mutation hand multisets
// share work. The cache's RWMutex serialises stores but lookups remain parallel. The
// cache is dropped when the function returns, capping memory at one round's worth.
//
// Annealing: at temperature == 0 only strict improvements clearing the minImprovement
// margin are accepted (hill climb with noise floor). At temperature > 0 worse mutations
// are also accepted with probability exp((avg - baseline) / temperature) — the
// Metropolis-style SA gate bypasses minImprovement so the SA walk retains its
// escape-local-maxima behaviour even when the floor is non-zero.
//
// minImprovement is the noise floor on strict improvements. Pass 0 to disable.
//
// FIFO pull order makes earliest-position-wins generally hold, but a stuck worker at
// position 20 doesn't block position 25 — a later position can occasionally win first.
//
// bestAvg is the SA "current state" (not the all-time best). seed is a base; worker w uses
// a derived stream. completed is an optional live-progress counter. adaptive=true stops
// per-mutation evals early at the requested precision (SE ≤ precision/4), capped by
// deck.adaptiveShufflesCap. precision is ignored when adaptive=false.
//
// Returns (acceptedDeck, acceptedStats, acceptedAvg, acceptedIndex, true) on first
// acceptance, or (nil, zero, bestAvg, -1, false) on no-acceptance or ctx cancellation.
//
// cache, when non-nil, is the hand-eval cache used by every shuffle in this round.
// Nil constructs a fresh per-round cache.
func RunMutationRound(
	ctx context.Context,
	mutations []deck.Mutation,
	bestAvg float64,
	temperature float64,
	minImprovement float64,
	shuffles int,
	mp Matchup,
	mutationWorkers, shuffleWorkers int,
	seed int64,
	completed *atomic.Int64,
	adaptive bool,
	precision float64,
	cache *Cache,
) (*deck.Deck, deck.Stats, float64, int, bool) {
	if mutationWorkers <= 0 {
		// 1 mutation worker is the empirical default — see the BenchmarkAnnealWorkerSweep
		// table on the RunMutationRound docstring for the rationale.
		mutationWorkers = 1
	}
	if shuffleWorkers <= 0 {
		shuffleWorkers = defaultWorkers()
	}
	if len(mutations) == 0 {
		return nil, deck.Stats{}, bestAvg, -1, false
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Buffer sized to mutationWorkers so every worker can land an acceptance without
	// blocking even if all of them race past cancellation simultaneously — only the
	// first one drained by the coordinator's select wins, the rest are GC'd with the
	// channel.
	improvementCh := make(chan mutationImprovement, mutationWorkers)

	jobs := make(chan int, len(mutations))
	for i := range mutations {
		jobs <- i
	}
	close(jobs)

	if cache == nil {
		cache = NewCache()
	}
	cfg := deckEvalConfig{
		mutations:      mutations,
		bestAvg:        bestAvg,
		temperature:    temperature,
		minImprovement: minImprovement,
		shuffles:       shuffles,
		matchup:        mp,
		shuffleWorkers: shuffleWorkers,
		seed:           seed,
		completed:      completed,
		cache:          cache,
		adaptive:       adaptive,
		precision:      precision,
	}

	var wg sync.WaitGroup
	for w := 0; w < mutationWorkers; w++ {
		wg.Add(1)
		go func(workerIdx int) {
			defer wg.Done()
			runDeckEvalWorker(innerCtx, cancel, workerIdx, cfg, jobs, improvementCh)
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
		return imp.candidate, imp.stats, imp.avg, imp.idx, true
	case <-workersDone:
		// A last-moment acceptance may have landed just before all senders returned.
		select {
		case imp := <-improvementCh:
			return imp.candidate, imp.stats, imp.avg, imp.idx, true
		default:
		}
		return nil, deck.Stats{}, bestAvg, -1, false
	}
}

// runDeckEvalWorker pulls mutation indices from jobs, evaluates the corresponding deck,
// and on a passing result sends a mutationImprovement and cancels the shared context.
// Each worker owns private scratch but shares the round's Cache. Returns when jobs is
// drained or the context is cancelled.
func runDeckEvalWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	workerIdx int,
	cfg deckEvalConfig,
	jobs <-chan int,
	improvementCh chan<- mutationImprovement,
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
		d := mut.Deck.Copy()
		var stats deck.Stats
		if cfg.adaptive {
			stats = ev.EvaluateAdaptive(d, cfg.precision, cfg.matchup, rng)
		} else {
			stats = ev.Evaluate(d, cfg.shuffles, cfg.matchup, rng)
		}
		avg := stats.Mean()
		if cfg.completed != nil {
			cfg.completed.Add(1)
		}
		if !acceptMutation(avg, cfg.bestAvg, cfg.temperature, cfg.minImprovement, rng) {
			continue
		}
		select {
		case improvementCh <- mutationImprovement{idx: i, avg: avg, candidate: d, stats: stats}:
		default:
			// Buffer is sized to mutationWorkers, so this default fires only if every
			// peer already filled the channel — coordinator drains exactly one anyway.
		}
		cancel()
		return
	}
}

// DefaultWorkers returns the recommended worker count for this CPU-bound workload. Capping
// at physical cores beats defaulting to GOMAXPROCS by ~20% on a typical consumer CPU (SMT
// siblings fight for cache and execution units rather than adding throughput). Clamped by
// GOMAXPROCS so a lower user/cgroup override wins.
func DefaultWorkers() int { return defaultWorkers() }

func defaultWorkers() int {
	maxProcs := runtime.GOMAXPROCS(0)
	physical := cpuid.CPU.PhysicalCores
	if physical <= 0 || physical > maxProcs {
		return maxProcs
	}
	return physical
}

// acceptMutation implements the Metropolis acceptance rule with a noise-floor guard.
// Strict improvements clearing the minImprovement margin always pass. Worse-or-marginal
// mutations pass with probability exp((deepAvg - bestAvg) / T) when T > 0; at T == 0
// they're rejected (hill-climb).
//
// minImprovement guards against shuffle-noise infinite loops where repeated near-zero
// "wins" keep accepting. The probabilistic gate intentionally ignores it so SA can still
// walk through ties / shallow dips to escape local maxima.
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
