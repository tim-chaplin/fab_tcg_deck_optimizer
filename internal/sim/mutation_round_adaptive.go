package sim

// Adaptive mutation-round runner. Each candidate is screened by an SPRT on coupled per-shuffle
// value deltas (see sprt.go) instead of a fixed shuffle budget: clear rejects are dropped in a
// handful of shuffles, only close calls run long, and there is no max-shuffle cap. A mutation the
// screen accepts is then confirmed with a high-precision coupled eval at statsShuffles before it
// can win the round — so a screen false-accept can't win, which is what guarantees the round
// terminates. The accept threshold is min-improvement for a hill climb and the random Metropolis
// threshold for an SA step, so this runner covers all temperatures.
//
// Parallelism here is across mutations (one worker per goroutine), not within an eval: a per-
// shuffle sequential screen can't fan a single shuffle across workers. The incumbent's per-shuffle
// values are shared across workers through incumbentValues.

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// perShuffleSeed derives a distinct, well-mixed seed for shuffle n (1-based) of a round. The
// incumbent and every mutation use perShuffleSeed(roundSeed, n) for their n-th screen shuffle, so
// the n-th shuffles are coupled (identical Fisher-Yates draws on position-aligned decks). n==0 is
// reserved for the confirm eval. splitmix64 finalizer over (roundSeed, n).
func perShuffleSeed(roundSeed int64, n int) int64 {
	x := uint64(roundSeed) + uint64(n)*0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return int64(x ^ (x >> 31))
}

// singleShuffleEval evaluates one shuffle of a fixed deck at a time, reusing its scratch and
// id-index across shuffles so a screen that runs many shuffles on one deck doesn't re-derive them
// each time.
type singleShuffleEval struct {
	d             *deck.Deck
	ev            *Evaluator
	mp            Matchup
	scratch       *shuffleScratch
	idIndex       map[ids.CardID]int
	handsPerCycle int
	handSize      int
	ok            bool
}

func newSingleShuffleEval(d *deck.Deck, ev *Evaluator, mp Matchup) *singleShuffleEval {
	handSize := d.Hero.(hero.Hero).Intelligence()
	deckSize := d.Size()
	if handSize <= 0 || deckSize < handSize {
		return &singleShuffleEval{ok: false}
	}
	uniqueIDs, idIndex := d.UniqueIDs()
	return &singleShuffleEval{
		d:             d,
		ev:            ev,
		mp:            mp,
		scratch:       newShuffleScratch(len(d.Weapons), deckSize, handSize, len(uniqueIDs)),
		idIndex:       idIndex,
		handsPerCycle: deckSize / handSize,
		handSize:      handSize,
		ok:            true,
	}
}

// value runs one shuffle seeded from seed and returns that shuffle's mean hand value.
func (e *singleShuffleEval) value(rng *rand.Rand, seed int64) float64 {
	if !e.ok {
		return 0
	}
	rng.Seed(seed)
	var s deck.Stats
	runOneShuffle(e.d, e.scratch, &s, e.idIndex, e.ev, rng, e.mp, e.handsPerCycle, e.handSize)
	return s.Mean()
}

// incumbentValues lazily computes and caches the incumbent's per-shuffle screen values, shared
// across all mutation workers in a round. The lock serialises the rare frontier extension;
// already-computed depths return immediately.
type incumbentValues struct {
	mu   sync.Mutex
	eval *singleShuffleEval
	rng  *rand.Rand
	seed int64
	vals []float64
}

func newIncumbentValues(incumbent *deck.Deck, ev *Evaluator, mp Matchup, seed int64) *incumbentValues {
	return &incumbentValues{
		eval: newSingleShuffleEval(incumbent, ev, mp),
		rng:  rand.New(rand.NewSource(1)),
		seed: seed,
	}
}

// at returns the incumbent's value on shuffle n (1-based), computing any not-yet-seen shuffles.
func (iv *incumbentValues) at(n int) float64 {
	iv.mu.Lock()
	defer iv.mu.Unlock()
	for len(iv.vals) < n {
		next := len(iv.vals) + 1
		iv.vals = append(iv.vals, iv.eval.value(iv.rng, perShuffleSeed(iv.seed, next)))
	}
	return iv.vals[n-1]
}

// confirmer lazily computes the incumbent's averaged value over shuffles at the confirm seed,
// once. A mutant confirm evaluated at the same seed (sequential) is coupled with it, so their
// difference carries the reduced common-random-numbers variance.
type confirmer struct {
	once      sync.Once
	incumbent *deck.Deck
	mp        Matchup
	shuffles  int
	seed      int64
	cache     *Cache
	incAvg    float64
}

func (c *confirmer) incumbentAvg() float64 {
	c.once.Do(func() {
		ev := NewEvaluatorWithCache(c.cache)
		c.incAvg = ev.Evaluate(c.incumbent, c.shuffles, c.mp, rand.New(rand.NewSource(c.seed))).Mean()
	})
	return c.incAvg
}

// AdaptiveProgress carries the live counters the anneal ticker reads while a round runs: how many
// candidates have reached a screen decision, and how many high-precision confirms are in flight
// (non-zero means the round paused to full-evaluate a promising candidate).
type AdaptiveProgress struct {
	Screened   atomic.Int64
	Confirming atomic.Int64
}

// AdaptiveRoundConfig bundles the inputs to RunMutationRoundAdaptive. The caller builds a fresh one
// each round; ctx stays a separate argument per the Go convention.
type AdaptiveRoundConfig struct {
	Mutations       []deck.Mutation
	Incumbent       *deck.Deck
	Threshold       float64 // the min-improvement gate
	Temperature     float64 // 0 = hill climb; > 0 = simulated annealing
	Matchup         Matchup
	StatsShuffles   int               // confirm / saved-stats budget
	MutationWorkers int               // 0 → DefaultWorkers()
	Seed            int64             // couples every screen shuffle and the confirm
	Progress        *AdaptiveProgress // optional live counters the ticker reads; nil to skip
	Cache           *Cache            // shared hand-eval cache; nil → a fresh one
}

// RunMutationRoundAdaptive runs one mutation round with SPRT screening + coupled confirm.
// cfg.MutationWorkers goroutines pull candidates from a shared queue; the first to clear both the
// screen and the confirm wins and cancels the rest.
//
// Each candidate is judged against a per-mutation accept threshold tau: a hill climb
// (cfg.Temperature == 0) sets tau = cfg.Threshold (the min-improvement gate); an SA step
// (cfg.Temperature > 0) sets tau = cfg.Temperature·ln(U) for a per-mutation uniform U, since the
// Metropolis rule accepts iff ΔV > cfg.Temperature·ln(U). Both are the same "is ΔV > tau?" test,
// screened by the SPRT and verified by the confirm. cfg.Seed couples every screen shuffle and the
// confirm across the incumbent and the mutants.
//
// Returns (acceptedDeck, acceptedStats, acceptedAvg, acceptedIndex, true) on the first confirmed
// acceptance, or (nil, zero, 0, -1, false) when every candidate is rejected or the round is
// cancelled.
func RunMutationRoundAdaptive(ctx context.Context, cfg AdaptiveRoundConfig) (*deck.Deck, deck.Stats, float64, int, bool) {
	if cfg.MutationWorkers <= 0 {
		cfg.MutationWorkers = defaultWorkers()
	}
	if len(cfg.Mutations) == 0 {
		return nil, deck.Stats{}, 0, -1, false
	}
	if cfg.Cache == nil {
		cfg.Cache = NewCache()
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	incVals := newIncumbentValues(cfg.Incumbent, NewEvaluatorWithCache(cfg.Cache), cfg.Matchup, cfg.Seed)
	conf := &confirmer{incumbent: cfg.Incumbent, mp: cfg.Matchup, shuffles: cfg.StatsShuffles, seed: perShuffleSeed(cfg.Seed, 0), cache: cfg.Cache}

	improvementCh := make(chan mutationImprovement, cfg.MutationWorkers)
	jobs := make(chan int, len(cfg.Mutations))
	for i := range cfg.Mutations {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < cfg.MutationWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runAdaptiveWorker(innerCtx, cancel, cfg, incVals, conf, jobs, improvementCh)
		}()
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
		select {
		case imp := <-improvementCh:
			return imp.candidate, imp.stats, imp.avg, imp.idx, true
		default:
		}
		return nil, deck.Stats{}, 0, -1, false
	}
}

// runAdaptiveWorker pulls candidates, screens each with the SPRT, and on a screen-accept runs the
// confirm. A confirmed improvement is sent and cancels the round; a screen false-accept that fails
// the confirm is dropped and the worker moves on.
func runAdaptiveWorker(ctx context.Context, cancel context.CancelFunc, cfg AdaptiveRoundConfig, incVals *incumbentValues, conf *confirmer, jobs <-chan int, improvementCh chan<- mutationImprovement) {
	ev := NewEvaluatorWithCache(cfg.Cache)
	rng := rand.New(rand.NewSource(1))  // reseeded per shuffle in value(); initial seed is irrelevant
	uRng := rand.New(rand.NewSource(1)) // reseeded per mutation for the SA acceptance draw
	for i := range jobs {
		if ctx.Err() != nil {
			return
		}
		// tau is the accept threshold for this candidate: min-improvement for a hill climb, the
		// random temperature·ln(U) Metropolis threshold for an SA step. U is drawn per mutation so
		// the SA acceptance is reproducible and independent of which worker pulls the candidate.
		tau := cfg.Threshold
		if cfg.Temperature > 0 {
			uRng.Seed(perShuffleSeed(cfg.Seed, -(i + 1)))
			tau = cfg.Temperature * math.Log(1-uRng.Float64())
		}
		d := cfg.Mutations[i].Deck()
		screen := newSingleShuffleEval(d, ev, cfg.Matchup)
		var acc sprtAccumulator
		verdict := sprtContinue
		for n := 1; verdict == sprtContinue; n++ {
			if ctx.Err() != nil {
				return
			}
			// Shift by (threshold - tau) so the SPRT's H0=0 / H1=threshold tests ΔV against
			// tau-threshold (reject) and tau (accept). At tau==threshold this is a no-op.
			delta := screen.value(rng, perShuffleSeed(cfg.Seed, n)) - incVals.at(n)
			acc.add(delta - tau + cfg.Threshold)
			verdict = acc.decision(cfg.Threshold, defaultSPRTConfig)
		}
		if cfg.Progress != nil {
			cfg.Progress.Screened.Add(1)
		}
		if verdict != sprtAccept {
			continue
		}
		// Screen says accept; confirm ΔV > tau at high precision, coupled to the incumbent baseline.
		// Confirming is the round's expensive moment — flag it so the ticker can report the full eval.
		if cfg.Progress != nil {
			cfg.Progress.Confirming.Add(1)
		}
		incAvg := conf.incumbentAvg()
		stats := ev.Evaluate(d, cfg.StatsShuffles, cfg.Matchup, rand.New(rand.NewSource(perShuffleSeed(cfg.Seed, 0))))
		if cfg.Progress != nil {
			cfg.Progress.Confirming.Add(-1)
		}
		mutAvg := stats.Mean()
		if mutAvg-incAvg <= tau {
			continue // screen false-accept — the confirm clears it, preserving termination
		}
		select {
		case improvementCh <- mutationImprovement{idx: i, avg: mutAvg, candidate: d, stats: stats}:
		default:
		}
		cancel()
		return
	}
}
