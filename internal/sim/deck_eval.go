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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
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
				workerEv := &Evaluator{cache: ev.cache, statePool: gameengine.NewPrewarmedPool()}
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
	handBuf     []card.Card
	presentBuf  []bool
	marginalBuf []deck.CardMarginalStats
}

// newShuffleScratch sizes the per-shuffle reusable buffers for a deck of (deckSize, handSize)
// shape. Called once per worker; the returned scratch is reused across every shuffle that
// worker runs.
func newShuffleScratch(_, _, handSize, numUniqueIDs int) *shuffleScratch {
	return &shuffleScratch{
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

	// Carry state borrows one pool slot; playOneTurn / Best mutate it in place across
	// turns. Put back at shuffle end before FreeAll.
	state := ev.statePool.Get()
	state.Reset(d.Hero.(hero.Hero), weaponsFromDeck(d), mp.IncomingPhysicalDamage, mp.IncomingArcaneDamage)

	// Initial hand drawn into the reusable handBuf, sorted so it is canonical from turn one.
	handBuf := scratch.handBuf
	h := handBuf[:handSize]
	for i, c := range d.Draw(handSize) {
		h[i] = c.(card.Card)
	}
	sortHandByID(h)

	maxHands := 2 * handsPerCycle
	for handIdx := 0; handIdx < maxHands; handIdx++ {
		preDeckSize := d.Size()
		state.SetHand(h)
		summary, snap := playOneTurn(state, d, ev, nil, nil)

		if recordTurnStats(stats, summary, handIdx, handsPerCycle) {
			recordBestTurnFromSnap(stats, summary, ev, snap)
		}
		tallyMarginalPresence(scratch.marginalBuf, idIndex, scratch.presentBuf, summary.BestLine, float64(summary.Value))

		state = summary.State
		d = summary.State.Deck()
		h = summary.State.Hand()
		// Stop when no fresh cards entered the next hand: either deck exhausted (len(h) <
		// handSize) or attack turn held everything (toDraw == 0, which would loop on identical
		// inputs forever). freshDrawn = preDeck + pitched - postDeck reduces to draw count.
		freshDrawn := preDeckSize + countPitched(summary.BestLine) - d.Size()
		if freshDrawn == 0 || len(h) < handSize {
			break
		}
	}
	ev.statePool.Put(state)
	ev.statePool.FreeAll()
}

// weaponsFromDeck equips the deck's weapons: each d.Weapons entry is a platonic weapon
// card (weapon.Card), so this builds the mutable engine-side weapon object for each — the
// game-start equip step. Fresh objects per call keep one shuffle's per-turn counter state
// from leaking into the next.
func weaponsFromDeck(d *deck.Deck) []gameengine.Weapon {
	cards := make([]weapon.Card, len(d.Weapons))
	for i, w := range d.Weapons {
		cards[i] = w.(weapon.Card)
	}
	return gameengine.EquipFromCards(cards)
}

// countPitched returns the number of Pitch-role entries in bestLine, excluding the arsenal
// slot (which never recycles).
func countPitched(bestLine []card.CardAssignment) int {
	n := 0
	for _, a := range bestLine {
		if a.Role == card.Pitch && !a.FromArsenal {
			n++
		}
	}
	return n
}

// playOneTurn drives one full turn: advance state, capture the start-of-turn snapshot,
// fire start-of-turn auras, run attack turn via Best, recycle pitched to deck bottom, draw the
// next hand (partial OK).
//
// Returned summary.State is the end-of-turn boundary (pitched recycled, next hand drawn,
// Value accrued) and threads directly into the next call as the carryover state. snap is
// the start-of-turn snapshot (independent of subsequent mutations to state/d) for
// record-if-best / replay; nil in replay mode.
//
// When snapshot is non-nil, runs in REPLAY mode: drives the attack turn through snapshot.bestLine
// + snapshot.cardsPlayed (no enumeration), streams emissions via logger (when non-nil), and
// returns the raw post-attack-turn per-perm state without recycle / next-draw cleanup.
func playOneTurn(
	state *gameengine.GameState,
	d *deck.Deck,
	ev *Evaluator,
	snapshot *turnSnapshot,
	logger card.Logger,
) (summary TurnSummary, snap *turnSnapshot) {
	advanceToNextTurn(state)

	if snapshot == nil {
		snap = &turnSnapshot{
			state: state.CopyPersistentState(),
			deck:  d.Copy(),
			hand:  append([]card.Card(nil), state.Hand()...),
		}
	}

	// Hand stays sorted by Card.ID() through the aura handlers, so the attack-turn runner and
	// cache key see a canonical multiset.
	processAurasAtStartOfTurn(state, d)
	if snapshot != nil {
		summary = runReplayForTurn(snapshot, ev, logger)
		// Skip end-of-turn cleanup so summary.State stays at the post-attack-turn per-perm state
		// the caller needs for its end-of-turn snapshot.
		return summary, nil
	}
	preAttackTurnHand := state.Hand()
	summary = runBestForTurn(state.Weapons(), preAttackTurnHand, d, state, ev)

	postAttackTurnDeck := recyclePitchedToDeckBottom(summary)
	intellect := state.Hero().(hero.Hero).Intelligence()
	nextHand := drawNextHand(postAttackTurnDeck, summary.State.Hand(), intellect)

	summary.State.SetHand(nextHand)
	verifyTurnInvariants(snap, preAttackTurnHand, summary)
	return summary, snap
}

// recyclePitchedToDeckBottom pushes the turn's pitched cards (excluding the arsenal-in
// slot, which never recycles) onto the bottom of the post-attack-turn deck and returns that
// deck. The attack turn ran on a shallow copy of the master deck that may have drawn mid-turn,
// so summary.State.Deck() — not the original d — is the deck the next turn inherits.
func recyclePitchedToDeckBottom(summary TurnSummary) *deck.Deck {
	postAttackTurnDeck := summary.State.Deck()
	pitched := pitchedFromBestLine(summary.BestLine)
	recycled := make([]deck.Card, len(pitched))
	for i, c := range pitched {
		recycled[i] = c
	}
	postAttackTurnDeck.PutBottom(recycled)
	return postAttackTurnDeck
}

// drawNextHand assembles the hand the next turn will see: starts from held (the attack turn's
// post-attack-turn hand), tops it up with end-of-turn draws from postAttackTurnDeck (capped at the
// deck's remaining size), and returns the result sorted by Card.ID() so findBest's cache
// key sees a canonical multiset.
//
// A fresh slice is always allocated: held aliases the attack turn's per-perm hand buffer (which
// the next playOneTurn would overwrite) and the attack turn's RemoveFromHand reorders via
// swap-with-last so held isn't sorted by ID. The sort must run on every path, including
// the toDraw == 0 case where the attack turn's output order would otherwise leak through.
func drawNextHand(postAttackTurnDeck *deck.Deck, held []card.Card, intellect int) []card.Card {
	toDraw := endOfTurnDraws(len(held), intellect)
	if toDraw > postAttackTurnDeck.Size() {
		toDraw = postAttackTurnDeck.Size()
	}
	nextHand := make([]card.Card, len(held), len(held)+toDraw)
	copy(nextHand, held)
	if toDraw > 0 {
		for _, c := range postAttackTurnDeck.Draw(toDraw) {
			nextHand = append(nextHand, c.(card.Card))
		}
	}
	sortHandByID(nextHand)
	return nextHand
}

// endOfTurnDraws reports how many cards the end-of-turn refill draws to bring a hand of
// heldLen up to intellect (zero when the hand already meets it). Sole source of truth for
// the refill target — the real refill and the partition's totalCards tiebreaker both
// route through here so they cannot drift. The real refill caps the result at its
// remaining deck size; the tiebreaker projects an uncapped refill.
func endOfTurnDraws(heldLen, intellect int) int {
	if n := intellect - heldLen; n > 0 {
		return n
	}
	return 0
}

// advanceToNextTurn clears per-turn ephemerals (value, cardsPlayed, ...). Idempotent on
// a freshly-built state. The deck pointer is left untouched — pool slots own a *Deck
// wrapper across resets, and the attack-turn runner rebinds its contents via ShallowCopyFrom.
func advanceToNextTurn(state *gameengine.GameState) {
	state.ResetEphemeralState()
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
	weapons []gameengine.Weapon,
	h []card.Card,
	d *deck.Deck,
	state *gameengine.GameState,
	ev *Evaluator,
) TurnSummary {
	return ev.Best(weapons, h, d, state)
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

// processAurasAtStartOfTurn fires every StartOfTurn aura handler queued on state. Handlers
// read and mutate state.Hand directly — both the dealt hand (visible at fire time per FaB
// turn order: previous turn's end-step refill ran before this turn's start triggers) and
// any reveals / discards / draws the handlers want to apply. Re-arms FiredThisTurn.
// Cascading reveals: a handler that pops d shrinks the view for the next, so two
// reveal-capable auras see distinct tops.
func processAurasAtStartOfTurn(state *gameengine.GameState, d *deck.Deck) {
	if !state.AnyAurasInPlay() {
		return
	}
	// Swap state's deck pointer to d so aura handlers (DrawOne / RevealTopOfDeck /
	// BanishTopOfDeck) mutate the master deck directly, then restore. The pool slot's
	// owned wrapper is preserved across the swap.
	saved := state.Deck()
	state.SetDeck(d)
	state.Engine().FireTriggers(card.FireContext{FiringType: triggertype.StartOfTurn})
	state.SetDeck(saved)
}

// pitchedFromBestLine returns BestLine's Pitch-role cards (excluding the arsenal-in slot,
// which never recycles), sorted by ID so recycled decks stay byte-identical across cache
// hits and from-scratch searches that pick different equally-optimal partitions.
func pitchedFromBestLine(line []card.CardAssignment) []card.Card {
	var out []card.Card
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		if a.Role == card.Pitch {
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
// fills snap with the attack turn-produced bestLine / cardsPlayed / swungWeapons / value, and
// attaches stats.PrintBest as the deferred replay closure.
func recordBestTurnFromSnap(stats *deck.Stats, summary TurnSummary, ev *Evaluator, snap *turnSnapshot) {
	lineCopy := make([]card.CardAssignment, len(summary.BestLine))
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
// whether the card was present in the dealt hand or arsenal-in slot when Best ran. bestLine
// covers both — each hand card has one entry, plus one FromArsenal entry for the arsenal
// card. presentBuf is a scratch slice indexed parallel to marginalBuf; both are caller-
// owned across turns to keep this path allocation-free.
func tallyMarginalPresence(marginalBuf []deck.CardMarginalStats, idIndex map[ids.CardID]int, presentBuf []bool, bestLine []card.CardAssignment, value float64) {
	if len(marginalBuf) == 0 {
		return
	}
	clear(presentBuf)
	for _, a := range bestLine {
		if i, ok := idIndex[a.Card.ID()]; ok {
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
