package sim

// Hand-by-hand simulation of a Deck: (*Evaluator).Evaluate shuffles, walks two cycles of
// hands per run, and folds each turn's outcome into a fresh deck.Stats. All cross-turn
// bookkeeping (held cards, arsenal, runechant carryover, start-of-turn Aura handling) lives
// here.

import (
	"io"
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Evaluate simulates `runs` shuffles of d. Each run assembles successive hands of
// d.Hero.Intelligence() cards (Held + fresh top-of-deck draws), computes the optimal play
// against mp, and recycles Pitched cards to deck bottom. A run ends when the deck can't
// fill the next hand. A "cycle" is one pass through the original deck size.
func (ev *Evaluator) Evaluate(d *deck.Deck, runs int, mp Matchup, rng *rand.Rand) deck.Stats {
	return ev.evaluateImpl(d, runs, mp, rng, nil)
}

// EvaluateAdaptive runs shuffles until the per-turn mean's standard error drops to
// precision/4 — a ~95% confidence interval of ±precision/2 around the running mean — capped
// at adaptiveShufflesCap. SE is checked every adaptiveCheckInterval shuffles. Use when
// knowing the mean to within precision is enough; modes that need apples-to-apples shuffle
// counts (compare, explicit -shuffles) should use Evaluate with a fixed runs count.
// Order-of-magnitude scale on a Viserai deck: precision=0.1 ≈ 1k shuffles, precision=0.01
// ≈ 80k shuffles.
func (ev *Evaluator) EvaluateAdaptive(d *deck.Deck, precision float64, mp Matchup, rng *rand.Rand) deck.Stats {
	return ev.evaluateImpl(d, adaptiveShufflesCap, mp, rng, makeAdaptiveStop(precision/4))
}

// shuffleStopper is the early-stop policy for the eval shuffle loop. Called once after each
// shuffle's stats are recorded; returning true breaks the loop. nil disables early stop.
type shuffleStopper func(stats *deck.Stats, runs int) bool

const (
	// adaptiveCheckInterval is the per-worker chunk size in the parallel-shuffle path: every
	// numWorkers × adaptiveCheckInterval shuffles, the pool barrier-merges into the run's
	// aggregate stats and runs the adaptive stop check. 50 is the empirical sweet spot —
	// smaller chunks let random Viserai decks converge in ~2000 shuffles; below 50 the
	// barrier overhead dominates.
	adaptiveCheckInterval = 50
	// adaptiveShufflesCap caps the adaptive shuffle path so a pathological high-variance
	// regime terminates. Sized for precision=0.01 (~82k shuffles on a baseline Viserai
	// deck) with headroom for higher-variance heroes.
	adaptiveShufflesCap = 200000
)

// makeAdaptiveStop returns a shuffleStopper that fires when the per-turn mean's standard
// error drops below targetSE. Checks every adaptiveCheckInterval shuffles to amortise the
// histogram walk.
func makeAdaptiveStop(targetSE float64) shuffleStopper {
	return func(stats *deck.Stats, runs int) bool {
		if runs%adaptiveCheckInterval != 0 {
			return false
		}
		return meanStandardError(stats) <= targetSE
	}
}

// meanStandardError computes the standard error of the per-turn mean Value: sigma / sqrt(N)
// where sigma is the unbiased per-turn sample standard deviation. Walks the histogram so
// it's O(unique values) ~ O(30) per call rather than O(N). Returns +Inf when fewer than two
// turns have been simulated (variance is undefined).
func meanStandardError(stats *deck.Stats) float64 {
	n := float64(stats.Hands)
	if n < 2 {
		return math.Inf(1)
	}
	mean := stats.TotalValue / n
	sumSq := 0.0
	for v, count := range stats.Histogram {
		diff := float64(v) - mean
		sumSq += diff * diff * float64(count)
	}
	variance := sumSq / (n - 1)
	return math.Sqrt(variance / n)
}

func (ev *Evaluator) evaluateImpl(d *deck.Deck, maxRuns int, mp Matchup, rng *rand.Rand, stop shuffleStopper) deck.Stats {
	handSize := d.Hero.(hero.Hero).Intelligence()
	deckSize := d.Size()
	if handSize <= 0 || deckSize < handSize {
		return deck.Stats{}
	}
	if ev.numWorkers > 1 {
		return ev.evaluateParallelImpl(d, maxRuns, mp, rng, stop, handSize, deckSize)
	}
	return ev.evaluateSequentialImpl(d, maxRuns, mp, rng, stop, handSize, deckSize)
}

// evaluateSequentialImpl runs the shuffle loop in the calling goroutine — the
// deterministic-RNG path tests rely on.
func (ev *Evaluator) evaluateSequentialImpl(d *deck.Deck, maxRuns int, mp Matchup, rng *rand.Rand, stop shuffleStopper, handSize, deckSize int) deck.Stats {
	handsPerCycle := deckSize / handSize
	uniqueIDs, idIndex := d.UniqueIDs()
	scratch := newShuffleScratch(len(d.Weapons), deckSize, handSize, len(uniqueIDs))

	var stats deck.Stats
	actualRuns := 0
	for r := 0; r < maxRuns; r++ {
		runOneShuffle(d, scratch, &stats, idIndex, ev, rng, mp, handsPerCycle, handSize)
		actualRuns = r + 1
		if stop != nil && stop(&stats, actualRuns) {
			break
		}
	}
	stats.Runs += actualRuns
	mergeMarginalBuf(&stats, uniqueIDs, scratch.marginalBuf)
	return stats
}

// evaluateParallelImpl fans the shuffle loop across ev.numWorkers goroutines that share
// ev.cache (RWMutex-protected) but each carry their own per-call scratch. Shuffles process
// in chunks of (numWorkers × adaptiveCheckInterval); after each chunk the main goroutine
// merges every worker's local deck.Stats into the running aggregate and runs the adaptive
// stop check. Per-worker seeds derive from rng.Int63() so the distribution is deterministic
// given the input rng.
func (ev *Evaluator) evaluateParallelImpl(d *deck.Deck, maxRuns int, mp Matchup, rng *rand.Rand, stop shuffleStopper, handSize, deckSize int) deck.Stats {
	numWorkers := ev.numWorkers
	handsPerCycle := deckSize / handSize
	uniqueIDs, idIndex := d.UniqueIDs()
	aggregateMarginal := make([]deck.CardMarginalStats, len(uniqueIDs))

	chunkPerWorker := adaptiveCheckInterval
	maxChunk := numWorkers * chunkPerWorker

	type partial struct {
		stats    deck.Stats
		marginal []deck.CardMarginalStats
	}
	results := make(chan partial, numWorkers)

	var stats deck.Stats
	actualRuns := 0
	for actualRuns < maxRuns {
		sz := maxChunk
		if actualRuns+sz > maxRuns {
			sz = maxRuns - actualRuns
		}
		runsPerWorker := sz / numWorkers
		extras := sz % numWorkers

		spawned := 0
		var wg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			myRuns := runsPerWorker
			if w < extras {
				myRuns++
			}
			if myRuns == 0 {
				continue
			}
			mySeed := rng.Int63()
			spawned++
			wg.Add(1)
			go func(seed int64, runs int) {
				defer wg.Done()
				workerEv := &Evaluator{cache: ev.cache}
				workerRNG := rand.New(rand.NewSource(seed))
				scratch := newShuffleScratch(len(d.Weapons), deckSize, handSize, len(uniqueIDs))
				var local deck.Stats
				for r := 0; r < runs; r++ {
					runOneShuffle(d, scratch, &local, idIndex, workerEv, workerRNG, mp, handsPerCycle, handSize)
				}
				results <- partial{stats: local, marginal: scratch.marginalBuf}
			}(mySeed, myRuns)
		}
		wg.Wait()
		for i := 0; i < spawned; i++ {
			r := <-results
			mergeStatsInto(&stats, &r.stats)
			for j := range aggregateMarginal {
				aggregateMarginal[j].PresentTotal += r.marginal[j].PresentTotal
				aggregateMarginal[j].PresentHands += r.marginal[j].PresentHands
				aggregateMarginal[j].AbsentTotal += r.marginal[j].AbsentTotal
				aggregateMarginal[j].AbsentHands += r.marginal[j].AbsentHands
			}
		}
		actualRuns += sz
		if stop != nil && stop(&stats, actualRuns) {
			break
		}
	}
	stats.Runs += actualRuns
	mergeMarginalBuf(&stats, uniqueIDs, aggregateMarginal)
	return stats
}

// shuffleScratch holds the per-goroutine slabs the shuffle loop reuses across iterations.
// Cross-turn aura / item / banished / graveyard / opponentMarked state lives on the master
// *gameengine.GameState threaded across turns; only per-turn buffers sit here.
type shuffleScratch struct {
	weaponsBuf  []weapon.Weapon
	handBuf     []card.Card
	heldBuf     []card.Card
	presentBuf  []bool
	marginalBuf []deck.CardMarginalStats
	// recycledBuf backs the per-turn "pitched-to-deck-bottom" slice runOneShuffle hands
	// to d.PutBottom. Refilled each turn; PutBottom copies into its own backing so this
	// stays safe to overwrite next turn.
	recycledBuf []deck.Card
}

// newShuffleScratch sizes the per-shuffle reusable buffers for a deck of (weaponCount,
// deckSize, handSize) shape. Called once per worker; the returned scratch is reused across
// every shuffle that worker runs.
func newShuffleScratch(weaponCount, _, handSize, numUniqueIDs int) *shuffleScratch {
	return &shuffleScratch{
		weaponsBuf:  make([]weapon.Weapon, weaponCount),
		handBuf:     make([]card.Card, handSize, handSize+startOfTurnRevealRoom),
		heldBuf:     make([]card.Card, 0, handSize),
		presentBuf:  make([]bool, numUniqueIDs),
		marginalBuf: make([]deck.CardMarginalStats, numUniqueIDs),
	}
}

// runOneShuffle simulates one shuffle end-to-end (Copy → Shuffle → walk turns → record
// stats), accumulating into the caller-owned *deck.Stats. masterDeck is copied before
// shuffling so each trial gets an independent deck and the master's Cards order stays
// stable across goroutines.
func runOneShuffle(masterDeck *deck.Deck, scratch *shuffleScratch, stats *deck.Stats, idIndex map[ids.CardID]int, ev *Evaluator, rng *rand.Rand, mp Matchup, handsPerCycle, handSize int) {
	d := masterDeck.Copy()
	d.Shuffle(rng)

	heroVal := d.Hero.(hero.Hero)
	weapons := scratch.weaponsBuf
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}

	// master is the start-of-turn carryover state, built once per shuffle and threaded
	// across turns. play.State (the post-chain winner) becomes the next turn's master.
	master := gameengine.GameStateBuilder().
		SetHero(heroVal).
		SetIncomingDamage(mp.IncomingDamage).
		SetArcaneIncomingDamage(mp.ArcaneIncomingDamage).
		Build()

	handIdx := 0
	heldBuf := scratch.heldBuf[:0]
	handBuf := scratch.handBuf
	maxHands := 2 * handsPerCycle
	for handIdx < maxHands {
		drawCount := handSize - len(heldBuf)
		if drawCount <= 0 || d.Size() < drawCount {
			break
		}
		// Build hand: held prefix plus freshly drawn cards from the deck top.
		h := handBuf[:handSize]
		copy(h, heldBuf)
		drawn := d.Draw(drawCount)
		for i, c := range drawn {
			h[len(heldBuf)+i] = c.(card.Card)
		}
		// Reveals from start-of-turn auras (e.g. Sigil of the Arknight) pop d in place and
		// append onto h before the chain runs.
		trigDamage := processAurasAtStartOfTurn(master, d, &h)
		arsenalIn := master.Arsenal()
		sortHandByID(h)
		play := runBestForTurn(weapons, h, mp, d, master, ev)
		play.Value += trigDamage

		if recordTurnStats(stats, play, handIdx, handsPerCycle) {
			recordBestTurn(stats, play, ev, weapons, h, mp, d, master)
		}
		tallyMarginalPresence(scratch.marginalBuf, idIndex, scratch.presentBuf, h, arsenalIn, float64(play.Value))
		// Adopt the chain's post-mutation deck and recycle pitched cards onto the
		// bottom — FaB's end-of-turn pitch-zone-to-deck rule.
		d = play.State.Deck()
		pitched := pitchedFromBestLine(play.BestLine)
		recycled := scratch.recycledBuf[:0]
		for _, c := range pitched {
			recycled = append(recycled, c)
		}
		scratch.recycledBuf = recycled
		d.PutBottom(recycled)
		// Carry hand leftover into next turn's heldBuf; thread play.State forward as master.
		heldBuf = append(heldBuf[:0], play.State.Hand()...)
		master = play.State
		master.ResetEphemeralState()
		// Clear master.Deck — d already holds the post-turn deck and the chain runner will
		// install its own per-leaf copy, so leaving one on the master is dead weight.
		master.SetDeck(nil)
		handIdx++
	}
	scratch.heldBuf = heldBuf
}

// mergeStatsInto folds src's per-shuffle accumulators into dst. PerCardMarginal merging
// happens separately via mergeMarginalBuf at the end of the run.
func mergeStatsInto(dst, src *deck.Stats) {
	dst.Hands += src.Hands
	dst.TotalValue += src.TotalValue
	dst.FirstCycle.Hands += src.FirstCycle.Hands
	dst.FirstCycle.Total += src.FirstCycle.Total
	dst.SecondCycle.Hands += src.SecondCycle.Hands
	dst.SecondCycle.Total += src.SecondCycle.Total
	if dst.Histogram == nil && len(src.Histogram) > 0 {
		dst.Histogram = make(map[int]int, len(src.Histogram))
	}
	for v, c := range src.Histogram {
		dst.Histogram[v] += c
	}
	if len(src.Best.BestLine) > 0 &&
		(len(dst.Best.BestLine) == 0 || src.Best.Value > dst.Best.Value) {
		dst.Best = src.Best
		dst.PrintBest = src.PrintBest
	}
}

// runBestForTurn dispatches to ev.BestSkipLog — the hot-path case that skips log
// materialisation.
func runBestForTurn(
	weapons []weapon.Weapon,
	h []card.Card,
	mp Matchup,
	d *deck.Deck,
	master *gameengine.GameState,
	ev *Evaluator,
) TurnSummary {
	return ev.BestSkipLog(weapons, h, mp, d, master)
}

// recordTurnStats folds one resolved turn's accumulators into stats: bumps Hands /
// TotalValue, lazily initialises Histogram, and credits the value to FirstCycle /
// SecondCycle based on handIdx. Returns true when this turn's Value beats stats.Best.
func recordTurnStats(stats *deck.Stats, play TurnSummary, handIdx, handsPerCycle int) bool {
	v := float64(play.Value)
	stats.TotalValue += v
	stats.Hands++
	if stats.Histogram == nil {
		stats.Histogram = map[int]int{}
	}
	stats.Histogram[play.Value]++
	newBest := play.Value > stats.Best.Value || len(stats.Best.BestLine) == 0
	switch handIdx / handsPerCycle {
	case 0:
		stats.FirstCycle.Hands++
		stats.FirstCycle.Total += v
	case 1:
		stats.SecondCycle.Hands++
		stats.SecondCycle.Total += v
	}
	return newBest
}

// startOfTurnRevealRoom caps how many cards a start-of-turn Aura reveal can append to a
// turn's dealt hand. Sized large enough that handBuf never reallocates.
const startOfTurnRevealRoom = 8

// processAurasAtStartOfTurn fires every triggertype.StartOfTurn handler queued on master,
// returns the summed damage to fold into the turn's Value, and appends any cards the
// handlers revealed onto h. Re-arms FiredThisTurn at the same time.
//
// Mutates master in place: destroyed auras splice out, revealed cards push onto
// master.Hand() (then move into h). master.Deck is set to d for the duration so reveal
// handlers can PopDeckTop, then cleared on return. Cascading reveals: a handler that pops
// d shrinks the view for the next handler, so two reveal-capable auras see distinct tops.
func processAurasAtStartOfTurn(master *gameengine.GameState, d *deck.Deck, h *[]card.Card) int {
	if len(master.Auras()) == 0 {
		return 0
	}
	master.SetDeck(d)
	preHand := len(master.Hand())
	damage := master.Engine().FireStartOfTurn()
	if revealed := master.Hand(); len(revealed) > preHand {
		*h = append(*h, revealed[preHand:]...)
		master.SetHand(revealed[:preHand])
	}
	master.SetDeck(nil)
	return damage
}

// pitchedFromBestLine returns the cards in BestLine assigned the Pitch role (excluding the
// arsenal-in slot, which never recycles). Sorted by ID so the deck-bottom recycle order is
// canonical by multiset — cache-replay and from-scratch search may pick different tie-break
// winners among equally-optimal partitions, and canonical sorting keeps their recycled
// decks byte-identical.
func pitchedFromBestLine(line []deck.CardAssignment) []card.Card {
	var out []card.Card
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		if a.Role == deck.Pitch {
			out = append(out, a.Card)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// sortHandByID sorts hand in place by Card.ID() so the partition recurse enumerates against
// a canonical order — drops a positional-tie-break source from findBest's leaf comparator
// and makes the cache-off / cache-on paths produce byte-identical results for matching
// multisets. Insertion sort beats sort.SliceStable handily at the small hand sizes (~7)
// where the reflection-based interface comparator dominates wall-clock.
func sortHandByID(hand []card.Card) {
	for i := 1; i < len(hand); i++ {
		c := hand[i]
		id := c.ID()
		j := i - 1
		for j >= 0 && hand[j].ID() > id {
			hand[j+1] = hand[j]
			j--
		}
		hand[j+1] = c
	}
}

// recordBestTurn clones the winning BestLine into stats.Best and attaches stats.PrintBest
// with a snapshot of the carryover state / deck / hand / matchup that the print path can
// replay on demand.
func recordBestTurn(stats *deck.Stats, play TurnSummary, ev *Evaluator, weapons []weapon.Weapon, h []card.Card, mp Matchup, d *deck.Deck, master *gameengine.GameState) {
	lineCopy := make([]deck.CardAssignment, len(play.BestLine))
	copy(lineCopy, play.BestLine)
	stats.Best = deck.BestTurn{
		Value:    play.Value,
		BestLine: lineCopy,
	}
	weaponsCopy := append([]weapon.Weapon(nil), weapons...)
	handCopy := append([]card.Card(nil), h...)
	snap := &bestSnapshot{
		master:   master.CopyPersistentState(),
		deck:     d.Copy(),
		hand:     handCopy,
		weapons:  weaponsCopy,
		mp:       mp,
		bestLine: lineCopy,
		value:    play.Value,
	}
	stats.PrintBest = func(w io.Writer) { PrintBestTurn(ev, snap, w) }
}

// tallyMarginalPresence credits this turn's value to each entry in marginalBuf, bucketed by
// whether the card was present in the dealt hand or in the arsenal-in slot when Best ran.
// presentBuf is a scratch slice indexed parallel to marginalBuf; both are caller-owned
// across turns to keep this path allocation-free.
func tallyMarginalPresence(marginalBuf []deck.CardMarginalStats, idIndex map[ids.CardID]int, presentBuf []bool, dealt []card.Card, arsenalIn card.Card, value float64) {
	if len(marginalBuf) == 0 {
		return
	}
	clear(presentBuf)
	for _, c := range dealt {
		if i, ok := idIndex[c.ID()]; ok {
			presentBuf[i] = true
		}
	}
	if arsenalIn != nil {
		if i, ok := idIndex[arsenalIn.ID()]; ok {
			presentBuf[i] = true
		}
	}
	for i := range marginalBuf {
		if presentBuf[i] {
			marginalBuf[i].PresentTotal += value
			marginalBuf[i].PresentHands++
		} else {
			marginalBuf[i].AbsentTotal += value
			marginalBuf[i].AbsentHands++
		}
	}
}

// mergeMarginalBuf folds the per-Evaluate slice accumulator into PerCardMarginal. The map
// is lazily initialised so unscored decks don't pay for an empty map.
func mergeMarginalBuf(stats *deck.Stats, uniqueIDs []ids.CardID, marginalBuf []deck.CardMarginalStats) {
	if len(uniqueIDs) == 0 {
		return
	}
	if stats.PerCardMarginal == nil {
		stats.PerCardMarginal = make(map[ids.CardID]deck.CardMarginalStats, len(uniqueIDs))
	}
	for i, id := range uniqueIDs {
		m := stats.PerCardMarginal[id]
		m.PresentTotal += marginalBuf[i].PresentTotal
		m.PresentHands += marginalBuf[i].PresentHands
		m.AbsentTotal += marginalBuf[i].AbsentTotal
		m.AbsentHands += marginalBuf[i].AbsentHands
		stats.PerCardMarginal[id] = m
	}
}
