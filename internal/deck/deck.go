// Package deck represents a candidate FaB deck and the hand-value stats accumulated from simulating
// it. Future deck-search code will create many Decks, evaluate each, and compare their Stats.
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

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
// currently be in the deck; addID's post-mutation count must not exceed maxCopies. The outer
// loop iterates removeID, the inner loop addID, both sorted by card.ID; pairs that would leave
// the deck unchanged (removeID == addID) are skipped.
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
	sort.Slice(uniqueIDs, func(i, j int) bool { return uniqueIDs[i] < uniqueIDs[j] })

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
	// Best is the single highest-value hand seen across all runs (ties broken by first occurrence).
	// Hand is in canonical (post-sort) order aligned with Play.Roles. Zero-valued if no hands have
	// been evaluated.
	Best BestHand
	// PerCard attributes hand-level outcomes back to the cards that appeared in those hands. The
	// map is populated once per hand (after hand.Best picks the winning play), not per permutation
	// — attribution cost is negligible compared to the underlying search.
	PerCard map[card.ID]CardPlayStats
}

// BestHand records a single hand and its optimal play — used to surface the peak draw a deck saw
// during simulation.
type BestHand struct {
	Hand []card.Card
	Play hand.Play
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

// Evaluate simulates `runs` shuffles of the deck. For each run it draws successive hands of
// d.Hero.Intelligence() cards from the top, computes the optimal play against an opponent attacking
// for incomingDamage, and returns Pitched cards to the bottom of the deck (in hand order). Played
// and defended cards are spent. Each run ends when fewer than a full hand's worth of cards remain.
//
// A "cycle" is one pass through the original deck size: cumulative hands 0..(deckSize/handSize - 1)
// are cycle 1, the next deckSize/handSize hands are cycle 2.
//
// Results accumulate into d.Stats and are also returned for convenience.
func (d *Deck) Evaluate(runs int, incomingDamage int, rng *rand.Rand) Stats {
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
		for tail-head >= handSize {
			// Compact when there isn't room at the bottom to append a full hand's worth of
			// pitched cards without overrunning buf.
			if tail+handSize > len(buf) {
				copy(buf, buf[head:tail])
				tail -= head
				head = 0
			}
			h := buf[head : head+handSize]
			// Snapshot the starting carryover before Best overwrites it — the best-hand record
			// wants the count in play *when the hand was dealt*, not what remained after.
			startingRunechants := runechantCarryover
			play := hand.Best(d.Hero, d.Weapons, h, incomingDamage, buf[head+handSize:tail], runechantCarryover)
			runechantCarryover = play.LeftoverRunechants
			v := float64(play.Value)

			d.Stats.TotalValue += v
			d.Stats.Hands++
			if play.Value > d.Stats.Best.Play.Value || d.Stats.Best.Hand == nil {
				// Clone both slices — h aliases the working deck and play.Roles is owned by the
				// returned Play, which a later Best() call could reuse.
				handCopy := make([]card.Card, len(h))
				copy(handCopy, h)
				rolesCopy := make([]hand.Role, len(play.Roles))
				copy(rolesCopy, play.Roles)
				var weaponsCopy []string
				if len(play.Weapons) > 0 {
					weaponsCopy = make([]string, len(play.Weapons))
					copy(weaponsCopy, play.Weapons)
				}
				d.Stats.Best = BestHand{
					Hand:               handCopy,
					Play:               hand.Play{Roles: rolesCopy, Weapons: weaponsCopy, Value: play.Value},
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

			// Attribute per-card contribution using play.Contributions, which hand.Best fills in
			// from the winning attack-chain replay (accurate per-card damage including riders and
			// hero triggers) plus role-based shares for pitch and defense.
			if d.Stats.PerCard == nil {
				d.Stats.PerCard = map[card.ID]CardPlayStats{}
			}
			for i, c := range h {
				stat := d.Stats.PerCard[c.ID()]
				if play.Roles[i] == hand.Pitch {
					stat.Pitches++
				} else {
					stat.Plays++
				}
				if i < len(play.Contributions) {
					stat.TotalContribution += play.Contributions[i]
				}
				d.Stats.PerCard[c.ID()] = stat
			}

			// Recycle: pitched cards go to the bottom of the remaining deck (buf[tail:]) in hand
			// order; attacked and defended cards are spent. The backing array has room since the
			// cards being "moved" are a subset of those we just consumed.
			for i, c := range h {
				if play.Roles[i] == hand.Pitch {
					buf[tail] = c
					tail++
				}
			}
			head += handSize
			handIdx++
		}
	}
	return d.Stats
}
