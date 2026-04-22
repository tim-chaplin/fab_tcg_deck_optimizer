// Package deck represents a candidate FaB deck and the hand-value stats accumulated from
// simulating it. Search code creates many Decks, evaluates each, and compares their Stats.
package deck

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/klauspost/cpuid/v2"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Deck is a hero, equipped weapons, a deck of cards, and the simulated hand-value stats.
type Deck struct {
	Hero    hero.Hero
	Weapons []weapon.Weapon
	Cards   []card.Card
	Stats   Stats
}

// New constructs a Deck. Panics if the weapon loadout violates the "0–2 weapons; if 2, both 1H"
// equipment rule.
func New(h hero.Hero, weapons []weapon.Weapon, cards []card.Card) *Deck {
	validateWeapons(weapons)
	return &Deck{Hero: h, Weapons: weapons, Cards: cards}
}

// Random generates a random legal deck for h: a random weapon loadout from cards.AllWeapons
// (one 2H or two 1H; dual-wielding the same weapon allowed) and size cards drawn uniformly from
// cards.Deckable() one at a time, skipping any roll that would exceed maxCopies for the picked
// ID. Matches the single-slot granularity of deck.AllMutations so the hill-climb can explore
// the space the generator actually produces.
//
// legal filters the card pool: only IDs for which legal(cards.Get(id)) returns true are
// candidates. Pass nil for no filtering. Callers typically wire format.Format.IsLegal through
// here to restrict generation to a constructed format's banlist.
func Random(h hero.Hero, size, maxCopies int, rng *rand.Rand, legal func(card.Card) bool) *Deck {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 1 (got %d)", maxCopies))
	}
	loadouts := weaponLoadouts(cards.AllWeapons)
	weapons := loadouts[rng.Intn(len(loadouts))]

	pool := legalPool(legal)
	if len(pool) == 0 {
		panic("deck: Random's legal filter rejected every card — cannot build a deck")
	}
	counts := map[cards.ID]int{}
	picks := make([]card.Card, 0, size)
	for len(picks) < size {
		id := pool[rng.Intn(len(pool))]
		if counts[id]+1 > maxCopies {
			continue
		}
		counts[id]++
		picks = append(picks, cards.Get(id))
	}
	return New(h, weapons, picks)
}

// legalPool returns cards.Deckable() filtered by legal, or the full list if legal is nil.
// Shared by Random and AllMutations so both apply the same filter.
func legalPool(legal func(card.Card) bool) []cards.ID {
	pool := cards.Deckable()
	if legal == nil {
		return pool
	}
	filtered := pool[:0]
	for _, id := range pool {
		if legal(cards.Get(id)) {
			filtered = append(filtered, id)
		}
	}
	return filtered
}

// Mutation is one candidate single-slot change: the mutated Deck plus a human-readable summary
// (e.g. "swapped Aether Slash (Red) for Arcanic Spike (Red)"). Consumers use Deck to evaluate
// and Description for logging.
type Mutation struct {
	Deck        *Deck
	Description string
}

// AllMutations returns every single-card mutation of d in a deterministic order: first every
// alternative weapon loadout (sorted by loadout key), then every (removeID, addID) pair where
// one copy of removeID is dropped and one copy of addID is added. removeID must be in the deck;
// addID's post-mutation count must not exceed maxCopies. Pairs with removeID == addID are
// skipped.
//
// Card-mutation ordering: the outer loop iterates uniqueIDs by ascending per-card average
// contribution (d.Stats.PerCard[id].Avg()), so low-value cards get swap candidates tried first.
// Cards without stats tie at 0 and fall back to card.ID. The inner loop iterates the addID pool
// by card.ID. Favouring low-value removal slots surfaces useful swaps early for a
// first-improvement hill climb.
//
// Single-card swaps (not paired swaps) let the hill climber reach decks with odd per-card counts
// (e.g. 1× X + 3× Y at maxCopies=3).
//
// legal filters the addition pool: only accepted IDs become swap-in candidates, so format-banned
// cards can't be introduced. Removal targets aren't filtered — a deck that entered the climb
// holding a banned card can still have it swapped out. Pass nil to skip filtering.
//
// Returned decks have zero Stats and share no backing slices with d or each other.
func AllMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	out := weaponLoadoutMutations(d)
	out = append(out, cardSwapMutations(d, maxCopies, legal)...)
	return out
}

// weaponLoadoutMutations emits one Mutation per distinct weapon loadout that isn't the current
// one. Loadouts are canonicalised by weaponKey (names sorted) and processed in key order so
// the output is deterministic regardless of map-iteration randomness.
func weaponLoadoutMutations(d *Deck) []Mutation {
	loadouts := weaponLoadouts(cards.AllWeapons)
	currentKey := weaponKey(d.Weapons)
	type keyedLoadout struct {
		key     string
		weapons []weapon.Weapon
	}
	sortedLoadouts := make([]keyedLoadout, 0, len(loadouts))
	for _, l := range loadouts {
		sortedLoadouts = append(sortedLoadouts, keyedLoadout{key: weaponKey(l), weapons: l})
	}
	sort.Slice(sortedLoadouts, func(i, j int) bool { return sortedLoadouts[i].key < sortedLoadouts[j].key })
	var out []Mutation
	for _, l := range sortedLoadouts {
		if l.key == currentKey {
			continue
		}
		newCards := make([]card.Card, len(d.Cards))
		copy(newCards, d.Cards)
		out = append(out, Mutation{
			Deck:        New(d.Hero, l.weapons, newCards),
			Description: fmt.Sprintf("swapped weapons from %s to %s", loadoutLabel(d.Weapons), loadoutLabel(l.weapons)),
		})
	}
	return out
}

// cardSwapMutations emits every single-card remove+add mutation the deck admits. Remove targets
// iterate in ascending per-card avg contribution so the hill climb spends its budget on the
// currently-worst cards first; with no Stats yet the tiebreak falls through to stable card.ID
// order. Add candidates skip no-ops (same ID) and entries already at maxCopies.
func cardSwapMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	uniqueIDs := make([]card.ID, 0, len(counts))
	for id := range counts {
		uniqueIDs = append(uniqueIDs, id)
	}
	sort.Slice(uniqueIDs, func(i, j int) bool {
		ai := d.Stats.PerCard[uniqueIDs[i]].Avg()
		aj := d.Stats.PerCard[uniqueIDs[j]].Avg()
		if ai != aj {
			return ai < aj
		}
		return uniqueIDs[i] < uniqueIDs[j]
	})

	// legalPool returns IDs in ascending order (cards.Deckable() iterates byID).
	pool := legalPool(legal)

	var out []Mutation
	for _, removeID := range uniqueIDs {
		removed := cards.Get(removeID)
		for _, addID := range pool {
			if addID == removeID {
				continue // no-op: remove one and add one of the same card.
			}
			if counts[addID] >= maxCopies {
				continue // at max copies.
			}
			replacement := cards.Get(addID)
			newCards := make([]card.Card, 0, len(d.Cards))
			removed1 := false
			for _, c := range d.Cards {
				if !removed1 && c.ID() == removeID {
					removed1 = true
					continue
				}
				newCards = append(newCards, c)
			}
			newCards = append(newCards, replacement)
			out = append(out, Mutation{
				Deck:        New(d.Hero, d.Weapons, newCards),
				Description: fmt.Sprintf("-1 %s, +1 %s", removed.Name(), replacement.Name()),
			})
		}
	}
	return out
}

// sortedWeaponNames returns the weapon names in ascending order. The canonical form both
// loadoutLabel and weaponKey build on so two loadouts with the same weapons in different orders
// compare equal.
func sortedWeaponNames(ws []weapon.Weapon) []string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	sort.Strings(names)
	return names
}

// loadoutLabel formats a weapon loadout for mutation descriptions, e.g. "[Nebula Blade]" or
// "[Reaping Blade, Scepter of Pain]".
func loadoutLabel(ws []weapon.Weapon) string {
	if len(ws) == 0 {
		return "[]"
	}
	return "[" + strings.Join(sortedWeaponNames(ws), ", ") + "]"
}

// weaponKey returns a comparable string for a weapon loadout so we can check equality.
func weaponKey(ws []weapon.Weapon) string {
	return strings.Join(sortedWeaponNames(ws), ",")
}

// weaponLoadouts enumerates every legal equip combination from ws: each 2H weapon as a solo
// loadout, plus every unordered pair of 1H weapons (including dual-wielding the same weapon).
func weaponLoadouts(ws []weapon.Weapon) [][]weapon.Weapon {
	var oneHand, twoHand []weapon.Weapon
	for _, w := range ws {
		if w.Hands() == 1 {
			oneHand = append(oneHand, w)
		} else {
			twoHand = append(twoHand, w)
		}
	}
	var out [][]weapon.Weapon
	for _, w := range twoHand {
		out = append(out, []weapon.Weapon{w})
	}
	for i := 0; i < len(oneHand); i++ {
		for j := i; j < len(oneHand); j++ {
			out = append(out, []weapon.Weapon{oneHand[i], oneHand[j]})
		}
	}
	return out
}

func validateWeapons(weapons []weapon.Weapon) {
	switch len(weapons) {
	case 0, 1:
		return
	case 2:
		if weapons[0].Hands() != 1 || weapons[1].Hands() != 1 {
			panic("deck: two-weapon loadout requires both weapons to be 1H")
		}
	default:
		panic(fmt.Sprintf("deck: invalid weapon count %d (max 2)", len(weapons)))
	}
}

// Stats holds aggregate hand-value statistics across all simulated runs.
type Stats struct {
	Runs        int
	Hands       int
	TotalValue  float64
	FirstCycle  CycleStats
	SecondCycle CycleStats
	// Best is the single highest-value hand seen across all runs (ties broken by first
	// occurrence). Summary.BestLine is in canonical (post-sort) order. Zero-valued if no hands
	// have been evaluated.
	Best BestTurn
	// PerCard attributes hand-level outcomes back to the cards that appeared in those hands.
	// Populated once per hand after hand.Best picks the winner — attribution cost is negligible
	// next to the underlying search.
	PerCard map[card.ID]CardPlayStats
	// Histogram counts hands seen at each integer Value. Keyed by TurnSummary.Value so Min /
	// Median can be derived without retaining every hand's value. Nil until the first hand is
	// evaluated.
	Histogram map[int]int
}

// BestTurn records a single hand and its optimal turn — the peak draw a deck saw during
// simulation. Summary.BestLine carries the cards and roles in canonical order.
type BestTurn struct {
	Summary hand.TurnSummary
	// StartingRunechants is the Runechant count carried in from the previous turn when this hand
	// was played. Only meaningful for Runeblade heroes.
	StartingRunechants int
}

// CardPlayStats captures how a single card contributed across hands it appeared in. Plays counts
// hands where it attacked or defended; Pitches counts hands where it was spent for resources.
// TotalContribution sums role-specific credit from the winning-line replay:
//
//   - Pitch   → Card.Pitch() (1/2/3 resource value, damage-equivalent by convention).
//   - Attack  → Card.Play() return plus the hero's OnCardPlayed trigger chained off it, at the
//     moment the card resolved in the winning permutation.
//   - Defend  → proportional share of min(sumDefense, incomingDamage), plus the card's own
//     Play return if it's a defense reaction.
//
// Useful as a directional per-card signal. The Defense share is proportional not causal: a
// defender soaking the whole block looks equal to a weaker one padding the same partition.
type CardPlayStats struct {
	Plays             int
	Pitches           int
	TotalContribution float64
}

// Avg returns mean per-card contribution across every hand where this card appeared (Plays +
// Pitches). Returns 0 when the card was never seen.
func (c CardPlayStats) Avg() float64 {
	n := c.Plays + c.Pitches
	if n == 0 {
		return 0
	}
	return c.TotalContribution / float64(n)
}

// CycleStats tracks total value and hand count for a single deck cycle.
type CycleStats struct {
	Hands int
	Total float64
}

// Mean returns the arithmetic mean hand value for this cycle.
func (c CycleStats) Mean() float64 {
	if c.Hands == 0 {
		return 0
	}
	return c.Total / float64(c.Hands)
}

// Mean returns the overall arithmetic mean hand value.
func (s Stats) Mean() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
}

// Min returns the lowest Value any simulated hand produced. Zero when no hands have been seen.
func (s Stats) Min() int {
	if len(s.Histogram) == 0 {
		return 0
	}
	first := true
	m := 0
	for v := range s.Histogram {
		if first || v < m {
			m = v
			first = false
		}
	}
	return m
}

// Max returns the highest Value any simulated hand produced. Zero when no hands have been seen.
func (s Stats) Max() int {
	m := 0
	for v := range s.Histogram {
		if v > m {
			m = v
		}
	}
	return m
}

// Median returns the median hand value. With an even number of hands it's the mean of the two
// middle values (so it can be fractional). Zero when no hands have been seen.
func (s Stats) Median() float64 {
	if s.Hands == 0 || len(s.Histogram) == 0 {
		return 0
	}
	keys := make([]int, 0, len(s.Histogram))
	for v := range s.Histogram {
		keys = append(keys, v)
	}
	sort.Ints(keys)
	// Walk the sorted values in order, counting cumulative hands until we pass the median
	// rank(s). lower = rank s.Hands/2 (0-indexed); upper = rank (s.Hands-1)/2 for even Hands.
	lowerRank := (s.Hands - 1) / 2
	upperRank := s.Hands / 2
	var lower, upper int
	cum := 0
	foundLower := false
	for _, v := range keys {
		cum += s.Histogram[v]
		if !foundLower && cum > lowerRank {
			lower = v
			foundLower = true
		}
		if cum > upperRank {
			upper = v
			break
		}
	}
	return float64(lower+upper) / 2
}

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

	// buf is a single-allocation slab holding deck state for the run. [head:tail] is the
	// remaining deck in top-to-bottom order. Dealt cards advance head; pitched cards are
	// re-appended at tail. Sized 2×deckSize so there's always room to append before compacting;
	// compaction (shifting [head:tail] down) happens at most once per deckSize/handSize
	// iterations. The head/tail pointers keep the per-hand path allocation-free.
	buf := make([]card.Card, deckSize*2)
	// handBuf is the per-turn working hand: Held prefix + fresh draws. heldBuf holds Held cards
	// between turns. Sized once per Evaluate so the inner loop stays allocation-free. handBuf's
	// capacity exceeds handSize so a DelayedPlay reveal (Sigil of the Arknight) can append the
	// revealed card to the dealt hand without reallocating.
	handBuf := make([]card.Card, handSize, handSize+delayedHandRoom)
	heldBuf := make([]card.Card, 0, handSize)
	nextHeld := make([]card.Card, 0, handSize)
	// delayedBuf queues cards played last turn that implement card.DelayedPlay; their
	// PlayNextTurn fires at the start of this turn. Double-buffered with nextDelayed like the
	// held slices so the swap is allocation-free.
	delayedBuf := make([]card.Card, 0, handSize)
	nextDelayed := make([]card.Card, 0, handSize)
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
		delayedBuf = delayedBuf[:0]
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
			// Fire any PlayNextTurn callbacks queued from last turn before the best-line search.
			// The TurnState sees the post-draw deck so top-of-deck peeks (Sigil of the Arknight)
			// read the actual card about to be revealed. A ToHand result pops the revealed card
			// off the deck and appends it to the hand — the reveal cascades across queued cards
			// because runDelayedPlays advances its own Deck view between callbacks.
			delayedContribs, delayedDamage, delayedRunes, revealed := runDelayedPlays(delayedBuf, buf[head+drawCount:tail])
			for range revealed {
				h = append(h, buf[head+drawCount])
				drawCount++
			}
			runechantCarryover += delayedRunes
			var play hand.TurnSummary
			if ev != nil {
				play = ev.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head+drawCount:tail], runechantCarryover, arsenalCard)
			} else {
				play = hand.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head+drawCount:tail], runechantCarryover, arsenalCard)
			}
			runechantCarryover = play.LeftoverRunechants
			arsenalCard = play.ArsenalCard
			// Delayed credit is a flat additive on Value; every partition benefits equally so
			// Best's ranking is unaffected, but the resulting Value must include it so the
			// best-hand pick and cycle averages reflect the real total.
			play.Value += delayedDamage
			play.DelayedFromLastTurn = delayedContribs
			v := float64(play.Value)

			d.Stats.TotalValue += v
			d.Stats.Hands++
			if d.Stats.Histogram == nil {
				d.Stats.Histogram = map[int]int{}
			}
			d.Stats.Histogram[play.Value]++
			if play.Value > d.Stats.Best.Summary.Value || len(d.Stats.Best.Summary.BestLine) == 0 {
				// Clone BestLine, AttackChain, and Drawn — all three alias memo-owned storage a
				// later Best call may reuse. Drawn carries the turn's mid-turn-drawn cards
				// (Pitch / Attack / Held / Arsenal dispositions) so FormatBestTurn's per-card
				// breakdown reconciles with the displayed Value; LeftoverRunechants is captured
				// so the header's "carryover runechants" count is real instead of always zero.
				lineCopy := make([]hand.CardAssignment, len(play.BestLine))
				copy(lineCopy, play.BestLine)
				var chainCopy []hand.AttackChainEntry
				if len(play.AttackChain) > 0 {
					chainCopy = make([]hand.AttackChainEntry, len(play.AttackChain))
					copy(chainCopy, play.AttackChain)
				}
				var drawnCopy []hand.CardAssignment
				if len(play.Drawn) > 0 {
					drawnCopy = make([]hand.CardAssignment, len(play.Drawn))
					copy(drawnCopy, play.Drawn)
				}
				var delayedCopy []hand.DelayedContribution
				if len(play.DelayedFromLastTurn) > 0 {
					delayedCopy = make([]hand.DelayedContribution, len(play.DelayedFromLastTurn))
					copy(delayedCopy, play.DelayedFromLastTurn)
				}
				d.Stats.Best = BestTurn{
					Summary: hand.TurnSummary{
						BestLine:            lineCopy,
						AttackChain:         chainCopy,
						Drawn:               drawnCopy,
						Value:               play.Value,
						LeftoverRunechants:  play.LeftoverRunechants,
						ArsenalCard:         play.ArsenalCard,
						DelayedFromLastTurn: delayedCopy,
					},
					StartingRunechants: startingRunechants,
				}
			}
			switch handIdx / handsPerCycle {
			case 0:
				d.Stats.FirstCycle.Hands++
				d.Stats.FirstCycle.Total += v
			case 1:
				d.Stats.SecondCycle.Hands++
				d.Stats.SecondCycle.Total += v
			}

			attributePlayStats(&d.Stats, play.BestLine)
			nextHeld = applyTurnResult(play, buf, &head, &tail, drawCount, nextHeld[:0])
			nextDelayed = collectDelayedPlays(play.BestLine, nextDelayed[:0])
			handIdx++
			heldBuf, nextHeld = nextHeld, heldBuf
			delayedBuf, nextDelayed = nextDelayed, delayedBuf
		}
	}
	return d.Stats
}

// delayedHandRoom caps how many cards a DelayedPlay reveal can append to a turn's dealt hand.
// Set larger than any plausible number of queued reveal-capable cards so the per-turn handBuf
// never reallocates.
const delayedHandRoom = 8

// runDelayedPlays fires a PlayNextTurn callback on each queued card, passing a TurnState whose
// Deck is the post-draw deck so top-of-deck peeks read the actual card about to be revealed.
// Returns the per-card contributions for FormatBestTurn, the summed damage to fold into Value,
// the Runechant tokens created during the callbacks (folded into the next turn's carryover so
// same-hand variable-cost checks and ArcaneDamageDealt triggers see them), and the list of
// cards the callbacks want appended to the hand (ToHand results) in the order they were
// consumed from the deck top. All zero/nil when no cards are queued.
//
// Cascading reveals: each ToHand result shrinks the view of Deck by one entry so the next
// callback peeks at the new top. Matches the real sequencing where sigils destroy one after
// another at the start of the action phase.
func runDelayedPlays(queued []card.Card, postDrawDeck []card.Card) ([]hand.DelayedContribution, int, int, []card.Card) {
	if len(queued) == 0 {
		return nil, 0, 0, nil
	}
	contribs := make([]hand.DelayedContribution, 0, len(queued))
	var revealed []card.Card
	total := 0
	ts := card.TurnState{Deck: postDrawDeck}
	for _, c := range queued {
		dp := c.(card.DelayedPlay)
		r := dp.PlayNextTurn(&ts)
		total += r.Damage
		contribs = append(contribs, hand.DelayedContribution{Card: c, Damage: r.Damage, ToHand: r.ToHand})
		if r.ToHand != nil {
			revealed = append(revealed, r.ToHand)
			ts.Deck = ts.Deck[1:]
		}
	}
	return contribs, total, ts.Runechants, revealed
}

// collectDelayedPlays scans a turn's BestLine for cards whose Play ran (Role == Attack) and that
// implement card.DelayedPlay. Those are queued so their PlayNextTurn fires at the top of the
// next turn. Pitched / arsenaled / held copies are skipped — only the played instance has its
// aura in the arena.
func collectDelayedPlays(line []hand.CardAssignment, dst []card.Card) []card.Card {
	for _, a := range line {
		if a.Role != hand.Attack {
			continue
		}
		if _, ok := a.Card.(card.DelayedPlay); ok {
			dst = append(dst, a.Card)
		}
	}
	return dst
}

// applyTurnResult folds a completed turn's outcome into cross-turn state: pitched hand cards
// recycle to the deck bottom (via recycleCardStates), head advances past initial draws and
// every mid-turn-drawn card, and each drawn card is routed by disposition. Drawn-card Held
// entries append into nextHeld; Arsenal flows through play.ArsenalCard and needs no
// bookkeeping here.
func applyTurnResult(play hand.TurnSummary, buf []card.Card, head, tail *int, drawCount int, nextHeld []card.Card) []card.Card {
	nextHeld = recycleCardStates(play.BestLine, buf, tail, nextHeld)
	*head += drawCount + len(play.Drawn)
	for _, d := range play.Drawn {
		if d.Role == hand.Held {
			nextHeld = append(nextHeld, d.Card)
		}
	}
	return nextHeld
}

// dealNextHand fills handBuf with this turn's dealt hand: the held prefix from heldBuf followed
// by fresh top-of-deck draws, totaling handSize cards. Compacts buf[head:tail] down to buf[0:]
// when the tail doesn't have room for a full hand of pitched cards on the upcoming recycle.
// Returns the dealt hand (aliasing handBuf — successive calls overwrite it), the number of
// fresh draws consumed, and ok=false when the run can't progress (deck exhausted, or the whole
// hand is already held with no room to draw).
func dealNextHand(buf, handBuf, heldBuf []card.Card, head, tail *int, handSize int) ([]card.Card, int, bool) {
	drawCount := handSize - len(heldBuf)
	if drawCount == 0 || *tail-*head < drawCount {
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
	// previous turn's attack chain plus any tokens freshly created by PlayNextTurn callbacks
	// (e.g. Blessing of Occult's start-of-turn rune creation).
	Runechants int
	// PrevTurnValue is the total Value (damage dealt + damage prevented) the previous turn
	// produced — the same number hand.Best reports as TurnSummary.Value for that turn.
	PrevTurnValue int
	// PrevTurnBestLine is the winning role assignment from turn 1, so tests can assert which
	// card took which role (e.g. sigil played as Role=Attack vs. pitched vs. held).
	PrevTurnBestLine []hand.CardAssignment
}

// EvalOneTurnForTesting runs one turn against d.Cards in source order (no shuffle) and returns
// the turn-2 start state: the hand just dealt, the arsenal slot, the remaining deck, and the
// runechant carryover. arsenalIn seeds turn 1's arsenal slot (nil for empty). initialHand
// sets turn 1's starting hand; nil takes d.Cards[:handSize] as the hand and treats the rest
// as the deck, non-nil uses the slice directly (may be shorter than handSize) and treats
// d.Cards as the deck entirely. Test-only — pins cross-turn sim behaviour against a known
// deck layout. Production callers should use Evaluate, which shuffles and loops.
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
	// handBuf capacity matches Evaluate's so DelayedPlay reveals can append without realloc.
	handBuf := make([]card.Card, handSize, handSize+delayedHandRoom)
	tail := deckSize

	h := handBuf[:len(turn1Hand)]
	copy(h, turn1Hand)
	play := hand.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head:tail], 0, arsenalIn)
	// drawCount=0: head already points past the starting hand, so applyTurnResult only needs
	// to advance past mid-turn draws.
	nextHeld := applyTurnResult(play, buf, &head, &tail, 0, nil)
	delayedQueue := collectDelayedPlays(play.BestLine, nil)

	// Deal turn 2's hand but stop short of running Best — the caller wants the pre-Best state.
	turn2Hand, drawCount2, ok := dealNextHand(buf, handBuf, nextHeld, &head, &tail, handSize)
	if !ok {
		return TurnStartState{
			ArsenalCard:   play.ArsenalCard,
			Runechants:    play.LeftoverRunechants,
			PrevTurnValue: play.Value,
		}
	}
	// Fire turn-1 DelayedPlay callbacks now so top-of-deck reveals (Sigil of the Arknight) show
	// up in the returned turn-2 hand, matching production behaviour where Best sees the
	// augmented hand. Runes created by the callbacks fold into the turn-2 starting carryover.
	_, _, delayedRunes, revealed := runDelayedPlays(delayedQueue, buf[head+drawCount2:tail])
	for range revealed {
		turn2Hand = append(turn2Hand, buf[head+drawCount2])
		drawCount2++
	}
	handCopy := append([]card.Card(nil), turn2Hand...)
	deckLeft := append([]card.Card(nil), buf[head+drawCount2:tail]...)
	lineCopy := append([]hand.CardAssignment(nil), play.BestLine...)

	return TurnStartState{
		Hand:             handCopy,
		ArsenalCard:      play.ArsenalCard,
		Deck:             deckLeft,
		Runechants:       play.LeftoverRunechants + delayedRunes,
		PrevTurnValue:    play.Value,
		PrevTurnBestLine: lineCopy,
	}
}

// attributePlayStats folds the winning BestLine into per-card aggregates. hand.Best already
// filled Contribution on each assignment. Held / Arsenal entries don't tick either counter
// (Arsenal's real contribution accrues when it's played out of the slot on a later turn);
// FromArsenal entries belong to a previous turn's hand and don't contribute to this hand.
func attributePlayStats(stats *Stats, line []hand.CardAssignment) {
	if stats.PerCard == nil {
		stats.PerCard = map[card.ID]CardPlayStats{}
	}
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		stat := stats.PerCard[a.Card.ID()]
		switch a.Role {
		case hand.Pitch:
			stat.Pitches++
		case hand.Attack, hand.Defend:
			stat.Plays++
		}
		stat.TotalContribution += a.Contribution
		stats.PerCard[a.Card.ID()] = stat
	}
}

// recycleCardStates prepares next turn's draw queue from this turn's assignments: pitched
// cards go to the bottom of buf[*tail:] (the backing array has room since moved cards are a
// subset of those just consumed); Held cards go into nextHeld for the next turn; attacked and
// defended cards are spent. Arsenal / arsenal-in entries thread through arsenalCard separately,
// not here. Returns the updated nextHeld slice (pass a nil/empty slice or nextHeld[:0] to
// start).
func recycleCardStates(line []hand.CardAssignment, buf []card.Card, tail *int, nextHeld []card.Card) []card.Card {
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		switch a.Role {
		case hand.Pitch:
			buf[*tail] = a.Card
			*tail++
		case hand.Held:
			nextHeld = append(nextHeld, a.Card)
		}
	}
	return nextHeld
}

// IterateParallel runs one iterate-mode round. Workers share a queue and each goroutine does
// both the shallow screen and — if shallow beats the effective threshold — the deep-shuffles
// confirmation for the same mutation. The first worker to land an acceptable mutation wins; a
// cancellation atomic stops the others. Parallelising deep confirms makes rounds with noisy
// shallow screens finish in max(shallow wall, deeps/workers × deep wall) instead of
// shallow wall + passes × deep.
//
// Annealing: at temperature == 0 the function accepts only strict improvements (deepAvg >
// baseline) — the classical hill climb. At temperature > 0 it also accepts worse mutations
// with probability exp((deepAvg - baseline) / temperature), implementing a Metropolis-style
// simulated-annealing acceptance gate. The shallow pre-screen is widened proportionally
// (threshold = baseline - 3·T) so mutations likely to clear the probabilistic gate aren't cut
// off before they get to deep confirm.
//
// Mutations are pulled FIFO so the earliest-position-wins heuristic of serial iterate generally
// holds, but a worker locked on a deep confirm at position 20 doesn't block position 25 — a
// later-position mutation can occasionally win if its deep confirm finishes first.
//
// ctx: aborts the round when Done; workers exit and IterateParallel returns found=false.
// mutations: ordered candidate list.
// bestAvg: the current deck's avg (the "current state" in SA terms, not the all-time best).
// temperature: SA temperature for this round; 0 disables annealing.
// shallowShuffles / deepShuffles / incoming / numWorkers: eval settings.
// seed: base seed; worker w uses (seed + w) for shallow and a derived stream for deep and
// acceptance rolls.
// shallowCompleted / deepsCompleted: optional atomic counters for live progress.
//
// Returns (acceptedDeck, acceptedAvg, acceptedIndex, true) on first acceptance, or
// (nil, bestAvg, -1, false) if nothing cleared the gate or ctx was cancelled.
func IterateParallel(
	ctx context.Context,
	mutations []Mutation,
	bestAvg float64,
	temperature float64,
	shallowShuffles, deepShuffles, incoming, numWorkers int,
	seed int64,
	shallowCompleted *atomic.Int64,
	deepsCompleted *atomic.Int64,
) (*Deck, float64, int, bool) {
	if numWorkers <= 0 {
		numWorkers = defaultWorkers()
	}
	if len(mutations) == 0 {
		return nil, bestAvg, -1, false
	}

	// Shallow threshold mirrors the deep acceptance gate's reach: at T=0 it's strict (> bestAvg),
	// at T>0 it's widened to let probabilistically-acceptable mutations through. 3·T covers
	// ~95% of mutations that would be accepted (exp(-3) ≈ 0.05).
	shallowThreshold := bestAvg
	if temperature > 0 {
		shallowThreshold = bestAvg - 3*temperature
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	type improvement struct {
		idx  int
		avg  float64
		deck *Deck
	}
	improvementCh := make(chan improvement, numWorkers)

	jobs := make(chan int, len(mutations))
	for i := range mutations {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerIdx int) {
			defer wg.Done()
			ev := hand.NewEvaluator()
			shallowRng := rand.New(rand.NewSource(seed + int64(workerIdx)))
			// Derive an independent deep stream so the two phases don't share rng state. The
			// acceptance-roll rng shares the deep stream — the deep eval has already happened
			// by the time the roll runs, so no cross-influence on the deep result.
			deepRng := rand.New(rand.NewSource(seed ^ (int64(workerIdx)+1)*int64(0x9e3779b9)))
			for i := range jobs {
				if innerCtx.Err() != nil {
					return
				}
				mut := mutations[i]
				d := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				shallowAvg := d.EvaluateWith(shallowShuffles, incoming, shallowRng, ev).Mean()
				if shallowCompleted != nil {
					shallowCompleted.Add(1)
				}
				if shallowAvg <= shallowThreshold {
					continue
				}
				if innerCtx.Err() != nil {
					return
				}
				// Fresh Deck for the deep pass so d.Stats from the shallow run doesn't leak in.
				dd := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				deepAvg := dd.EvaluateWith(deepShuffles, incoming, deepRng, ev).Mean()
				if deepsCompleted != nil {
					deepsCompleted.Add(1)
				}
				if !acceptMutation(deepAvg, bestAvg, temperature, deepRng) {
					continue
				}
				select {
				case improvementCh <- improvement{idx: i, avg: deepAvg, deck: dd}:
				default:
					// Another worker already filled the buffer; drop silently.
				}
				cancel()
				return
			}
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

// defaultWorkers returns the worker count to use when callers pass numWorkers<=0. The workload
// is purely CPU-bound (shallow/deep sim, zero I/O), so SMT siblings end up fighting for the same
// cache and execution units rather than adding throughput — benchmarks on an 8-core / 16-thread
// Ryzen showed ~20% speedup capping at physical cores vs defaulting to GOMAXPROCS. cpuid exposes
// the physical count portably; we still clamp by GOMAXPROCS so an explicit GOMAXPROCS<physical
// (user override, container cgroup) wins.
func defaultWorkers() int {
	maxProcs := runtime.GOMAXPROCS(0)
	physical := cpuid.CPU.PhysicalCores
	if physical <= 0 || physical > maxProcs {
		return maxProcs
	}
	return physical
}

// acceptMutation implements the Metropolis acceptance rule. Strict improvements (deepAvg >
// bestAvg) always pass. Worse mutations pass with probability exp((deepAvg - bestAvg) / T)
// when T > 0; at T == 0 they're rejected, recovering the classical hill-climb behaviour.
func acceptMutation(deepAvg, bestAvg, temperature float64, rng *rand.Rand) bool {
	if deepAvg > bestAvg {
		return true
	}
	if temperature <= 0 {
		return false
	}
	prob := math.Exp((deepAvg - bestAvg) / temperature)
	return rng.Float64() < prob
}
