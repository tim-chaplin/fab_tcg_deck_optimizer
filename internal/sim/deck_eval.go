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
// Cross-turn carryover lives on the master *gameengine.GameState; only per-turn buffers
// sit here.
type shuffleScratch struct {
	weaponsBuf  []weapon.Weapon
	handBuf     []card.Card
	presentBuf  []bool
	marginalBuf []deck.CardMarginalStats
	// recycledBuf backs the per-turn pitched-to-deck-bottom slice fed to d.PutBottom.
	// PutBottom copies into its own backing, so this stays safe to overwrite next turn.
	recycledBuf []deck.Card
}

// newShuffleScratch sizes the per-shuffle reusable buffers for a deck of (weaponCount,
// deckSize, handSize) shape. Called once per worker; the returned scratch is reused across
// every shuffle that worker runs.
func newShuffleScratch(weaponCount, _, handSize, numUniqueIDs int) *shuffleScratch {
	return &shuffleScratch{
		weaponsBuf:  make([]weapon.Weapon, weaponCount),
		handBuf:     make([]card.Card, handSize, handSize+startOfTurnRevealRoom),
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

	master := gameengine.GameStateBuilder().
		SetHero(heroVal).
		SetIncomingDamage(mp.IncomingDamage).
		SetArcaneIncomingDamage(mp.ArcaneIncomingDamage).
		Build()

	// Initial hand drawn into the reusable handBuf.
	handBuf := scratch.handBuf
	h := handBuf[:handSize]
	for i, c := range d.Draw(handSize) {
		h[i] = c.(card.Card)
	}

	// recordIfBest records-if-best and tallies marginal presence. handIdx is captured by
	// reference so each call sees the current iteration's value.
	handIdx := 0
	recordIfBest := func(summary TurnSummary, dealtHand []card.Card, arsenalIn card.Card, snap *bestSnapshot) {
		if recordTurnStats(stats, summary, handIdx, handsPerCycle) {
			recordBestTurnFromSnap(stats, summary, ev, snap)
		}
		tallyMarginalPresence(scratch.marginalBuf, idIndex, scratch.presentBuf, dealtHand, arsenalIn, float64(summary.Value))
	}

	maxHands := 2 * handsPerCycle
	for ; handIdx < maxHands; handIdx++ {
		preDeckSize := d.Size()
		var summary TurnSummary
		summary, scratch.recycledBuf = playOneTurn(master, h, d, mp, weapons, ev, handSize, scratch.recycledBuf, recordIfBest, nil)

		master = summary.State
		d = summary.State.Deck()
		h = summary.State.Hand()
		// Stop when no fresh cards entered the next hand: either deck exhausted (len(h) <
		// handSize) or chain held everything (toDraw == 0, which would loop on identical
		// inputs forever). freshDrawn = preDeck + pitched - postDeck reduces to draw count.
		freshDrawn := preDeckSize + len(scratch.recycledBuf) - d.Size()
		if freshDrawn == 0 || len(h) < handSize {
			break
		}
	}
}

// playOneTurn drives one full turn: advance master, snapshot for replay (if record != nil),
// fire start-of-turn auras, run chain via Best, invoke record, recycle pitched to deck
// bottom, draw the next hand (partial OK).
//
// Returned summary.State is the end-of-turn boundary (pitched recycled, next hand drawn,
// Value accrued) and threads directly into the next call as master. recycledBuf is reused
// in place when non-nil.
//
// When snapshot is non-nil, runs in replay mode: drives the chain through snapshot.bestLine
// + snapshot.cardsPlayed (no enumeration), streams emissions via snapshot.logger, and
// returns the raw post-chain per-perm state without recycle / next-draw cleanup.
func playOneTurn(
	master *gameengine.GameState,
	hand []card.Card,
	d *deck.Deck,
	mp Matchup,
	weapons []weapon.Weapon,
	ev *Evaluator,
	handSize int,
	recycledBuf []deck.Card,
	record turnRecord,
	snapshot *bestSnapshot,
) (summary TurnSummary, recycledOut []deck.Card) {
	advanceToNextTurn(master)

	var snap *bestSnapshot
	var arsenalIn card.Card
	if snapshot == nil && record != nil {
		arsenalIn = master.Arsenal()
		snap = &bestSnapshot{
			master:  master.CopyPersistentState(),
			deck:    d.Copy(),
			hand:    append([]card.Card(nil), hand...),
			weapons: append([]weapon.Weapon(nil), weapons...),
			mp:      mp,
		}
	}

	processAurasAtStartOfTurn(master, d, &hand)
	sortHandByID(hand)
	dealtHand := hand
	if snapshot != nil {
		summary = runReplayForTurn(snapshot)
	} else {
		summary = runBestForTurn(weapons, hand, mp, d, master, ev)
	}

	if record != nil {
		record(summary, dealtHand, arsenalIn, snap)
	}

	if snapshot != nil {
		// Skip end-of-turn cleanup so summary.State stays at the post-chain per-perm state
		// the caller needs for its end-of-turn snapshot.
		return summary, recycledBuf
	}

	// Chain ran on a shallow copy of d that may have drawn mid-turn; use the winner's
	// post-chain deck for recycle / next-turn draw.
	postChainDeck := summary.State.Deck()
	recycledOut = recycledBuf[:0]
	for _, c := range pitchedFromBestLine(summary.BestLine) {
		recycledOut = append(recycledOut, c)
	}
	postChainDeck.PutBottom(recycledOut)

	held := summary.State.Hand()
	toDraw := handSize - len(held)
	if toDraw > postChainDeck.Size() {
		toDraw = postChainDeck.Size()
	}
	nextHand := held
	if toDraw > 0 {
		// Allocate fresh: held aliases the chain's per-perm hand buffer, which the next
		// playOneTurn call may overwrite.
		nextHand = make([]card.Card, len(held), len(held)+toDraw)
		copy(nextHand, held)
		for _, c := range postChainDeck.Draw(toDraw) {
			nextHand = append(nextHand, c.(card.Card))
		}
	}

	summary.State.SetHand(nextHand)
	return summary, recycledOut
}

// turnRecord runs after the chain. snap is captured at start-of-turn (post-reset,
// pre-aura) so PrintBestTurn can replay against the exact deck the chain saw — critical
// when chain effects draw or tutor. Tests pass nil to skip the snapshot alloc.
type turnRecord func(summary TurnSummary, dealtHand []card.Card, arsenalIn card.Card, snap *bestSnapshot)

// advanceToNextTurn clears per-turn ephemerals (value, cardsPlayed, ...) and detaches any
// stale deck pointer. Idempotent on a freshly-built master.
func advanceToNextTurn(master *gameengine.GameState) {
	master.ResetEphemeralState()
	master.SetDeck(nil)
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

// runBestForTurn dispatches to ev.Best.
func runBestForTurn(
	weapons []weapon.Weapon,
	h []card.Card,
	mp Matchup,
	d *deck.Deck,
	master *gameengine.GameState,
	ev *Evaluator,
) TurnSummary {
	return ev.Best(weapons, h, mp, d, master)
}

// recordTurnStats folds one resolved turn's accumulators into stats: bumps Hands /
// TotalValue, lazily initialises Histogram, and credits the value to FirstCycle /
// SecondCycle based on handIdx. Returns true when this turn's Value beats stats.Best.
func recordTurnStats(stats *deck.Stats, summary TurnSummary, handIdx, handsPerCycle int) bool {
	v := float64(summary.Value)
	stats.TotalValue += v
	stats.Hands++
	if stats.Histogram == nil {
		stats.Histogram = map[int]int{}
	}
	stats.Histogram[summary.Value]++
	newBest := summary.Value > stats.Best.Value || len(stats.Best.BestLine) == 0
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

// processAurasAtStartOfTurn fires every StartOfTurn aura handler queued on master. Handlers
// write value gains directly to master.Value() and reveals append to h. Re-arms
// FiredThisTurn. Callers must refill h to full size first so reveal handlers see the
// post-draw deck top. Cascading reveals: a handler that pops d shrinks the view for the
// next, so two reveal-capable auras see distinct tops.
func processAurasAtStartOfTurn(master *gameengine.GameState, d *deck.Deck, h *[]card.Card) {
	if len(master.Auras()) == 0 {
		return
	}
	master.SetDeck(d)
	preHand := len(master.Hand())
	master.Engine().FireStartOfTurn()
	if revealed := master.Hand(); len(revealed) > preHand {
		*h = append(*h, revealed[preHand:]...)
		master.SetHand(revealed[:preHand])
	}
	master.SetDeck(nil)
}

// pitchedFromBestLine returns BestLine's Pitch-role cards (excluding the arsenal-in slot,
// which never recycles), sorted by ID so recycled decks stay byte-identical across cache
// hits and from-scratch searches that pick different equally-optimal partitions.
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

// sortHandByID sorts hand in place by Card.ID() to give findBest a canonical enumeration
// order, so cache-off and cache-on paths produce byte-identical results for matching
// multisets. Insertion sort wins at small hand sizes (~7) where the reflection-based
// interface comparator dominates wall-clock.
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

// recordBestTurnFromSnap commits a winning turn to stats: clones BestLine into stats.Best,
// fills snap with the chain-produced bestLine / cardsPlayed / swungWeapons / value, and
// attaches stats.PrintBest as the deferred replay closure.
func recordBestTurnFromSnap(stats *deck.Stats, summary TurnSummary, ev *Evaluator, snap *bestSnapshot) {
	lineCopy := make([]deck.CardAssignment, len(summary.BestLine))
	copy(lineCopy, summary.BestLine)
	stats.Best = deck.BestTurn{
		Value:    summary.Value,
		BestLine: lineCopy,
	}
	snap.bestLine = lineCopy
	snap.swungWeapons = append([]string(nil), summary.SwungWeapons...)
	if summary.State != nil {
		snap.cardsPlayed = append([]card.Card(nil), summary.State.CardsPlayed()...)
	}
	snap.value = summary.Value
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
