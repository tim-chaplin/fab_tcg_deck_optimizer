package sim

// Mutation-round runner: RunMutationRound walks candidate mutations and applies the
// Metropolis acceptance rule, returning the first acceptance.
//
// Two parallelism knobs:
//   - mutationWorkers fans the mutation queue across N goroutines (shared FIFO channel);
//     each worker evaluates one mutation end-to-end. Round Cache is shared.
//   - shuffleWorkers fans every per-mutation Evaluate's shuffle loop across W goroutines.
//
// Shapes: (1, W) for single-deck workloads, (M, 1) for parallel mutations, (M, W) for both
// (oversubscription helps when shuffle workers stall on a barrier).

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
}

// RunMutationRound runs one mutation-round pass. mutationWorkers pull indices from a shared
// FIFO queue, evaluate each via a per-worker Evaluator on the round's shared Cache, and
// apply the Metropolis acceptance gate. First acceptance wins; the others cancel.
//
// Defaults (0): mutationWorkers=1, shuffleWorkers=DefaultWorkers(). The (1, DefaultWorkers())
// shape wins the worker_sweep benchmark on a Viserai deck by ~20% — sequential mutations
// let the cache fill with one deck's hand multisets at a time (~70% hit rate within a
// mutation), and the per-shuffle barrier balances variance better than per-mutation queues.
//
// Cache is round-scoped and shared across workers, then dropped on return. RWMutex
// serialises stores; lookups remain parallel.
//
// Annealing: at temperature == 0 only strict improvements clearing minImprovement pass
// (hill climb with noise floor). At temperature > 0, worse mutations pass with probability
// exp((avg - baseline) / temperature) — the SA gate bypasses minImprovement to keep
// escape-local-maxima behaviour. minImprovement=0 disables the floor.
//
// FIFO pull order makes earliest-position-wins usually hold, but a stuck early worker
// doesn't block later positions from occasionally winning first.
//
// bestAvg is the SA current state (not all-time best): the baseline each mutation's avg is
// compared against. Every mutation is evaluated with the fixed shuffles budget, shuffled from
// seed — so the comparison is paired (common random numbers) when bestAvg was itself measured
// on this seed and shuffleWorkers, and an ordinary independent comparison otherwise. cache,
// when non-nil, is the hand-eval cache for every shuffle this round; nil allocates one.
//
// Returns (acceptedDeck, acceptedStats, acceptedAvg, acceptedIndex, true) on first
// acceptance, or (nil, zero, bestAvg, -1, false) on no acceptance / cancellation.
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
	cache *Cache,
) (*deck.Deck, deck.Stats, float64, int, bool) {
	if mutationWorkers <= 0 {
		// 1 worker is the empirical default — see the RunMutationRound docstring.
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

	// Buffer = mutationWorkers so every worker can land an acceptance without blocking
	// even if all race past cancellation; only the first drain by the coordinator wins.
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

// runDeckEvalWorker pulls mutation indices from jobs, evaluates each deck, and on a passing
// result sends a mutationImprovement and cancels the shared context. Private scratch,
// shared cache. Returns on jobs drain or context cancellation.
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
	// Every mutation shuffles from cfg.seed, the same seed the caller used for the incumbent
	// (cfg.bestAvg). Common random numbers: the shuffle noise shared by incumbent and mutant
	// cancels in the avg-vs-baseline comparison, so a real improvement resolves in far fewer
	// shuffles, and the shared 39 cards draw identically across mutants to lift the cache hit
	// rate. shuffleRNG is reseeded per mutation; coinRNG is a separate per-worker stream so the
	// SA acceptance coin still varies between mutations.
	shuffleRNG := rand.New(rand.NewSource(cfg.seed))
	coinRNG := rand.New(rand.NewSource(cfg.seed ^ (int64(workerIdx)+1)*int64(0x9e3779b9)))
	for i := range jobs {
		if ctx.Err() != nil {
			return
		}
		mut := cfg.mutations[i]
		d := mut.Deck()
		shuffleRNG.Seed(cfg.seed)
		stats := ev.Evaluate(d, cfg.shuffles, cfg.matchup, shuffleRNG)
		avg := stats.Mean()
		if cfg.completed != nil {
			cfg.completed.Add(1)
		}
		if !acceptMutation(avg, cfg.bestAvg, cfg.temperature, cfg.minImprovement, coinRNG) {
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

// DefaultWorkers returns the recommended worker count: physical cores (SMT siblings cost
// ~20% throughput on this CPU-bound workload), clamped by GOMAXPROCS so user / cgroup
// overrides win.
func DefaultWorkers() int { return defaultWorkers() }

func defaultWorkers() int {
	maxProcs := runtime.GOMAXPROCS(0)
	physical := cpuid.CPU.PhysicalCores
	if physical <= 0 || physical > maxProcs {
		return maxProcs
	}
	return physical
}

// acceptMutation is the Metropolis acceptance rule with a noise-floor guard. Strict
// improvements clearing minImprovement always pass. Worse / marginal pass with
// probability exp((deepAvg - bestAvg) / T) when T > 0, else rejected.
//
// minImprovement guards against shuffle-noise loops of near-zero "wins". The probabilistic
// gate ignores it so SA can still walk through ties / shallow dips.
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
