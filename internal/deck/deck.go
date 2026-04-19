// Package deck represents a candidate FaB deck and the hand-value stats accumulated from
// simulating it. Search code creates many Decks, evaluates each, and compares their Stats.
package deck

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

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
// cards.Deckable() in pairs — every included printing appears exactly twice, or more up to
// maxCopies if the same printing is rolled multiple times. size must be even and maxCopies ≥ 2.
//
// legal filters the card pool: only IDs for which legal(cards.Get(id)) returns true are
// candidates. Pass nil for no filtering. Callers typically wire format.Format.IsLegal through
// here to restrict generation to a constructed format's banlist.
func Random(h hero.Hero, size, maxCopies int, rng *rand.Rand, legal func(card.Card) bool) *Deck {
	if size%2 != 0 {
		panic(fmt.Sprintf("deck: Random requires even size (got %d) — cards are added in pairs", size))
	}
	if maxCopies < 2 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 2 (got %d) — cards are added in pairs", maxCopies))
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
		if counts[id]+2 > maxCopies {
			continue
		}
		counts[id] += 2
		c := cards.Get(id)
		picks = append(picks, c, c)
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

// Avg returns the average hand value for this cycle.
func (c CycleStats) Avg() float64 {
	if c.Hands == 0 {
		return 0
	}
	return c.Total / float64(c.Hands)
}

// Avg returns the overall average hand value.
func (s Stats) Avg() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
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
	// between turns. Sized once per Evaluate so the inner loop stays allocation-free.
	handBuf := make([]card.Card, handSize)
	heldBuf := make([]card.Card, 0, handSize)
	nextHeld := make([]card.Card, 0, handSize)
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
		// Cap the run at two full cycles. A pitch-everything-swing-a-weapon loop recycles the
		// same cards forever (hand.Best returns identical summaries each iteration, so head and
		// tail advance in lockstep); two cycles also match FirstCycle / SecondCycle stats.
		maxHands := 2 * handsPerCycle
		for handIdx < maxHands {
			// Fresh draws fill whatever Held cards didn't take. If every slot is already held
			// (pathological for tiny hands), the hand doesn't progress — stop the run.
			drawCount := handSize - len(heldBuf)
			if drawCount == 0 || tail-head < drawCount {
				break
			}
			// Compact when the tail has no room for a full hand's pitched cards.
			if tail+handSize > len(buf) {
				copy(buf, buf[head:tail])
				tail -= head
				head = 0
			}
			// Assemble the hand: held prefix, then fresh draws. Best() sorts the hand into
			// canonical order and Roles align to that order, so slot position here is irrelevant
			// for anything downstream.
			h := handBuf[:handSize]
			copy(h, heldBuf)
			copy(h[len(heldBuf):], buf[head:head+drawCount])
			// Snapshot the starting carryover before Best overwrites it — the best-hand record
			// wants the count in play when the hand was dealt, not what remained after.
			startingRunechants := runechantCarryover
			var play hand.TurnSummary
			if ev != nil {
				play = ev.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head+drawCount:tail], runechantCarryover, arsenalCard)
			} else {
				play = hand.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head+drawCount:tail], runechantCarryover, arsenalCard)
			}
			runechantCarryover = play.LeftoverRunechants
			arsenalCard = play.ArsenalCard
			v := float64(play.Value)

			d.Stats.TotalValue += v
			d.Stats.Hands++
			if play.Value > d.Stats.Best.Summary.Value || len(d.Stats.Best.Summary.BestLine) == 0 {
				// Clone BestLine and AttackChain — both alias memo-owned storage a later Best
				// call may reuse.
				lineCopy := make([]hand.CardAssignment, len(play.BestLine))
				copy(lineCopy, play.BestLine)
				var chainCopy []hand.AttackChainEntry
				if len(play.AttackChain) > 0 {
					chainCopy = make([]hand.AttackChainEntry, len(play.AttackChain))
					copy(chainCopy, play.AttackChain)
				}
				d.Stats.Best = BestTurn{
					Summary: hand.TurnSummary{
						BestLine:    lineCopy,
						AttackChain: chainCopy,
						Value:       play.Value,
						ArsenalCard: play.ArsenalCard,
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
			nextHeld = recyclePlayedCards(play.BestLine, buf, &tail, nextHeld[:0])
			head += drawCount
			handIdx++
			heldBuf, nextHeld = nextHeld, heldBuf
		}
	}
	return d.Stats
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

// recyclePlayedCards prepares next turn's draw queue from this turn's assignments: pitched
// cards go to the bottom of buf[*tail:] (the backing array has room since moved cards are a
// subset of those just consumed); Held cards go into nextHeld for the next turn; attacked and
// defended cards are spent. Arsenal / arsenal-in entries thread through arsenalCard separately,
// not here. Returns the updated nextHeld slice (pass a nil/empty slice or nextHeld[:0] to
// start).
func recyclePlayedCards(line []hand.CardAssignment, buf []card.Card, tail *int, nextHeld []card.Card) []card.Card {
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
// both the shallow screen and — if shallow beats bestAvg — the deep-shuffles confirmation for
// the same mutation. The first worker to land a confirmed improvement wins; a cancellation
// atomic stops the others. Parallelising deep confirms makes rounds with noisy shallow screens
// finish in max(shallow wall, deeps/workers × deep wall) instead of shallow wall + passes × deep.
//
// Mutations are pulled FIFO so the earliest-position-wins heuristic of serial iterate generally
// holds, but a worker locked on a deep confirm at position 20 doesn't block position 25 — a
// later-position mutation can occasionally win if its deep confirm finishes first.
//
// ctx: aborts the round when Done; workers exit and IterateParallel returns found=false. The
// caller distinguishes "aborted" from "local max" via ctx.Err().
// mutations: ordered candidate list.
// bestAvg: current baseline (at deep-shuffles depth).
// shallowShuffles / deepShuffles / incoming: eval settings.
// numWorkers: goroutines; 0 uses runtime.GOMAXPROCS(0).
// seed: base seed; worker w uses (seed + w) for shallow and a derived stream for deep.
// shallowCompleted / deepsCompleted: optional atomic counters incremented per shallow eval and
// per attempted deep confirm, so callers can render live progress. Nil to opt out.
//
// Returns (improvedDeck, improvedAvg, improvedIndex, true) on first confirmed improvement, or
// (nil, bestAvg, -1, false) if none was found or ctx was cancelled.
func IterateParallel(
	ctx context.Context,
	mutations []Mutation,
	bestAvg float64,
	shallowShuffles, deepShuffles, incoming, numWorkers int,
	seed int64,
	shallowCompleted *atomic.Int64,
	deepsCompleted *atomic.Int64,
) (*Deck, float64, int, bool) {
	if numWorkers <= 0 {
		numWorkers = runtime.GOMAXPROCS(0)
	}
	if len(mutations) == 0 {
		return nil, bestAvg, -1, false
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	type improvement struct {
		idx  int
		avg  float64
		deck *Deck
	}
	// Buffered to numWorkers so the first sender never blocks and later senders can drop their
	// improvement without waiting once the main goroutine has taken one.
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
			// Derive an independent deep stream so the two phases don't share rng state.
			deepRng := rand.New(rand.NewSource(seed ^ (int64(workerIdx)+1)*int64(0x9e3779b9)))
			for i := range jobs {
				if innerCtx.Err() != nil {
					return
				}
				mut := mutations[i]
				d := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				shallowAvg := d.EvaluateWith(shallowShuffles, incoming, shallowRng, ev).Avg()
				if shallowCompleted != nil {
					shallowCompleted.Add(1)
				}
				if shallowAvg <= bestAvg {
					continue
				}
				if innerCtx.Err() != nil {
					return
				}
				// Fresh Deck for the deep pass so d.Stats from the shallow run doesn't leak in.
				dd := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				deepAvg := dd.EvaluateWith(deepShuffles, incoming, deepRng, ev).Avg()
				if deepsCompleted != nil {
					deepsCompleted.Add(1)
				}
				if deepAvg > bestAvg {
					select {
					case improvementCh <- improvement{idx: i, avg: deepAvg, deck: dd}:
					default:
						// Another worker already filled the buffer; drop silently.
					}
					cancel()
					return
				}
				// Deep rejected; keep pulling more shallow jobs.
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
		// A last-moment improvement may have landed just before all senders returned.
		select {
		case imp := <-improvementCh:
			return imp.deck, imp.avg, imp.idx, true
		default:
		}
		return nil, bestAvg, -1, false
	}
}
