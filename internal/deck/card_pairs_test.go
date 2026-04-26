package deck

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestCardPairMutations_EmitsBothHalvesAbsent: a deck holding neither pair member produces
// pair mutations that add both halves and remove two distinct existing cards. With 2 unique
// non-pair IDs in deck the inner C(uniques, 2) loop emits exactly one candidate per absent
// pair.
func TestCardPairMutations_EmitsBothHalvesAbsent(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	muts := cardPairMutations(d, 2, nil)
	if len(muts) != len(cardPairs) {
		t.Fatalf("got %d pair mutations, want %d (one per absent pair, C(2,2)=1 candidates each)",
			len(muts), len(cardPairs))
	}

	// Each emitted mutation must (a) include both pair members exactly once, (b) drop one of
	// the deck's unique IDs each, (c) keep size stable at 4.
	for i, m := range muts {
		if len(m.Deck.Cards) != 4 {
			t.Errorf("mutation %d (%s): card count %d, want 4", i, m.Description, len(m.Deck.Cards))
		}
		counts := map[card.ID]int{}
		for _, c := range m.Deck.Cards {
			counts[c.ID()]++
		}
		// Find the matching pair by checking which one the mutation introduced.
		pairIdx := -1
		for j, p := range cardPairs {
			if counts[p.First] == 1 && counts[p.Second] == 1 {
				pairIdx = j
				break
			}
		}
		if pairIdx < 0 {
			t.Errorf("mutation %d (%s): no registered pair has both members at exactly 1 copy; counts=%v",
				i, m.Description, counts)
		}
	}
}

// TestCardPairMutations_SkipsWhenEitherHalfPresent: if the deck already holds one half of a
// pair, the pair generator skips that pair entirely (single-slot mutations cover the
// remaining half). With the only registered pair being Sun Kiss / Moon Wish, dropping one
// Sun Kiss into the deck zeroes out that pair's contribution.
func TestCardPairMutations_SkipsWhenEitherHalfPresent(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	sk := cards.Get(card.SunKissRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, sk})

	muts := cardPairMutations(d, 2, nil)
	for i, m := range muts {
		if strings.Contains(m.Description, "Moon Wish") || strings.Contains(m.Description, "Sun Kiss") {
			t.Errorf("mutation %d (%s): Sun Kiss/Moon Wish pair should be skipped (Sun Kiss already present)",
				i, m.Description)
		}
	}
}

// TestCardPairMutations_ResultDifferentFromSource: the no-op safeguard must hold — every
// emitted pair mutation produces a deck with a different card multiset than the source.
// Defensive against a future bug where a pair member overlaps with a removal target.
func TestCardPairMutations_ResultDifferentFromSource(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})
	srcKey := cardMultisetKey(d.Cards)
	for i, m := range cardPairMutations(d, 2, nil) {
		if cardMultisetKey(m.Deck.Cards) == srcKey {
			t.Errorf("mutation %d (%s) produced a no-op (same multiset as source)", i, m.Description)
		}
	}
}

// TestCardPairMutations_RespectsLegalFilter: a legal predicate that rejects either pair half
// suppresses pair mutations for that pair. Prevents the pair generator from quietly
// reintroducing format-banned cards via a pair add when single-slot mutations correctly
// filter them out.
func TestCardPairMutations_RespectsLegalFilter(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	// Reject Sun Kiss (Red), the second half of the pilot pair.
	legal := func(c card.Card) bool { return c.ID() != card.SunKissRed }
	muts := cardPairMutations(d, 2, legal)
	if len(muts) != 0 {
		t.Errorf("legal filter rejecting Sun Kiss (Red) should suppress all pair mutations; got %d",
			len(muts))
	}
}

// TestCardPairMutations_DeterministicOrdering: two back-to-back calls must produce the same
// mutation sequence. AllMutations consumers (the iterate-mode worker pool) rely on stable
// indexing for reproducibility under a fixed seed.
func TestCardPairMutations_DeterministicOrdering(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := cardPairMutations(d, 2, nil)
	second := cardPairMutations(d, 2, nil)
	if len(first) != len(second) {
		t.Fatalf("call counts differ: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].Description != second[i].Description {
			t.Errorf("mutation %d descriptions differ: %q vs %q",
				i, first[i].Description, second[i].Description)
		}
	}
}
