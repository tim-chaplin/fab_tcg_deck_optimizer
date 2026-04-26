package deck

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestCardPairMutations_EnumeratesAllVariantCrossProducts: with neither pair half present,
// the generator must emit a candidate per (firstVariant, secondVariant) cross-product per
// (i, j) low-avg removal pair. For the Sun Kiss / Moon Wish pair (3 × 3 = 9 cross-products),
// a deck with 2 unique non-pair IDs (1 (i,j) pair) yields 9 candidates per pair entry.
func TestCardPairMutations_EnumeratesAllVariantCrossProducts(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	muts := cardPairMutations(d, 2, nil)
	const variantCombosPerPair = 3 * 3
	const removalCombos = 1 // C(2, 2)
	want := len(cardPairs) * variantCombosPerPair * removalCombos
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (%d pairs × %d variant combos × %d removal pairs)",
			len(muts), want, len(cardPairs), variantCombosPerPair, removalCombos)
	}

	// Every (firstID, secondID) cross-product from cardPairs[0] must appear at least once.
	type combo struct{ first, second card.ID }
	seen := map[combo]bool{}
	for _, m := range muts {
		counts := map[card.ID]int{}
		for _, c := range m.Deck.Cards {
			counts[c.ID()]++
		}
		for _, fID := range cardPairs[0].First {
			for _, sID := range cardPairs[0].Second {
				if counts[fID] == 1 && counts[sID] == 1 {
					seen[combo{fID, sID}] = true
				}
			}
		}
	}
	wantCombos := len(cardPairs[0].First) * len(cardPairs[0].Second)
	if len(seen) != wantCombos {
		t.Errorf("variant cross-product coverage: saw %d distinct (first, second) pairs, want %d",
			len(seen), wantCombos)
	}
}

// TestCardPairMutations_FiresWhenOneHalfAlreadyPresent: pair mutations fire whenever the
// per-variant maxCopies cap allows the add. With one half partially present, the climber
// can still grow the OTHER variant of that side as a pair-shape mutation rather than two
// sequential single-slot swaps.
//
// Same-variant overlap suppression (e.g. -1 SunKissRed + +1 SunKissRed reducing to a
// single-slot) is the orthogonal optimisation tested in
// TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps; here we check that
// non-overlapping variant combinations still emit despite Sun Kiss (Red) being a removal
// candidate.
func TestCardPairMutations_FiresWhenOneHalfAlreadyPresent(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	sk := cards.Get(card.SunKissRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, sk})

	muts := cardPairMutations(d, 2, nil)
	if len(muts) == 0 {
		t.Fatal("expected pair mutations even with one half present")
	}
	// Sun Kiss (Yellow) and Sun Kiss (Blue) are absent (count 0), so cross-product adds for
	// those variants must still fire — at least one must combine with any Moon Wish variant.
	sawDifferentSunKissVariantAdd := false
	for _, m := range muts {
		if strings.Contains(m.Description, "+1 Sun Kiss (Yellow)") ||
			strings.Contains(m.Description, "+1 Sun Kiss (Blue)") {
			sawDifferentSunKissVariantAdd = true
			break
		}
	}
	if !sawDifferentSunKissVariantAdd {
		t.Error("expected at least one mutation adding a non-Red Sun Kiss variant " +
			"when Red is already present")
	}
}

// TestCardPairMutations_RespectsMaxCopiesPerVariant: a pair-add variant whose post-mutation
// count would exceed maxCopies must be skipped. Drives this with an unbalanced deck: 2 copies
// of Sun Kiss (Red) (saturated at maxCopies=2) means no SunKissRed adds, but SunKissYellow /
// SunKissBlue adds remain valid. Mirror check on the Moon Wish side: 2 MoonWishRed saturates
// MoonWishRed adds.
func TestCardPairMutations_RespectsMaxCopiesPerVariant(t *testing.T) {
	skR := cards.Get(card.SunKissRed)
	mwR := cards.Get(card.MoonWishRed)
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
		[]card.Card{skR, skR, mwR, mwR, a, b})

	muts := cardPairMutations(d, 2, nil)
	for i, m := range muts {
		// No mutation should add SunKissRed or MoonWishRed (both at cap).
		if strings.Contains(m.Description, "+1 Sun Kiss (Red)") {
			t.Errorf("mutation %d (%s): added Sun Kiss (Red) despite saturation at maxCopies=2",
				i, m.Description)
		}
		if strings.Contains(m.Description, "+1 Moon Wish (Red)") {
			t.Errorf("mutation %d (%s): added Moon Wish (Red) despite saturation at maxCopies=2",
				i, m.Description)
		}
	}
	// The non-saturated cross-products (Yellow/Blue × Yellow/Blue = 4 combos) should still
	// emit. With 4 unique IDs in deck (sk-red, mw-red, a, b), C(4, 2) = 6 (i,j) pairs.
	// Overlap suppression only fires when a removal target equals an ADD id; here the adds
	// are Yellow/Blue variants so none of the deck removals overlap, leaving all 6 (i,j)
	// pairs eligible. Total = 4 combos × 6 removals = 24.
	const wantCombos = 4
	const wantRemovals = 6
	want := wantCombos * wantRemovals
	if len(muts) != want {
		t.Errorf("got %d pair mutations with saturated red variants, want %d (%d combos × %d removals)",
			len(muts), want, wantCombos, wantRemovals)
	}
}

// TestCardPairMutations_HandlesUnbalancedHalfCounts: the generator should work with arbitrary
// per-variant counts of each half. Drives this with a deck holding 5 Moon Wish (across
// variants) and 3 Sun Kiss (across variants) at maxCopies=3 — a realistic "we found the
// pair, now anneal can fine-tune the per-variant balance" mid-climb state. Pair mutations
// should fire only for variants below cap, and the generated decks must remain at the
// original card count.
func TestCardPairMutations_HandlesUnbalancedHalfCounts(t *testing.T) {
	mwR := cards.Get(card.MoonWishRed)
	mwY := cards.Get(card.MoonWishYellow)
	mwB := cards.Get(card.MoonWishBlue)
	skR := cards.Get(card.SunKissRed)
	skY := cards.Get(card.SunKissYellow)
	skB := cards.Get(card.SunKissBlue)
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	// 2 + 2 + 1 = 5 Moon Wish; 1 + 1 + 1 = 3 Sun Kiss; padding stays under maxCopies=3 per
	// non-pair card (3 Arcanic Crackle, 3 Arcanic Spike).
	cardsList := []card.Card{
		mwR, mwR, mwY, mwY, mwB,
		skR, skY, skB,
		a, a, a, b, b, b,
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, cardsList)

	muts := cardPairMutations(d, 3, nil)
	if len(muts) == 0 {
		t.Fatal("expected pair mutations on unbalanced deck")
	}
	for i, m := range muts {
		if len(m.Deck.Cards) != len(cardsList) {
			t.Errorf("mutation %d (%s): card count %d, want %d (size must stay stable)",
				i, m.Description, len(m.Deck.Cards), len(cardsList))
		}
		// At maxCopies=3, MoonWishRed+1=3 ≤ 3 (eligible), MoonWishYellow+1=3 ≤ 3 (eligible),
		// MoonWishBlue+1=2 ≤ 3 (eligible). Same for Sun Kiss (all 1+1=2 ≤ 3). So no variant
		// is saturated at maxCopies=3 — every cross-product should be a candidate add.
		// Mostly we're checking the generator doesn't panic or emit malformed mutations on
		// the unbalanced shape.
		counts := map[card.ID]int{}
		for _, c := range m.Deck.Cards {
			counts[c.ID()]++
			if counts[c.ID()] > 3 {
				t.Errorf("mutation %d (%s): card %s exceeds maxCopies=3 (count %d)",
					i, m.Description, c.Name(), counts[c.ID()])
			}
		}
	}
}

// TestCardPairMutations_ResultDifferentFromSource: every emitted pair mutation produces a
// deck with a different card multiset than the source. Defensive against a future bug where
// the overlap-suppression check misses a path that ends up at the source composition.
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

// TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps: when a removal target is
// itself a pair member, the resulting mutation reduces to a single-slot swap (the matching
// pair member's count is unchanged after -1 +1). Single-slot already covers that, so the
// pair generator skips those combos. Drives this with a deck containing Sun Kiss (Red) as a
// removal candidate and verifies no mutation removes and re-adds Sun Kiss (Red).
func TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps(t *testing.T) {
	skR := cards.Get(card.SunKissRed)
	a := cards.Get(card.ArcanicCrackleRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{skR, a, a, a})
	for i, m := range cardPairMutations(d, 2, nil) {
		if strings.Contains(m.Description, "-1 Sun Kiss (Red)") &&
			strings.Contains(m.Description, "+1 Sun Kiss (Red)") {
			t.Errorf("mutation %d (%s): redundant -1/+1 of Sun Kiss (Red) — overlap suppression failed",
				i, m.Description)
		}
	}
}

// TestCardPairMutations_RespectsLegalFilter: a legal predicate that rejects a single pair
// variant suppresses only that variant's combos, not the whole pair. Sun Kiss (Yellow) gets
// rejected; the remaining 3 × 2 = 6 cross-products still emit — matches how
// cardSwapMutations treats per-printing legality.
func TestCardPairMutations_RespectsLegalFilter(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	legal := func(c card.Card) bool { return c.ID() != card.SunKissYellow }
	muts := cardPairMutations(d, 2, legal)
	for i, m := range muts {
		if strings.Contains(m.Description, "Sun Kiss (Yellow)") {
			t.Errorf("mutation %d (%s): added rejected Sun Kiss (Yellow)", i, m.Description)
		}
	}
	// 9 - 3 (Sun Kiss Yellow combos with the 3 Moon Wish variants) = 6 remaining cross-products.
	const wantCombos = 6
	if len(muts) != wantCombos {
		t.Errorf("got %d mutations after rejecting Sun Kiss (Yellow), want %d", len(muts), wantCombos)
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
