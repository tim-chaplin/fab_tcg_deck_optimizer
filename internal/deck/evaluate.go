package deck

// Hand-by-hand simulation of a Deck: Evaluate / EvaluateWith shuffle, walk two cycles of hands
// per run, and fold each turn's outcome into Stats; EvalOneTurnForTesting runs a single turn
// against a fixed card order for assertion-style tests. All cross-turn bookkeeping (held cards,
// arsenal, runechant carryover, start-of-turn AuraTrigger handling) lives here.

import (
	"fmt"
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

// Evaluate simulates runs shuffles of the deck. For each run it assembles successive hands of
// d.Hero.Intelligence() cards (Held cards from last turn plus fresh top-of-deck draws), computes
// the optimal play against an opponent attacking for incomingDamage, and recycles Pitched cards
// to the bottom of the deck in hand order. Played and defended cards are spent; Held cards carry
// into the next hand. A run ends when the deck can't fill the next hand.
//
// A "cycle" is one pass through the original deck size: hands 0..(deckSize/handSize - 1) are
// cycle 1, the next deckSize/handSize are cycle 2.
//
// Results accumulate into d.Stats and are returned for convenience.
//
// Uses the package-level shared hand.Evaluator. Concurrent callers must use EvaluateWith with a
// goroutine-local Evaluator — the shared buffers have no internal synchronisation.
func (d *Deck) Evaluate(runs int, incomingDamage int, rng *rand.Rand) Stats {
	return d.EvaluateWith(runs, incomingDamage, rng, nil)
}

// EvaluateWith is Evaluate using the given hand.Evaluator. Pass a dedicated Evaluator per
// goroutine for parallel runs; nil reuses the package-level shared Evaluator.
func (d *Deck) EvaluateWith(runs int, incomingDamage int, rng *rand.Rand, ev *hand.Evaluator) Stats {
	d.Stats.Runs += runs
	simstate.CurrentHero = d.Hero
	handSize := d.Hero.Intelligence()
	deckSize := len(d.Cards)
	if handSize <= 0 || deckSize < handSize {
		return d.Stats
	}
	handsPerCycle := deckSize / handSize

	// uniqueIDs / idIndex / presentBuf / marginalBuf back the per-turn marginal-stats
	// accounting. uniqueIDs lists every distinct card.ID that appears in d.Cards (one entry
	// per ID, in deck order of first appearance). idIndex maps an ID back to its position so
	// the per-turn presence walk over the dealt hand is O(handSize) map lookups instead of
	// an O(handSize × uniqueIDs) scan. presentBuf is reused each turn — zeroed via clear()
	// — to mark which uniqueIDs sat in this turn's dealt hand or arsenal-in slot.
	// marginalBuf accumulates the with/without sums in a flat slice so the inner loop avoids
	// per-turn map churn (~30ns × 2 ops × 21 IDs/turn would dominate Evaluate's hot path on
	// large anneal benchmarks); the slice is folded into Stats.PerCardMarginal once after
	// every shuffle finishes.
	uniqueIDs, idIndex := uniqueDeckIDs(d.Cards)
	presentBuf := make([]bool, len(uniqueIDs))
	marginalBuf := make([]CardMarginalStats, len(uniqueIDs))

	// buf is a single-allocation slab holding deck state for the run. [head:tail] is the
	// remaining deck in top-to-bottom order. Dealt cards advance head; pitched cards are
	// re-appended at tail. Sized 2×deckSize so there's always room to append before compacting;
	// compaction (shifting [head:tail] down) happens at most once per deckSize/handSize
	// iterations. The head/tail pointers keep the per-hand path allocation-free.
	buf := make([]card.Card, deckSize*2)
	// handBuf is the per-turn working hand: Held prefix + fresh draws. heldBuf holds Held
	// cards between turns. Sized once per Evaluate so the inner loop stays allocation-free.
	// handBuf's capacity exceeds handSize so a start-of-turn AuraTrigger reveal can append
	// the revealed card to the dealt hand without reallocating.
	handBuf := make([]card.Card, handSize, handSize+startOfTurnRevealRoom)
	heldBuf := make([]card.Card, 0, handSize)
	nextHeld := make([]card.Card, 0, handSize)
	// auraTriggerBuf carries AuraTriggers left alive at the end of last turn. Double-buffered
	// with nextAuraTrigger like heldBuf so the swap is allocation-free.
	auraTriggerBuf := make([]card.AuraTrigger, 0, handSize)
	nextAuraTrigger := make([]card.AuraTrigger, 0, handSize)
	for r := 0; r < runs; r++ {
		copy(buf, d.Cards)
		// Inline Fisher-Yates: rng.Shuffle would heap-allocate a closure over buf every run.
		for i := deckSize - 1; i > 0; i-- {
			j := rng.Intn(i + 1)
			buf[i], buf[j] = buf[j], buf[i]
		}

		head, tail := 0, deckSize
		handIdx := 0
		runechantCarryover := 0
		var arsenalCard card.Card
		heldBuf = heldBuf[:0]
		auraTriggerBuf = auraTriggerBuf[:0]
		// Cap the run at two full cycles. A pitch-everything-swing-a-weapon loop recycles the
		// same cards forever (hand.Best returns identical summaries each iteration, so head and
		// tail advance in lockstep); two cycles also match FirstCycle / SecondCycle stats.
		maxHands := 2 * handsPerCycle
		for handIdx < maxHands {
			h, drawCount, ok := dealNextHand(buf, handBuf, heldBuf, &head, &tail, handSize)
			if !ok {
				break
			}
			// Snapshot the starting carryover before Best overwrites it — the best-hand record
			// wants the count in play when the hand was dealt, not what remained after.
			startingRunechants := runechantCarryover
			// Snapshot the aura cards in play at the top of this turn (one entry per queued
			// AuraTrigger) before processTriggersAtStartOfTurn potentially destroys any. A
			// fresh slice keeps the snapshot stable once auraTriggerBuf is rewritten with the
			// survivors.
			var startOfTurnAuras []card.Card
			if len(auraTriggerBuf) > 0 {
				startOfTurnAuras = make([]card.Card, len(auraTriggerBuf))
				for i, t := range auraTriggerBuf {
					startOfTurnAuras[i] = t.Self
				}
			}
			// Process AuraTriggers carried in from last turn before the best-line search.
			// Survivors become this turn's priorAuraTriggers. Reveal handlers pop the deck top
			// and append it to the hand so the best-line search sees the augmented hand.
			var trigContribs []hand.TriggerContribution
			var trigDamage, trigRunes int
			var trigRevealed []card.Card
			auraTriggerBuf, trigContribs, trigDamage, trigRunes, trigRevealed, _ = processTriggersAtStartOfTurn(auraTriggerBuf, buf[head+drawCount:tail])
			for range trigRevealed {
				h = append(h, buf[head+drawCount])
				drawCount++
			}
			runechantCarryover += trigRunes
			// arsenalIn snapshots the arsenal slot's contents at the top of this turn, before
			// Best decides what to put in arsenal-out. Marginal stats key on arsenalIn so the
			// "card present in this turn's hand" set covers everything the solver had access to.
			arsenalIn := arsenalCard
			var play hand.TurnSummary
			if ev != nil {
				play = ev.BestWithTriggers(d.Hero, d.Weapons, h, incomingDamage, buf[head+drawCount:tail], runechantCarryover, arsenalCard, auraTriggerBuf)
			} else {
				play = hand.BestWithTriggers(d.Hero, d.Weapons, h, incomingDamage, buf[head+drawCount:tail], runechantCarryover, arsenalCard, auraTriggerBuf)
			}
			runechantCarryover = play.State.Runechants
			arsenalCard = play.State.Arsenal
			// Start-of-turn trigger credit is a flat additive on Value. Every partition
			// benefits equally so Best's ranking is unaffected, but Value must include it so
			// the best-hand pick and cycle averages reflect the real total.
			play.Value += trigDamage
			play.TriggersFromLastTurn = trigContribs
			play.StartOfTurnAuras = startOfTurnAuras
			v := float64(play.Value)

			d.Stats.TotalValue += v
			d.Stats.Hands++
			if d.Stats.Histogram == nil {
				d.Stats.Histogram = map[int]int{}
			}
			d.Stats.Histogram[play.Value]++
			if play.Value > d.Stats.Best.Summary.Value || len(d.Stats.Best.Summary.BestLine) == 0 {
				recordBestTurn(&d.Stats, play, startingRunechants)
			}
			switch handIdx / handsPerCycle {
			case 0:
				d.Stats.FirstCycle.Hands++
				d.Stats.FirstCycle.Total += v
			case 1:
				d.Stats.SecondCycle.Hands++
				d.Stats.SecondCycle.Total += v
			}

			tallyMarginalPresence(marginalBuf, idIndex, presentBuf, h, arsenalIn, v)
			nextHeld = applyTurnResult(play, buf, &head, &tail, nextHeld[:0])
			nextAuraTrigger = append(nextAuraTrigger[:0], play.State.AuraTriggers...)
			handIdx++
			heldBuf, nextHeld = nextHeld, heldBuf
			auraTriggerBuf, nextAuraTrigger = nextAuraTrigger, auraTriggerBuf
		}
	}
	mergeMarginalBuf(&d.Stats, uniqueIDs, marginalBuf)
	// Assemble the best turn's structured log once, after the loop, so the in-memory snapshot
	// and the on-disk JSON both carry the same shape. The JSON layer round-trips Log
	// verbatim; printing routes through hand.FormatTurnLog so live and reloaded decks render
	// identically.
	if len(d.Stats.Best.Summary.BestLine) > 0 {
		d.Stats.Best.Log = hand.BuildTurnLog(d.Stats.Best.Summary, d.Stats.Best.StartingRunechants)
	}
	return d.Stats
}

// startOfTurnRevealRoom caps how many cards a start-of-turn AuraTrigger reveal can append
// to a turn's dealt hand. Set larger than any plausible number of queued reveal-capable
// triggers so the per-turn handBuf never reallocates.
const startOfTurnRevealRoom = 8

// processTriggersAtStartOfTurn walks every AuraTrigger queued from last turn and does all
// the bookkeeping a turn boundary requires:
//
//   - Clears FiredThisTurn on every trigger regardless of Type, re-arming OncePerTurn gates.
//   - Fires every TriggerStartOfTurn handler against a shared TurnState seeded with the
//     post-draw deck, so handlers that peek the top read the card about to be revealed.
//   - Decrements Count on each fired trigger, drops the entry when Count hits zero, and
//     adds the destroyed aura to the start-of-turn graveyard so subsequent handlers see
//     it in state.Graveyard.
//   - Passes non-start-of-turn triggers through unchanged so they can fire mid-chain.
//
// Returns the survivor list, per-aura contributions for FormatBestTurn, the summed damage
// to fold into Value, Runechants created during the handlers (fed into next turn's
// carryover), cards the handlers moved from the deck top into the hand (ts.Revealed) in
// reveal order, and auras destroyed this pass in destroy order.
//
// Cascading reveals: a handler that pops s.Deck shrinks the view for the next handler, so
// two reveal-capable auras see distinct tops.
func processTriggersAtStartOfTurn(queued []card.AuraTrigger, postDrawDeck []card.Card) (
	survivors []card.AuraTrigger,
	contribs []hand.TriggerContribution,
	damage int,
	runes int,
	revealed []card.Card,
	graveyarded []card.Card,
) {
	if len(queued) == 0 {
		return queued[:0], nil, 0, 0, nil, nil
	}
	ts := card.TurnState{Deck: postDrawDeck}
	survivors = queued[:0]
	for _, t := range queued {
		// Re-arm the OncePerTurn gate before the start-of-turn fire so handlers that read
		// FiredThisTurn see the cleared state.
		t.FiredThisTurn = false
		if t.Type != card.TriggerStartOfTurn {
			survivors = append(survivors, t)
			continue
		}
		preReveal := len(ts.Revealed)
		d := t.Handler(&ts)
		damage += d
		// Attribute any newly-revealed card to this trigger so the best-turn printout can
		// show what the handler drew (e.g. Sigil of the Arknight: "drew X into hand"). Taking
		// ts.Revealed[preReveal] instead of counting from the end handles cascading reveals
		// where a later handler also appends — each trigger sees its own first-appended card.
		var revealed card.Card
		if len(ts.Revealed) > preReveal {
			revealed = ts.Revealed[preReveal]
		}
		contribs = append(contribs, hand.TriggerContribution{Card: t.Self, Damage: d, Revealed: revealed})
		t.Count--
		if t.Count > 0 {
			survivors = append(survivors, t)
			continue
		}
		// Aura destroyed — Self joins the start-of-turn graveyard so subsequent handlers see
		// it in state.Graveyard.
		ts.AddToGraveyard(t.Self)
	}
	return survivors, contribs, damage, ts.Runechants, ts.Revealed, ts.Graveyard
}

// applyTurnResult folds a completed turn's outcome into cross-turn state. The deck loop
// adopts play.State.Deck wholesale (cards mutated freely during the chain — DrawOne pops,
// alt-cost prepends, tutor removals — and the snapshot reflects every change), then
// recycles pitched-role cards to the bottom of buf per FaB rules. nextHeld is replaced with
// play.State.Hand, which carries partition Held cards plus anything tutored that didn't get
// played. Panics if buf is undersized — the standard 2×deckSize sizing leaves enough room
// for any plausible mid-chain growth, so a too-small buf signals a sizing bug at the caller.
func applyTurnResult(play hand.TurnSummary, buf []card.Card, head, tail *int, nextHeld []card.Card) []card.Card {
	newDeck := play.State.Deck
	pitched := pitchedFromBestLine(play.BestLine)
	totalLen := len(newDeck) + len(pitched)
	if cap(buf) < totalLen {
		panic(fmt.Sprintf("applyTurnResult: buf cap %d < required %d (newDeck=%d + pitched=%d) — caller under-sized buf",
			cap(buf), totalLen, len(newDeck), len(pitched)))
	}
	*head = 0
	copy(buf[:len(newDeck)], newDeck)
	copy(buf[len(newDeck):totalLen], pitched)
	*tail = totalLen
	nextHeld = nextHeld[:0]
	nextHeld = append(nextHeld, play.State.Hand...)
	return nextHeld
}

// pitchedFromBestLine returns the cards in BestLine assigned the Pitch role (excluding the
// arsenal-in slot, which never recycles into the deck). Used by applyTurnResult to put
// pitched cards on the deck bottom per FaB's end-of-turn pitch-zone-to-deck rule.
func pitchedFromBestLine(line []hand.CardAssignment) []card.Card {
	var out []card.Card
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		if a.Role == hand.Pitch {
			out = append(out, a.Card)
		}
	}
	return out
}

// dealNextHand fills handBuf with this turn's dealt hand: the held prefix from heldBuf followed
// by fresh top-of-deck draws, totaling handSize cards. Compacts buf[head:tail] down to buf[0:]
// when the tail doesn't have room for a full hand of pitched cards on the upcoming recycle.
// Returns the dealt hand (aliasing handBuf — successive calls overwrite it), the number of
// fresh draws consumed, and ok=false when the run can't progress: deck exhausted, the whole
// hand is already held with no room to draw, or last turn's start-of-turn reveal padded the
// hand past handSize and enough of those extras got Held to overflow handSize this turn.
func dealNextHand(buf, handBuf, heldBuf []card.Card, head, tail *int, handSize int) ([]card.Card, int, bool) {
	drawCount := handSize - len(heldBuf)
	if drawCount <= 0 || *tail-*head < drawCount {
		return nil, 0, false
	}
	if *tail+handSize > len(buf) {
		copy(buf, buf[*head:*tail])
		*tail -= *head
		*head = 0
	}
	h := handBuf[:handSize]
	copy(h, heldBuf)
	copy(h[len(heldBuf):], buf[*head:*head+drawCount])
	return h, drawCount, true
}

// TurnStartState captures the game state at the start of a turn: the hand just dealt, the card
// in the arsenal slot, the deck cards still to be drawn (top-to-bottom), the live Runechant
// count at the start of this turn, and the Value dealt by the previous turn (damage +
// prevention). Returned by EvalOneTurnForTesting.
type TurnStartState struct {
	Hand        []card.Card
	ArsenalCard card.Card
	Deck        []card.Card
	// Runechants is the live Runechant count at the start of this turn — leftover from the
	// previous turn's attack chain plus any tokens freshly created by start-of-turn
	// AuraTrigger handlers.
	Runechants int
	// PrevTurnValue is the total Value (damage dealt + damage prevented) the previous turn
	// produced — the same number hand.Best reports as TurnSummary.Value for that turn.
	PrevTurnValue int
	// PrevTurnBestLine is the winning role assignment from turn 1, so tests can assert which
	// card took which role.
	PrevTurnBestLine []hand.CardAssignment
	// PrevTurnGraveyard is the cards that ended up in the graveyard at the end of turn 1, in
	// the order they landed there. Sourced from hand.TurnSummary.Graveyard so tests can
	// distinguish "this card is in the graveyard" from "this card is just absent from the
	// next-turn surfaces."
	PrevTurnGraveyard []card.Card
	// StartOfTurnTriggerDamage is the damage-equivalent credited by turn-2's start-of-turn
	// AuraTrigger handlers — triggers registered during turn 1 that fired at the top of
	// turn 2. Zero when no trigger survived into the pass. Production callers fold this
	// into turn 2's Value; exposed here so tests can assert the cross-turn credit without
	// running turn 2 to completion.
	StartOfTurnTriggerDamage int
	// StartOfTurnGraveyard is the auras destroyed during turn-2's start-of-turn AuraTrigger
	// pass, in destroy order.
	StartOfTurnGraveyard []card.Card
}

// EvalOneTurnForTesting runs one turn against d.Cards in source order (no shuffle) and
// returns the turn-2 start state: the hand just dealt, the arsenal slot, the remaining
// deck, and the runechant carryover. arsenalIn seeds turn 1's arsenal slot (nil for empty).
// initialHand sets turn 1's starting hand; nil takes d.Cards[:handSize] as the hand and
// treats the rest as the deck, non-nil uses the slice directly (may be shorter than
// handSize) and treats d.Cards as the deck entirely. Test-only — production callers use
// Evaluate, which shuffles and loops.
func (d *Deck) EvalOneTurnForTesting(incomingDamage int, arsenalIn card.Card, initialHand []card.Card) TurnStartState {
	simstate.CurrentHero = d.Hero
	handSize := d.Hero.Intelligence()
	if handSize <= 0 {
		return TurnStartState{}
	}

	// Resolve turn 1's hand and the head offset. No caller-supplied hand: d.Cards[:handSize]
	// is the hand (default layout). Caller-supplied: d.Cards is the deck entirely, and the
	// hand is exactly what the caller handed in.
	var turn1Hand []card.Card
	var head int
	if initialHand == nil {
		if len(d.Cards) < handSize {
			return TurnStartState{}
		}
		turn1Hand = d.Cards[:handSize]
		head = handSize
	} else {
		if len(initialHand) == 0 || len(initialHand) > handSize {
			return TurnStartState{}
		}
		turn1Hand = initialHand
		head = 0
	}

	deckSize := len(d.Cards)
	// Oversized buf: 2×deckSize matches Evaluate's layout. Add a handSize cushion so small
	// decks still have room for mid-turn pitches (hand + drawn) without overflowing tail.
	buf := make([]card.Card, deckSize*2+handSize*2)
	copy(buf, d.Cards)
	// handBuf capacity matches Evaluate's so start-of-turn AuraTrigger reveals can append
	// without realloc.
	handBuf := make([]card.Card, handSize, handSize+startOfTurnRevealRoom)
	tail := deckSize

	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	play := hand.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head:tail], 0, arsenalIn)
	// drawCount=0: head already points past the starting hand, so applyTurnResult only needs
	// to advance past mid-turn draws.
	nextHeld := applyTurnResult(play, buf, &head, &tail, nil)
	triggerQueue := append([]card.AuraTrigger(nil), play.State.AuraTriggers...)

	// Deal turn 2's hand but stop short of running Best — the caller wants the pre-Best state.
	turn2Hand, drawCount2, ok := dealNextHand(buf, handBuf, nextHeld, &head, &tail, handSize)
	if !ok {
		return TurnStartState{
			ArsenalCard:       play.State.Arsenal,
			Runechants:        play.State.Runechants,
			PrevTurnValue:     play.Value,
			PrevTurnGraveyard: append([]card.Card(nil), play.State.Graveyard...),
		}
	}
	// Process turn-1 AuraTriggers at the turn-2 boundary the same way Evaluate does:
	// fire start-of-turn handlers, re-arm OncePerTurn gates, drop exhausted entries.
	// Reveals into the hand are consumed here so the returned turn-2 Hand matches what
	// Best would see.
	_, _, trigDamage, trigRunes, trigRevealed, trigGraveyarded := processTriggersAtStartOfTurn(triggerQueue, buf[head+drawCount2:tail])
	for range trigRevealed {
		turn2Hand = append(turn2Hand, buf[head+drawCount2])
		drawCount2++
	}
	handCopy := append([]card.Card(nil), turn2Hand...)
	deckLeft := append([]card.Card(nil), buf[head+drawCount2:tail]...)
	lineCopy := append([]hand.CardAssignment(nil), play.BestLine...)

	return TurnStartState{
		Hand:                     handCopy,
		ArsenalCard:              play.State.Arsenal,
		Deck:                     deckLeft,
		Runechants:               play.State.Runechants + trigRunes,
		PrevTurnValue:            play.Value,
		PrevTurnBestLine:         lineCopy,
		PrevTurnGraveyard:        append([]card.Card(nil), play.State.Graveyard...),
		StartOfTurnTriggerDamage: trigDamage,
		StartOfTurnGraveyard:     trigGraveyarded,
	}
}

// recordBestTurn clones the winning turn's slices into fresh storage and stamps stats.Best
// with the resulting BestTurn. Every slice in play (BestLine, SwungWeapons,
// TriggersFromLastTurn, StartOfTurnAuras, State.*) aliases scratch hand.Best may rewrite on
// the next call, so retaining them directly would let a later evaluation mutate the saved
// peak. Nil-length slices skip the clone so the captured hand.TurnSummary holds nil rather
// than a zero-length allocation.
func recordBestTurn(stats *Stats, play hand.TurnSummary, startingRunechants int) {
	lineCopy := make([]hand.CardAssignment, len(play.BestLine))
	copy(lineCopy, play.BestLine)
	var swungCopy []string
	if len(play.SwungWeapons) > 0 {
		swungCopy = append([]string(nil), play.SwungWeapons...)
	}
	var trigCopy []hand.TriggerContribution
	if len(play.TriggersFromLastTurn) > 0 {
		trigCopy = make([]hand.TriggerContribution, len(play.TriggersFromLastTurn))
		copy(trigCopy, play.TriggersFromLastTurn)
	}
	var aurasCopy []card.Card
	if len(play.StartOfTurnAuras) > 0 {
		aurasCopy = make([]card.Card, len(play.StartOfTurnAuras))
		copy(aurasCopy, play.StartOfTurnAuras)
	}
	stats.Best = BestTurn{
		Summary: hand.TurnSummary{
			BestLine:             lineCopy,
			SwungWeapons:         swungCopy,
			Value:                play.Value,
			State:                cloneCarryState(play.State),
			TriggersFromLastTurn: trigCopy,
			StartOfTurnAuras:     aurasCopy,
		},
		StartingRunechants: startingRunechants,
	}
}

// cloneCarryState deep-copies every slice in cs so the returned snapshot survives later
// hand.Best calls overwriting their backing arrays.
func cloneCarryState(cs hand.CarryState) hand.CarryState {
	out := hand.CarryState{
		Arsenal:    cs.Arsenal,
		Runechants: cs.Runechants,
	}
	if len(cs.Hand) > 0 {
		out.Hand = append([]card.Card(nil), cs.Hand...)
	}
	if len(cs.Deck) > 0 {
		out.Deck = append([]card.Card(nil), cs.Deck...)
	}
	if len(cs.Graveyard) > 0 {
		out.Graveyard = append([]card.Card(nil), cs.Graveyard...)
	}
	if len(cs.Banish) > 0 {
		out.Banish = append([]card.Card(nil), cs.Banish...)
	}
	if len(cs.AuraTriggers) > 0 {
		out.AuraTriggers = append([]card.AuraTrigger(nil), cs.AuraTriggers...)
	}
	if len(cs.Log) > 0 {
		out.Log = append([]string(nil), cs.Log...)
	}
	return out
}

// uniqueDeckIDs returns the distinct card IDs in cs (in deck order of first appearance) and
// a position-lookup map keyed by ID. The caller uses uniqueIDs to iterate every card the deck
// could ever score against and idIndex to flip per-turn presence flags from the dealt hand.
func uniqueDeckIDs(cs []card.Card) ([]card.ID, map[card.ID]int) {
	ids := make([]card.ID, 0, len(cs))
	idx := make(map[card.ID]int, len(cs))
	for _, c := range cs {
		id := c.ID()
		if _, seen := idx[id]; seen {
			continue
		}
		idx[id] = len(ids)
		ids = append(ids, id)
	}
	return ids, idx
}

// tallyMarginalPresence credits this turn's value to each entry in marginalBuf, bucketed by
// whether the card was present in the dealt hand or in the arsenal-in slot when hand.Best
// ran. presentBuf is a scratch slice indexed parallel to marginalBuf; the caller owns both
// across turns to keep this path allocation-free. Operates entirely on slices so the inner
// loop avoids the per-turn map churn a direct Stats.PerCardMarginal[id] update would cost.
func tallyMarginalPresence(marginalBuf []CardMarginalStats, idIndex map[card.ID]int, presentBuf []bool, dealt []card.Card, arsenalIn card.Card, value float64) {
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

// mergeMarginalBuf folds the per-Evaluate slice accumulator into Stats.PerCardMarginal,
// summing into existing entries so multiple Evaluate calls accumulate the same way PerCard
// does. The map is lazily initialised so decks that never get evaluated don't pay for an
// empty map.
func mergeMarginalBuf(stats *Stats, uniqueIDs []card.ID, marginalBuf []CardMarginalStats) {
	if len(uniqueIDs) == 0 {
		return
	}
	if stats.PerCardMarginal == nil {
		stats.PerCardMarginal = make(map[card.ID]CardMarginalStats, len(uniqueIDs))
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

// recycleCardStates prepares next turn's draw queue from this turn's assignments: pitched
// cards go to the bottom of buf[*tail:] (the backing array has room since moved cards are a
// subset of those just consumed); Held cards go into nextHeld for the next turn; attacked and
// defended cards are spent. Cards listed in returnedToTopOfDeck are skipped on the Held
// branch — applyTurnResult inserts those copies on top of the next-turn deck, so an
// additional nextHeld carry would double-count them. Arsenal / arsenal-in entries thread
// through arsenalCard separately, not here. Returns the updated nextHeld slice (pass a
// nil/empty slice or nextHeld[:0] to start).
func recycleCardStates(line []hand.CardAssignment, returnedToTopOfDeck []card.Card, buf []card.Card, tail *int, nextHeld []card.Card) []card.Card {
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		switch a.Role {
		case hand.Pitch:
			buf[*tail] = a.Card
			*tail++
		case hand.Held:
			if containsCardOnce(returnedToTopOfDeck, a.Card) {
				returnedToTopOfDeck = removeCardOnce(returnedToTopOfDeck, a.Card)
				continue
			}
			nextHeld = append(nextHeld, a.Card)
		}
	}
	return nextHeld
}

// containsCardOnce reports whether cs holds at least one occurrence of c (by ID). Linear
// scan; returnedToTopOfDeck lists are tiny (one entry per alt-cost-using card per chain) so a map
// would just add overhead.
func containsCardOnce(cs []card.Card, c card.Card) bool {
	for _, x := range cs {
		if x.ID() == c.ID() {
			return true
		}
	}
	return false
}

// removeCardOnce returns cs with the first occurrence of c (by ID) removed. Used by
// recycleCardStates to consume a returnedToTopOfDeck entry exactly once per matching BestLine slot,
// so a deck that holds two copies of a card and consumes only one via alt cost still carries
// the other to nextHeld.
func removeCardOnce(cs []card.Card, c card.Card) []card.Card {
	for i, x := range cs {
		if x.ID() == c.ID() {
			return append(cs[:i:i], cs[i+1:]...)
		}
	}
	return cs
}
