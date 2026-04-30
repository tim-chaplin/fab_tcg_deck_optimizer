package sim

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for
// the given chain length. Tests use this instead of hand-rolling the context fields so the
// common shape is centralised. Lives in a sim test file (rather than sim_test) so
// exports_test.go's NewSequenceContextForTest wrapper can reach it.
func newSequenceContextForTest(h Hero, pitched, deck []Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	bufs := newAttackBufs(chainLen, 0, nil)
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               deck,
		bufs:               bufs,
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
		carryWinner:        &bufs.carryWinnerScratch,
	}
}

// deckFingerprint builds a comparable summary of a deck for equality checks in tests. Hashes
// the weapon loadout and a sorted card-count histogram so decks compare equal iff they would
// produce identical simulations. Lives in package sim (rather than testutils) because it
// uses the unexported weaponKey helper.
func deckFingerprint(d *Deck) string {
	s := weaponKey(d.Weapons) + "|"
	counts := map[ids.CardID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
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
