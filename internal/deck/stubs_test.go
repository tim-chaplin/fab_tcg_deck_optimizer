package deck

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// int1StubHero is a no-op hero with Intelligence=1 so tests can isolate per-hand behaviour
// without interaction between multiple drawn cards.
var int1StubHero = testutils.Hero{Intel: 1}

// deckFingerprint builds a comparable summary of a deck for equality checks in tests. Hashes
// the weapon loadout and a sorted card-count histogram so decks compare equal iff they would
// produce identical simulations.
func deckFingerprint(d *Deck) string {
	s := weaponKey(d.Weapons) + "|"
	counts := map[ids.CardID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	// Stable ordering — iterate over all possible IDs in byID order isn't exposed, so use a
	// sorted slice of (id, count).
	type pair struct {
		id ids.CardID
		n  int
	}
	var pairs []pair
	for id, n := range counts {
		pairs = append(pairs, pair{id, n})
	}
	// Insertion sort by id (tiny N).
	for i := 1; i < len(pairs); i++ {
		for j := i; j > 0 && pairs[j-1].id > pairs[j].id; j-- {
			pairs[j-1], pairs[j] = pairs[j], pairs[j-1]
		}
	}
	for _, p := range pairs {
		s += string(rune(p.id)) + ":" + string(rune(p.n)) + ","
	}
	return s
}
