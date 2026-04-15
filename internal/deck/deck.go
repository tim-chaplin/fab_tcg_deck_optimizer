// Package deck represents a candidate FaB deck and the hand-value stats accumulated from simulating
// it. Future deck-search code will create many Decks, evaluate each, and compare their Stats.
package deck

import (
	"fmt"
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
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
			play := hand.Best(d.Hero, d.Weapons, h, incomingDamage)
			v := float64(play.Value())

			d.Stats.TotalValue += v
			d.Stats.Hands++
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
