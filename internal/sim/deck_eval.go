package sim

// Hand-by-hand simulation of a Deck: (*Evaluator).Evaluate shuffles, walks two cycles of hands
// per run, and folds each turn's outcome into a fresh deck.Stats. All cross-turn bookkeeping
// (held cards, arsenal, runechant carryover, start-of-turn Aura handling) lives here. The
// single-turn assertion-style entry point EvalOneTurnForTesting lives in
// eval_one_turn_for_testing.go.

import (
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/weapon"
)

// Evaluate simulates `runs` shuffles of d. Each run assembles successive hands of
// d.Hero.Intelligence() cards (Held + fresh top-of-deck draws), computes the optimal play
// against mp, and recycles Pitched cards to deck bottom. A run ends when the deck can't
// fill the next hand. A "cycle" is one pass through the original deck size.
//
// Returns a fresh deck.Stats — callers that want to accumulate across multiple Evaluate
// calls maintain their own deck.Stats and merge returned values in.
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
	// adaptiveCheckInterval is the per-worker chunk size in the parallel-shuffle path —
	// after every numWorkers × adaptiveCheckInterval shuffles, the worker pool barrier-
	// merges into the run's aggregate stats and runs the adaptive stop check. 50 is the
	// empirical sweet spot:
	// dropping from 1000 → 50 cut anneal-bench wall-clock by 2.4× because random Viserai
	// decks converge in ~2000 shuffles instead of being forced to 8000 by a too-large chunk.
	// Going below 50 stops paying off as barrier overhead dominates.
	adaptiveCheckInterval = 50
	// adaptiveShufflesCap is the upper bound on the adaptive shuffle path. Caps a
	// pathological high-variance regime that doesn't converge — the run terminates at this
	// many shuffles even if the SE target was never hit. Sized for precision=0.01 (~82k
	// shuffles on a baseline Viserai deck) with headroom for higher-variance hero.
	adaptiveShufflesCap = 200000
)

// makeAdaptiveStop returns a shuffleStopper that fires when the per-turn mean's standard
// error drops below targetSE. Checks every adaptiveCheckInterval shuffles so the
// histogram walk doesn't run on every iteration.
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

// evaluateSequentialImpl runs the shuffle loop in the calling goroutine, using ev's
// cachedBufs scratch directly. This is the deterministic-RNG path tests rely on.
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
// ev.cache (RWMutex-protected) but each carry their own per-call scratch. Shuffles are
// processed in chunks of (numWorkers × adaptiveCheckInterval); after each chunk the main
// goroutine merges every worker's local deck.Stats into the running aggregate and runs the
// adaptive stop check. Per-worker RNG seeds are derived from rng.Int63() so the chunk
// distribution is deterministic given the input rng.
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
		// Drain spawned-many results from the buffered channel.
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
// Cross-turn aura / item / banished / graveyard / opponentMarked state is carried via the
// *gameengine.GameState master threaded across turns, so only the per-turn dealt-hand and
// the leftover held-card buffer live here.
type shuffleScratch struct {
	weaponsBuf  []weapon.Weapon
	handBuf     []card.Card
	heldBuf     []card.Card
	presentBuf  []bool
	marginalBuf []deck.CardMarginalStats
}

// newShuffleScratch sizes the per-shuffle reusable buffers for a deck of
// (weaponCount, deckSize, handSize) shape. Called once per worker (or once per evaluate-call
// for the sequential path); the returned scratch is hot-loop reused across every shuffle the
// worker runs.
func newShuffleScratch(weaponCount, _, handSize, numUniqueIDs int) *shuffleScratch {
	return &shuffleScratch{
		weaponsBuf:  make([]weapon.Weapon, weaponCount),
		handBuf:     make([]card.Card, handSize, handSize+startOfTurnRevealRoom),
		heldBuf:     make([]card.Card, 0, handSize),
		presentBuf:  make([]bool, numUniqueIDs),
		marginalBuf: make([]deck.CardMarginalStats, numUniqueIDs),
	}
}

// runOneShuffle simulates a single shuffle of the deck end-to-end (Copy → Shuffle → walk
// turns → record stats). Accumulates results into the caller-owned *deck.Stats. Both the
// sequential and parallel paths pass a local deck.Stats here; the parallel path merges
// per-worker totals at chunk boundaries via mergeStatsInto.
//
// masterDeck is the per-evaluation deck shared with mutation enumeration; runOneShuffle
// copies it before shuffling so each shuffle trial gets an independent deck and the
// master's Cards order stays stable across goroutines.
func runOneShuffle(masterDeck *deck.Deck, scratch *shuffleScratch, stats *deck.Stats, idIndex map[ids.CardID]int, ev *Evaluator, rng *rand.Rand, mp Matchup, handsPerCycle, handSize int) {
	d := masterDeck.Copy()
	d.Shuffle(rng)

	heroVal := d.Hero.(hero.Hero)
	weapons := scratch.weaponsBuf
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}

	// master is the start-of-turn carryover state — built once per shuffle and threaded
	// across turns via PrepareNextTurn. play.State (the post-chain winner) becomes the
	// next turn's master.
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
		// Snapshot the pre-process aura / item set for recordBestTurn (uses the start-of-
		// turn list, not the post-process survivors). startOfTurnAuras feeds the formatter's
		// "Auras at start of turn:" line.
		startingAuras := concreteAuras(master.Auras())
		startingItems := concreteItems(master.Items())
		startOfTurnAuras := snapshotStartOfTurnAuras(startingAuras)
		dealtHand := append([]card.Card(nil), h...)
		// processAurasAtStartOfTurn drives the deck through PopDeckTop on reveal-handling
		// auras (Sigil of the Arknight). Pass d directly so reveals mutate the eval-loop
		// deck in place — trigRevealed lists what got drawn.
		queued := concreteAuras(master.Auras())
		survivors, trigContribs, trigDamage, trigRevealed, _ := processAurasAtStartOfTurn(queued, d)
		// Replace master's aura set with survivors so Best sees the post-process state.
		master.ClearAuras()
		for _, a := range survivors {
			master.CreateAura(a)
		}
		for _, c := range trigRevealed {
			h = append(h, c)
		}
		arsenalIn := master.Arsenal()
		sortHandByID(h)
		play := runBestForTurn(weapons, h, mp, d, master, ev)
		play.Value += trigDamage
		play.TriggersFromLastTurn = trigContribs
		play.StartOfTurnAuras = startOfTurnAuras
		play.DealtHand = dealtHand

		if recordTurnStats(stats, play, handIdx, handsPerCycle) {
			replay := replayBestForTurnWithLog(weapons, h, mp, d, master, ev)
			replay.Value = play.Value
			replay.TriggersFromLastTurn = trigContribs
			replay.StartOfTurnAuras = startOfTurnAuras
			replay.DealtHand = dealtHand
			recordBestTurn(stats, replay, startingAuras, startingItems)
		}
		tallyMarginalPresence(scratch.marginalBuf, idIndex, scratch.presentBuf, h, arsenalIn, float64(play.Value))
		// Adopt the chain's post-mutation deck and recycle pitched cards onto the
		// bottom — FaB's end-of-turn pitch-zone-to-deck rule.
		d = play.State.Deck()
		pitched := pitchedFromBestLine(play.BestLine)
		recycled := make([]deck.Card, len(pitched))
		for i, c := range pitched {
			recycled[i] = c
		}
		d.PutBottom(recycled)
		// Carry hand leftover into next turn's heldBuf; thread play.State forward as master.
		heldBuf = append(heldBuf[:0], play.State.Hand()...)
		master = play.State
		master.ResetEphemeralState()
		// The chain runner installs a fresh per-leaf deck copy; a deck on the master is
		// dead weight every per-leaf Copy would duplicate. d already holds the post-turn
		// deck for the next iteration.
		master.SetDeck(nil)
		handIdx++
	}
	scratch.heldBuf = heldBuf
}

// concreteAuras casts a []gameengine.Aura slice to []*aura.Aura, allocating a fresh slice.
// Callers that need a snapshot frozen against subsequent master.ClearAuras / CreateAura
// mutations use this; lazy iteration over master.Auras() doesn't need the cast.
func concreteAuras(in []gameengine.Aura) []*aura.Aura {
	if len(in) == 0 {
		return nil
	}
	out := make([]*aura.Aura, 0, len(in))
	for _, a := range in {
		out = append(out, a.(*aura.Aura))
	}
	return out
}

// concreteItems is the items counterpart of concreteAuras.
func concreteItems(in []gameengine.Item) []*item.Item {
	if len(in) == 0 {
		return nil
	}
	out := make([]*item.Item, 0, len(in))
	for _, it := range in {
		out = append(out, it.(*item.Item))
	}
	return out
}

// mergeStatsInto folds src's per-shuffle accumulators into dst. Used by the parallel path
// to merge each worker's local deck.Stats into the run's aggregate after a chunk barrier.
// Histogram / PerCardMarginal merging is handled separately (the latter via
// mergeMarginalBuf at the end of the run).
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
	}
}

// snapshotStartOfTurnAuras returns a fresh slice of the Self cards backing every queued
// card-style Aura at the top of the turn — i.e. the card auras in play before
// processAurasAtStartOfTurn fires and potentially destroys any. Token-style auras
// (Runechants, …) are skipped: the formatter renders them separately. Returns nil when
// the queue holds no card auras.
func snapshotStartOfTurnAuras(queued []*aura.Aura) []card.Card {
	if len(queued) == 0 {
		return nil
	}
	var out []card.Card
	for _, t := range queued {
		if src := t.SourceCard(); src != nil {
			out = append(out, src.(card.Card))
		}
	}
	return out
}

// runBestForTurn dispatches to ev.BestSkipLog — the hot-path goroutine-local case. The
// returned TurnSummary has State.Log empty; replayBestForTurnWithLog re-runs with full
// Log materialisation when a turn becomes the new deck-best.
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

// replayBestForTurnWithLog re-runs the Best search with full Log materialisation. Same
// inputs and same algorithm as runBestForTurn — Best is deterministic given the inputs, so
// the returned TurnSummary has identical Value, BestLine, and CarryState to the SkipLog
// run, plus a fully populated State.Log. Used only when a turn becomes the new deck-best,
// so the replay cost amortises across the bulk of turns that don't.
func replayBestForTurnWithLog(
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
// TotalValue, lazily initialises the Histogram, and credits the value to FirstCycle /
// SecondCycle based on where handIdx sits relative to the deck's hands-per-cycle boundary.
//
// Returns true when this turn's Value beats the current stats.Best — the caller is then
// responsible for calling recordBestTurn with a TurnSummary that has its State.Log fully
// populated (replayed via replayBestForTurnWithLog when the SkipLog path was used). Keeping
// the recordBestTurn clone out of here means the SkipLog run isn't cloned uselessly when
// the caller plans to overwrite with the replayed result.
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

// startOfTurnRevealRoom caps how many cards a start-of-turn Aura reveal can append
// to a turn's dealt hand. Set larger than any plausible number of queued reveal-capable
// triggers so the per-turn handBuf never reallocates.
const startOfTurnRevealRoom = 8

// processAurasAtStartOfTurn walks every Aura queued from last turn and does all the
// bookkeeping a turn boundary requires:
//
//   - Clears FiredThisTurn on every trigger regardless of TriggerType, re-arming
//     OncePerTurn gates.
//   - Fires every TriggerStartOfTurn handler against a shared TurnState seeded with the
//     post-draw deck, so handlers that peek the top read the card about to be revealed.
//     Handlers that destroy themselves call ge.DestroyAura, which splices ts.auras
//     immediately and (when addToGraveyard) appends Self to the start-of-turn graveyard.
//   - Leaves non-start-of-turn auras in place so they can fire mid-chain.
//
// Returns the survivor list, per-aura contributions for FormatBestTurn, the summed damage
// to fold into Value, cards the handlers drew into the hand (ts.hand) in draw order, and
// auras destroyed this pass in destroy order.
//
// Cascading reveals: a handler that pops gs.Deck shrinks the view for the next handler, so
// two reveal-capable auras see distinct tops.
func processAurasAtStartOfTurn(queued []*aura.Aura, d *deck.Deck) (
	survivors []*aura.Aura,
	contribs []deck.TriggerContribution,
	damage int,
	revealed []card.Card,
	graveyarded []card.Card,
) {
	if len(queued) == 0 {
		return queued[:0], nil, 0, nil, nil
	}
	// Drive a fresh state + engine over the queued aura list. Reveal handlers (Sigil of
	// the Arknight) read the deck via PopDeckTop, which mutates d in place. Adopting
	// queued onto the state'gs aura list lets handlers' Destroy splice the live list
	// directly.
	gs := gameengine.GameStateBuilder().Build()
	gs.SetDeck(d)
	for _, a := range queued {
		gs.CreateAura(a)
	}
	ge := gs.Engine()
	// Walk queued in lockstep with FireStartOfTurn'gs callback: FireStartOfTurn visits
	// auras in ge.auras order, firing each TriggerStartOfTurn entry. We pre-capture each
	// firing aura'gs SourceCard so the contribution carries source identity even after
	// the aura destroys itself (and disappears from ge.auras).
	sourceByFireIdx := make([]card.Card, 0, len(queued))
	for _, a := range queued {
		if a.TriggerType() == triggertype.StartOfTurn {
			var src card.Card
			if s := a.SourceCard(); s != nil {
				src = s.(card.Card)
			}
			sourceByFireIdx = append(sourceByFireIdx, src)
		}
	}
	fireIdx := 0
	ge.FireStartOfTurn(func(_, dmg int, drawn card.Card, newEntries []turnlogger.LogEntry) {
		var text string
		if len(newEntries) > 0 {
			text = newEntries[0].Text
		}
		var src card.Card
		if fireIdx < len(sourceByFireIdx) {
			src = sourceByFireIdx[fireIdx]
		}
		contribs = append(contribs, deck.TriggerContribution{
			Card:     src,
			Damage:   dmg,
			Revealed: drawn,
			Text:     text,
		})
		damage += dmg
		fireIdx++
	})
	out := make([]*aura.Aura, 0, len(gs.Auras()))
	for _, a := range gs.Auras() {
		out = append(out, a.(*aura.Aura))
	}
	return out, contribs, damage, gs.Hand(), append([]card.Card(nil), gs.Graveyard()...)
}

// pitchedFromBestLine returns the cards in BestLine assigned the Pitch role (excluding the
// arsenal-in slot, which never recycles into the deck). Used by the eval loop to put
// pitched cards on the deck bottom per FaB's end-of-turn pitch-zone-to-deck rule. Sorted
// by ID so the deck-bottom recycle order is canonical-by-multiset rather than dependent
// on BestLine positional ordering — a same-multiset partition produced by the cache-
// replay path may have different positional roles than the from-scratch search would
// pick (different tie-break winner among optimal partitions), and without canonical
// sorting the two paths' recycled decks would diverge by card position even though the
// chain output is identical.
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

// sortHandByID sorts hand in place by Card.ID(). Called right before each turn's Best /
// runBestForTurn so the partition recurse always enumerates against the same canonical hand
// order — that drops a positional-tie-break source from findBest's leaf comparator and makes
// the cache-off and cache-on paths produce byte-identical results for matching multisets.
// Insertion sort instead of sort.SliceStable: hand sizes top out around 7 and the standard
// library's reflection-based interface comparator is ~30% of the bench's wall-clock at
// these sizes; insertion sort is stable and beats it handily for small N.
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

// recordBestTurn clones the winning turn's slices into fresh storage and stamps stats.Best
// with the resulting deck.BestTurn. Every slice in play (BestLine, SwungWeapons,
// TriggersFromLastTurn, StartOfTurnAuras, State.*) aliases scratch Best may rewrite on
// the next call, so retaining them directly would let a later evaluation mutate the saved
// peak. Nil-length slices skip the clone so the captured TurnSummary holds nil rather
// than a zero-length allocation.
func recordBestTurn(stats *deck.Stats, play TurnSummary, startingAuras []*aura.Aura, startingItems []*item.Item) {
	lineCopy := make([]deck.CardAssignment, len(play.BestLine))
	copy(lineCopy, play.BestLine)
	stats.Best = deck.BestTurn{
		Value:    play.Value,
		BestLine: lineCopy,
		Log:      BuildTurnLog(play, startingAuras, startingItems),
	}
}

// tallyMarginalPresence credits this turn's value to each entry in marginalBuf, bucketed by
// whether the card was present in the dealt hand or in the arsenal-in slot when Best
// ran. presentBuf is a scratch slice indexed parallel to marginalBuf; the caller owns both
// across turns to keep this path allocation-free. Operates entirely on slices so the inner
// loop avoids the per-turn map churn a direct PerCardMarginal[id] update would cost.
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

// mergeMarginalBuf folds the per-Evaluate slice accumulator into PerCardMarginal on the
// supplied deck.Stats. The map is lazily initialised so unscored decks don't pay for an
// empty map.
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
