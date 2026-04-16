// Package deck represents a candidate FaB deck and the hand-value stats accumulated from simulating
// it. Future deck-search code will create many Decks, evaluate each, and compare their Stats.
package deck

import (
	"fmt"
	"math/rand"

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

// Mutate creates a new deck by swapping one card (both copies) from `d` with a random card from
// the deckable pool that isn't already in the deck. The new deck has fresh (zero) stats. Weapons
// and hero are preserved.
func Mutate(d *Deck, rng *rand.Rand) *Deck {
	// Build a set of unique card names in the deck.
	inDeck := map[string]bool{}
	for _, c := range d.Cards {
		inDeck[c.Name()] = true
	}

	// Pick which unique card to remove.
	uniq := make([]string, 0, len(inDeck))
	for name := range inDeck {
		uniq = append(uniq, name)
	}
	removeName := uniq[rng.Intn(len(uniq))]

	// Pick a replacement from the pool that isn't already in the deck.
	pool := cards.Deckable()
	var replaceID cards.ID
	for {
		id := pool[rng.Intn(len(pool))]
		if !inDeck[cards.Get(id).Name()] {
			replaceID = id
			break
		}
	}
	replacement := cards.Get(replaceID)

	// Build the new card list: drop both copies of removeName, add two of replacement.
	newCards := make([]card.Card, 0, len(d.Cards))
	for _, c := range d.Cards {
		if c.Name() != removeName {
			newCards = append(newCards, c)
		}
	}
	newCards = append(newCards, replacement, replacement)
	return New(d.Hero, d.Weapons, newCards)
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
}

// BestHand records a single hand and its optimal play — used to surface the peak draw a deck saw
// during simulation.
type BestHand struct {
	Hand []card.Card
	Play hand.Play
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

	working := make([]card.Card, 0, deckSize)
	for r := 0; r < runs; r++ {
		working = append(working[:0], d.Cards...)
		rng.Shuffle(len(working), func(i, j int) {
			working[i], working[j] = working[j], working[i]
		})

		handIdx := 0
		for len(working) >= handSize {
			h := working[:handSize]
			play := hand.Best(d.Hero, d.Weapons, h, incomingDamage, working[handSize:])
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
					Hand: handCopy,
					Play: hand.Play{Roles: rolesCopy, Weapons: weaponsCopy, Value: play.Value},
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

			// Recycle: pitched cards go to the bottom (in hand order); attacked and defended cards are
			// spent. If nothing was pitched, the deck shrinks by handSize this turn.
			pitched := make([]card.Card, 0, handSize)
			for i, c := range h {
				if play.Roles[i] == hand.Pitch {
					pitched = append(pitched, c)
				}
			}
			working = append(working[handSize:], pitched...)
			handIdx++
		}
	}
	return d.Stats
}
