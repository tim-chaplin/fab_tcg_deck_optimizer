package sim

// Adaptive mutation-round runner for the T==0 hill climb. Each candidate is screened by an SPRT
// on coupled per-shuffle value deltas (see sprt.go) instead of a fixed shuffle budget: clear
// non-improvements are dropped in a handful of shuffles, only close calls run long, and there is
// no max-shuffle cap. A mutation the screen accepts is then confirmed with a high-precision
// coupled eval at statsShuffles before it can win the round — so a screen false-accept can't be
// taken as a real improvement, which is what guarantees the round terminates.
//
// Parallelism here is across mutations (one worker per goroutine), not within an eval: a per-
// shuffle sequential screen can't fan a single shuffle across workers. The incumbent's per-shuffle
// values are shared across workers through incumbentValues.

import (
	"context"
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

// RunMutationRoundAdaptive runs one T==0 hill-climb round with SPRT screening + coupled confirm.
// mutationWorkers goroutines pull candidates from a shared queue; the first to clear both the
// screen and the confirm wins and cancels the rest. threshold is the min-improvement gate, applied
// at both the screen (as the SPRT's H1) and the confirm. statsShuffles is the confirm / saved-stats
// budget. seed couples every screen shuffle and the confirm across the incumbent and the mutants.
//
// Returns (acceptedDeck, acceptedStats, acceptedAvg, acceptedIndex, true) on the first confirmed
// improvement, or (nil, zero, 0, -1, false) when every candidate is rejected or the round is
// cancelled.
func RunMutationRoundAdaptive(
	ctx context.Context,
	mutations []deck.Mutation,
	incumbent *deck.Deck,
	threshold float64,
	mp Matchup,
	statsShuffles int,
	mutationWorkers int,
	seed int64,
	completed *atomic.Int64,
	cache *Cache,
) (*deck.Deck, deck.Stats, float64, int, bool) {
	if mutationWorkers <= 0 {
		mutationWorkers = defaultWorkers()
	}
	if len(mutations) == 0 {
		return nil, deck.Stats{}, 0, -1, false
	}
	if cache == nil {
		cache = NewCache()
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	incVals := newIncumbentValues(incumbent, NewEvaluatorWithCache(cache), mp, seed)
	conf := &confirmer{incumbent: incumbent, mp: mp, shuffles: statsShuffles, seed: perShuffleSeed(seed, 0), cache: cache}

	improvementCh := make(chan mutationImprovement, mutationWorkers)
	jobs := make(chan int, len(mutations))
	for i := range mutations {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < mutationWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runAdaptiveWorker(innerCtx, cancel, mutations, incVals, conf, threshold, mp, statsShuffles, seed, completed, cache, jobs, improvementCh)
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
func runAdaptiveWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	mutations []deck.Mutation,
	incVals *incumbentValues,
	conf *confirmer,
	threshold float64,
	mp Matchup,
	statsShuffles int,
	seed int64,
	completed *atomic.Int64,
	cache *Cache,
	jobs <-chan int,
	improvementCh chan<- mutationImprovement,
) {
	ev := NewEvaluatorWithCache(cache)
	rng := rand.New(rand.NewSource(1)) // reseeded per shuffle in value(); initial seed is irrelevant
	for i := range jobs {
		if ctx.Err() != nil {
			return
		}
		d := mutations[i].Deck()
		screen := newSingleShuffleEval(d, ev, mp)
		var acc sprtAccumulator
		verdict := sprtContinue
		for n := 1; verdict == sprtContinue; n++ {
			if ctx.Err() != nil {
				return
			}
			delta := screen.value(rng, perShuffleSeed(seed, n)) - incVals.at(n)
			acc.add(delta)
			verdict = acc.decision(threshold, defaultSPRTConfig)
		}
		if completed != nil {
			completed.Add(1)
		}
		if verdict != sprtAccept {
			continue
		}
		// Screen says improvement; confirm at high precision, coupled to the incumbent baseline.
		stats := ev.Evaluate(d, statsShuffles, mp, rand.New(rand.NewSource(perShuffleSeed(seed, 0))))
		mutAvg := stats.Mean()
		if mutAvg-conf.incumbentAvg() <= threshold {
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
