// Package deck represents a candidate FaB deck and the hand-value stats accumulated from simulating
// it. Future deck-search code will create many Decks, evaluate each, and compare their Stats.
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

// Deck is a hero plus equipped weapons and a deck of cards, along with the hand-value stats
// accumulated from simulating it.
type Deck struct {
	Hero    hero.Hero
	Weapons []weapon.Weapon
	Cards   []card.Card
	Stats   Stats
}

// New constructs a Deck with the given hero, weapons, and cards. Panics if the weapon loadout
// violates the "0–2 weapons; if 2, both must be 1H" equipment rule.
func New(h hero.Hero, weapons []weapon.Weapon, cards []card.Card) *Deck {
	validateWeapons(weapons)
	return &Deck{Hero: h, Weapons: weapons, Cards: cards}
}

// Random generates a random legal deck for `h`: a random weapon loadout from cards.AllWeapons
// (one 2H or two 1H, dual-wielding the same weapon allowed) and `size` cards drawn uniformly
// from cards.Deckable() in pairs — every included printing appears exactly twice (or more, up
// to `maxCopies`, if the same printing is rolled on multiple picks). `size` must be even and
// `maxCopies` must be at least 2.
func Random(h hero.Hero, size, maxCopies int, rng *rand.Rand) *Deck {
	if size%2 != 0 {
		panic(fmt.Sprintf("deck: Random requires even size (got %d) — cards are added in pairs", size))
	}
	if maxCopies < 2 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 2 (got %d) — cards are added in pairs", maxCopies))
	}
	loadouts := weaponLoadouts(cards.AllWeapons)
	weapons := loadouts[rng.Intn(len(loadouts))]

	pool := cards.Deckable()
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

// Mutation is one candidate single-slot change to a deck: the mutated Deck plus a human-readable
// summary of what changed (e.g. "swapped Aether Slash (Red) for Arcanic Spike (Red)"). Consumers
// use Deck to evaluate and Description for logging.
type Mutation struct {
	Deck        *Deck
	Description string
}

// AllMutations returns every single-card mutation of d in a deterministic order: first every
// alternative weapon loadout (sorted by loadout key), then every (removeID, addID) pair where
// one copy of removeID is dropped from the deck and one copy of addID is added. removeID must
// currently be in the deck; addID's post-mutation count must not exceed maxCopies. Pairs that
// would leave the deck unchanged (removeID == addID) are skipped.
//
// Ordering of card mutations: the outer loop iterates uniqueIDs by ascending per-card average
// contribution (d.Stats.PerCard[id].Avg()), so the cards carrying the least value in the current
// deck get swap candidates tried first. Cards without stats (the deck hasn't been evaluated
// yet) all tie at 0 and fall back to the card.ID tiebreak. The inner loop iterates the addID
// pool by card.ID.
//
// Iterate mode's hill climb adopts the FIRST mutation that beats the current best, so putting
// cheap-to-replace cards up front makes the early-improvement rounds meaningful rather than
// churn through high-value cards that are unlikely to lose their slot.
//
// Single-card swaps let the hill climber reach decks with odd per-card counts (e.g. 1× X +
// 3× Y at maxCopies=3, or 1× X with a hole filled elsewhere). The earlier "swap a whole pair
// for a whole pair" rule enforced 2-per-card artificially — with the sim fast enough, we let
// composition fall out of which configurations actually score higher.
//
// The returned decks have fresh (zero) stats and share no backing slices with d or each other.
func AllMutations(d *Deck, maxCopies int) []Mutation {
	var out []Mutation

	// Weapon mutations: every loadout different from the current one.
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

	// Card mutations: for each unique card in the deck, try adding any Deckable card whose
	// post-mutation count is still within maxCopies (including cards already in the deck below
	// the cap).
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	uniqueIDs := make([]card.ID, 0, len(counts))
	for id := range counts {
		uniqueIDs = append(uniqueIDs, id)
	}
	// Order removeID by ascending per-card avg contribution, tiebreaker on card.ID. When the
	// current deck has no stats (PerCard is nil / unseen id), Avg() returns 0 so every card ties
	// and the tiebreak falls through to stable card.ID order — same behaviour as before.
	sort.Slice(uniqueIDs, func(i, j int) bool {
		ai := d.Stats.PerCard[uniqueIDs[i]].Avg()
		aj := d.Stats.PerCard[uniqueIDs[j]].Avg()
		if ai != aj {
			return ai < aj
		}
		return uniqueIDs[i] < uniqueIDs[j]
	})

	pool := cards.Deckable()
	sort.Slice(pool, func(i, j int) bool { return pool[i] < pool[j] })

	for _, removeID := range uniqueIDs {
		removed := cards.Get(removeID)
		for _, addID := range pool {
			if addID == removeID {
				continue // no-op: remove one, add one of the same card.
			}
			if counts[addID] >= maxCopies {
				continue // already at max copies; can't add another.
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

// loadoutLabel formats a weapon loadout for mutation descriptions, e.g. "[Nebula Blade]" or
// "[Reaping Blade, Scepter of Pain]".
func loadoutLabel(ws []weapon.Weapon) string {
	if len(ws) == 0 {
		return "[]"
	}
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	sort.Strings(names)
	return "[" + strings.Join(names, ", ") + "]"
}

// weaponKey returns a comparable string for a weapon loadout so we can check equality.
func weaponKey(ws []weapon.Weapon) string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

// weaponLoadouts enumerates every legal equip combination from `ws`: each 2H weapon as a solo
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
	// PerCard attributes hand-level outcomes back to the cards that appeared in those hands. The
	// map is populated once per hand (after hand.Best picks the winning play), not per permutation
	// — attribution cost is negligible compared to the underlying search.
	PerCard map[card.ID]CardPlayStats
}

// BestTurn records a single hand and its optimal turn — used to surface the peak draw a deck
// saw during simulation. Summary.BestLine carries the cards and their assigned roles in
// canonical order; no parallel Hand slice is needed.
type BestTurn struct {
	Summary hand.TurnSummary
	// StartingRunechants is the Runechant count carried in from the previous turn when this hand
	// was played. Only meaningful for Runeblade heroes.
	StartingRunechants int
}

// CardPlayStats captures how a single card contributed to the decks it appeared in. Plays counts
// hands where it was played as an attack or defense; Pitches counts hands where it was spent for
// resources. TotalContribution sums a per-role accounting of what the card did on each
// appearance, filled in by hand.Best's tracked replay of the winning line:
//
//   - Pitch   → Card.Pitch() (1/2/3 resource value, treated as damage-equivalent per convention).
//   - Attack  → Card.Play() return plus the hero's OnCardPlayed trigger chained off it, captured
//     at the moment the card resolved in the winning attacker permutation — so conditional
//     riders, Runechant creations, and all other Play-time damage are attributed to the card
//     that actually did them.
//   - Defend  → the card's proportional share of min(sum_defense, incomingDamage), plus the
//     card's own Play return if it's a defense reaction.
//
// So the metric is "how much value does this card usually contribute, itself, to its hand" — as
// opposed to the hand's total value lumping every card together. Useful as a directional
// per-card signal; the Defense share is proportional rather than causal, so a defender that
// soaks all the block will look equal to a weaker one padding the same partition.
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

// Evaluate simulates `runs` shuffles of the deck. For each run it assembles successive hands of
// d.Hero.Intelligence() cards — Held cards from the previous turn plus fresh draws from the top
// of the deck — computes the optimal play against an opponent attacking for incomingDamage, and
// returns Pitched cards to the bottom of the deck (in hand order). Played and defended cards are
// spent; Held cards carry into the next hand so only (handSize - heldCount) fresh cards are
// drawn. Each run ends when the deck no longer has enough cards left to fill out the next hand.
//
// A "cycle" is one pass through the original deck size: cumulative hands 0..(deckSize/handSize - 1)
// are cycle 1, the next deckSize/handSize hands are cycle 2.
//
// Results accumulate into d.Stats and are also returned for convenience.
//
// Uses the package-level shared hand.Evaluator, which means repeated Evaluate calls on the same
// (hero, weapons) share their memo cache. Callers that evaluate multiple decks concurrently must
// use EvaluateWith with a goroutine-local Evaluator instead — the shared one has no internal
// synchronisation.
func (d *Deck) Evaluate(runs int, incomingDamage int, rng *rand.Rand) Stats {
	return d.EvaluateWith(runs, incomingDamage, rng, nil)
}

// EvaluateWith is Evaluate using the given hand.Evaluator. Pass a dedicated Evaluator per
// goroutine when running Evaluate calls in parallel; pass nil to reuse the package-level shared
// Evaluator (equivalent to calling Evaluate directly).
func (d *Deck) EvaluateWith(runs int, incomingDamage int, rng *rand.Rand, ev *hand.Evaluator) Stats {
	d.Stats.Runs += runs
	simstate.CurrentHero = d.Hero
	handSize := d.Hero.Intelligence()
	deckSize := len(d.Cards)
	if handSize <= 0 || deckSize < handSize {
		return d.Stats
	}
	handsPerCycle := deckSize / handSize

	// `buf` is a single-allocation reusable slab holding the current "deck state" for the run.
	// [head:tail] is the remaining deck in top-to-bottom order. Each iteration dealt cards are
	// consumed by advancing head; pitched cards are re-appended at tail. Sized 2×deckSize so there
	// is always room to append pitched cards before we need to compact head back to 0; compaction
	// (which shifts [head:tail] down) happens at most once every deckSize/handSize iterations.
	// The head/tail pointers and one-shot allocation keep the per-hand path allocation-free.
	buf := make([]card.Card, deckSize*2)
	// handBuf is the per-turn working hand: Held cards from last turn (prefix) + fresh draws.
	// heldBuf holds the Held cards between turns; next iteration copies them into handBuf. Both
	// are sized once per Evaluate so the inner loop stays allocation-free.
	handBuf := make([]card.Card, handSize)
	heldBuf := make([]card.Card, 0, handSize)
	nextHeld := make([]card.Card, 0, handSize)
	for r := 0; r < runs; r++ {
		copy(buf, d.Cards)
		// Inline Fisher-Yates: the closure-based rng.Shuffle would heap-allocate a func value
		// capturing buf on every run.
		for i := deckSize - 1; i > 0; i-- {
			j := rng.Intn(i + 1)
			buf[i], buf[j] = buf[j], buf[i]
		}

		head, tail := 0, deckSize
		handIdx := 0
		runechantCarryover := 0
		var arsenalCard card.Card
		heldBuf = heldBuf[:0]
		// maxHands caps the run at two full cycles through the deck. Without this a partition
		// that pitches everything and swings a weapon every turn recycles the same cards back
		// into the deck forever — hand.Best returns an identical TurnSummary each iteration so
		// head and tail advance in lockstep and the run never terminates. Two cycles is enough
		// to observe the shuffle's early and late game and matches the FirstCycle / SecondCycle
		// stats the caller already tracks.
		maxHands := 2 * handsPerCycle
		for handIdx < maxHands {
			// Fresh draws fill whatever held cards didn't take. If every slot is already held —
			// pathological but possible for tiny hands — the hand doesn't progress, so stop the
			// run rather than spin.
			drawCount := handSize - len(heldBuf)
			if drawCount == 0 || tail-head < drawCount {
				break
			}
			// Compact when there isn't room at the bottom to append a full hand's worth of
			// pitched cards without overrunning buf.
			if tail+handSize > len(buf) {
				copy(buf, buf[head:tail])
				tail -= head
				head = 0
			}
			// Assemble the hand: held prefix first, then fresh draws. Best() sorts the hand in
			// canonical order and Roles align to that post-sort order, so which slot each card
			// ends up in doesn't affect the held-vs-drawn distinction for anything downstream.
			h := handBuf[:handSize]
			copy(h, heldBuf)
			copy(h[len(heldBuf):], buf[head:head+drawCount])
			// Snapshot the starting carryover before Best overwrites it — the best-hand record
			// wants the count in play *when the hand was dealt*, not what remained after.
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
				// Clone BestLine and AttackChain — both alias memo-owned storage that a later
				// Best call may reuse.
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

			// Attribute per-card contribution from the winning BestLine. hand.Best has already
			// filled Contribution on each assignment. Held and Arsenal entries are "not played,
			// not pitched" so neither counter ticks; the Arsenal card's real contribution will
			// accrue on the later turn when it's played out of the slot. Arsenal-in assignments
			// (FromArsenal=true) belong to a previous turn's hand, so they don't contribute to
			// THIS hand's per-card stats either.
			if d.Stats.PerCard == nil {
				d.Stats.PerCard = map[card.ID]CardPlayStats{}
			}
			for _, a := range play.BestLine {
				if a.FromArsenal {
					continue
				}
				stat := d.Stats.PerCard[a.Card.ID()]
				switch a.Role {
				case hand.Pitch:
					stat.Pitches++
				case hand.Attack, hand.Defend:
					stat.Plays++
				}
				stat.TotalContribution += a.Contribution
				d.Stats.PerCard[a.Card.ID()] = stat
			}

			// Recycle: pitched hand cards go to the bottom of the remaining deck (buf[tail:]) in
			// hand order; attacked and defended cards are spent. The backing array has room
			// since the cards being "moved" are a subset of those we just consumed. Held cards
			// stay in the player's hand and get copied into nextHeld for the next turn. Arsenal
			// and arsenal-in entries are handled via arsenalCard / the loop's top-level arsenal
			// threading, not here.
			nextHeld = nextHeld[:0]
			for _, a := range play.BestLine {
				if a.FromArsenal {
					continue
				}
				switch a.Role {
				case hand.Pitch:
					buf[tail] = a.Card
					tail++
				case hand.Held:
					nextHeld = append(nextHeld, a.Card)
				}
			}
			head += drawCount
			handIdx++
			heldBuf, nextHeld = nextHeld, heldBuf
		}
	}
	return d.Stats
}

// IterateParallel runs one iterate-mode round. Workers share a single pool and each goroutine
// does BOTH the shallow screen for its pulled mutation AND — if the shallow result beats bestAvg
// — the deep-shuffles confirmation for that same mutation. The first worker to land a confirmed
// improvement publishes it; a cancellation atomic stops everyone else. Because deep confirms
// parallelise instead of serialising on the main goroutine, iterate rounds on noisy shallow
// screens (many false-positive passes) finish in max(shallow wall, deeps/workers × deep wall)
// rather than shallow wall + passes × deep wall.
//
// Mutations are pulled FIFO from the shared queue and workers start at the front, so the
// earliest-position-wins heuristic of serial iterate generally holds — but a worker locked on a
// deep confirm at position 20 doesn't block others from picking up position 25, so a later-
// position mutation can occasionally win if its deep confirm finishes first. The hill-climb
// trajectory stays comparable.
//
// ctx: aborts the round when Done — workers exit early and IterateParallel returns with
// found=false; caller distinguishes "aborted" from "local max" via ctx.Err().
// mutations: ordered list of candidates.
// bestAvg: current baseline (at deep-shuffles depth).
// shallowShuffles / deepShuffles / incoming: eval settings.
// numWorkers: goroutines in the pool; 0 uses runtime.GOMAXPROCS(0).
// seed: base seed for worker RNGs; worker w uses (seed + w) for shallow and a derived stream for
// deep so the two phases don't alias.
// shallowCompleted: optional atomic counter incremented once per shallow eval the worker pool
// finishes, so callers can render live "tested N/total" progress from a separate goroutine.
// deepsCompleted: optional counter incremented once per attempted deep confirm (regardless of
// outcome) so callers can also show deep-phase progress. Nil to opt out.
//
// Returns (improvedDeck, improvedAvg, improvedIndex, true) on first confirmed improvement, or
// (nil, bestAvg, -1, false) if no improvement was found OR ctx was cancelled.
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
	// Buffered to numWorkers so the first-to-finish sender never blocks; later senders can also
	// drop their improvement without waiting once the main goroutine has taken one.
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
			// Derive an independent deep stream so the two phases don't share rng state across
			// a single worker.
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
				// Fresh Deck for the deep pass so d.Stats from the shallow run doesn't bleed in.
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
		// A last-moment improvement may have landed just before the senders all returned.
		select {
		case imp := <-improvementCh:
			return imp.deck, imp.avg, imp.idx, true
		default:
		}
		return nil, bestAvg, -1, false
	}
}
