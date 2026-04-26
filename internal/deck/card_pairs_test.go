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
// the generator emits a candidate per (firstVariant, secondVariant) cross-product per
// distinct removed-ID combo from the deck. For a [a, a, b, b] deck the unique removed-ID
// combos are {(a,a), (a,b), (b,b)} = 3; with 9 (Moon Wish, Sun Kiss) cross-products that's
// 3 × 9 = 27 mutations.
func TestCardPairMutations_EnumeratesAllVariantCrossProducts(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	muts := cardPairMutations(d, nil)
	const variantCombosPerPair = 3 * 3
	const dedupedRemovalCombos = 3 // (a,a), (a,b), (b,b)
	want := len(cardPairs) * variantCombosPerPair * dedupedRemovalCombos
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (%d pairs × %d variant combos × %d removal combos)",
			len(muts), want, len(cardPairs), variantCombosPerPair, dedupedRemovalCombos)
	}

	// Every (firstID, secondID) cross-product from cardPairs[0] must appear at least once.
	type combo struct{ first, second card.ID }
	seen := map[combo]bool{}
	for _, m := range muts {
		for _, fID := range cardPairs[0].First {
			for _, sID := range cardPairs[0].Second {
				if strings.Contains(m.Description, "+1 "+cards.Get(fID).Name()) &&
					strings.Contains(m.Description, "+1 "+cards.Get(sID).Name()) {
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

// TestCardPairMutations_RemovesBothCopiesOfDuplicate is the pilot for the index-based
// generator: a 2-card deck of [HocusPocusBlue, HocusPocusBlue] must yield mutations that
// remove BOTH copies of HocusPocusBlue and add a Moon Wish / Sun Kiss pair. A unique-ID
// generator would skip this (only one unique ID, no inter-ID pair exists); index-based
// iteration over (0, 1) reaches it directly.
func TestCardPairMutations_RemovesBothCopiesOfDuplicate(t *testing.T) {
	hp := cards.Get(card.HocusPocusBlue)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{hp, hp})

	muts := cardPairMutations(d, nil)
	// Exactly one removed-ID combo (HocusPocusBlue, HocusPocusBlue) × 9 variant combos.
	const want = 9
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (1 removal combo × 9 variant combos)",
			len(muts), want)
	}

	for i, m := range muts {
		if !strings.Contains(m.Description, "-1 Hocus Pocus (Blue), -1 Hocus Pocus (Blue)") {
			t.Errorf("mutation %d (%s): expected both copies of Hocus Pocus (Blue) removed",
				i, m.Description)
		}
		// Result deck has 2 cards (the new pair), zero HocusPocusBlue.
		if len(m.Deck.Cards) != 2 {
			t.Errorf("mutation %d (%s): card count %d, want 2", i, m.Description, len(m.Deck.Cards))
		}
		for _, c := range m.Deck.Cards {
			if c.ID() == card.HocusPocusBlue {
				t.Errorf("mutation %d (%s): result deck still holds Hocus Pocus (Blue)",
					i, m.Description)
			}
		}
	}
}

// TestCardPairMutations_FiresWhenOneHalfAlreadyPresent: pair mutations fire whenever the
// removal-pair / add-pair don't overlap on a card ID. With one half partially present, the
// climber can still grow the OTHER variant of that side as a pair-shape mutation rather than
// two sequential single-slot swaps.
//
// Same-ID overlap suppression (e.g. -1 SunKissRed + +1 SunKissRed reducing to a single-slot)
// is the orthogonal optimisation tested in
// TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps; here we check that
// non-overlapping variant combinations still emit despite Sun Kiss (Red) being a removal
// candidate.
func TestCardPairMutations_FiresWhenOneHalfAlreadyPresent(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	sk := cards.Get(card.SunKissRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, sk})

	muts := cardPairMutations(d, nil)
	if len(muts) == 0 {
		t.Fatal("expected pair mutations even with one half present")
	}
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

// TestCardPairMutations_GeneratesCapViolatingCandidates documents the contract change:
// cardPairMutations is cap-blind, so a deck holding the maximum copies of a pair variant
// can still see candidates that would push that variant past maxCopies. The shared
// filterMaxCopiesViolations post-pass strips them when AllMutations runs; the per-generator
// cap check is gone.
func TestCardPairMutations_GeneratesCapViolatingCandidates(t *testing.T) {
	skR := cards.Get(card.SunKissRed)
	a := cards.Get(card.ArcanicCrackleRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
		[]card.Card{skR, skR, a, a})

	// The pair generator sees 4 indices with removed-ID combos (skR,skR), (skR,a), (a,a) —
	// 3 unique combos. (skR,*) overlaps any combo whose secondID is SunKissRed (3 combos,
	// suppressed); (a,a) and (skR,skR) trigger no overlap. So combos that DO emit:
	//   (skR, skR) × 9 = 9 (all emit; overlap only fires on add-side IDs)
	//   wait — overlap fires when removeID == addID. SunKissRed IS an add ID for 3 combos.
	//   So (skR, skR) skipped for those 3 combos: 9 - 3 = 6 emit.
	//   (skR, a) skipped for 3 combos: 9 - 3 = 6 emit.
	//   (a, a) all 9 emit.
	// Total: 6 + 6 + 9 = 21.
	muts := cardPairMutations(d, nil)
	const want = 21
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (cap-blind enumeration)", len(muts), want)
	}

	// At least one of those mutations must add Sun Kiss (Red) again — pushing the count to 3
	// — which would violate maxCopies=2. cardPairMutations does NOT enforce that; the post-
	// filter in AllMutations does.
	sawCapViolator := false
	for _, m := range muts {
		if strings.Contains(m.Description, "+1 Sun Kiss (Red)") {
			counts := map[card.ID]int{}
			for _, c := range m.Deck.Cards {
				counts[c.ID()]++
			}
			if counts[card.SunKissRed] > 2 {
				sawCapViolator = true
				break
			}
		}
	}
	if !sawCapViolator {
		t.Error("expected at least one cap-violating candidate from cardPairMutations " +
			"(filterMaxCopiesViolations is the responsible gate)")
	}
}

// TestCardPairMutations_HandlesUnbalancedHalfCounts: the generator should work with arbitrary
// per-variant counts of each half. Drives this with a deck holding 5 Moon Wish (across
// variants) and 3 Sun Kiss (across variants) — a realistic mid-climb state. cap-blind
// enumeration emits candidates regardless of saturation; the resulting decks must remain at
// the original card count.
func TestCardPairMutations_HandlesUnbalancedHalfCounts(t *testing.T) {
	mwR := cards.Get(card.MoonWishRed)
	mwY := cards.Get(card.MoonWishYellow)
	mwB := cards.Get(card.MoonWishBlue)
	skR := cards.Get(card.SunKissRed)
	skY := cards.Get(card.SunKissYellow)
	skB := cards.Get(card.SunKissBlue)
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	cardsList := []card.Card{
		mwR, mwR, mwY, mwY, mwB,
		skR, skY, skB,
		a, a, a, b, b, b,
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, cardsList)

	muts := cardPairMutations(d, nil)
	if len(muts) == 0 {
		t.Fatal("expected pair mutations on unbalanced deck")
	}
	for i, m := range muts {
		if len(m.Deck.Cards) != len(cardsList) {
			t.Errorf("mutation %d (%s): card count %d, want %d (size must stay stable)",
				i, m.Description, len(m.Deck.Cards), len(cardsList))
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
	for i, m := range cardPairMutations(d, nil) {
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
	for i, m := range cardPairMutations(d, nil) {
		if strings.Contains(m.Description, "-1 Sun Kiss (Red)") &&
			strings.Contains(m.Description, "+1 Sun Kiss (Red)") {
			t.Errorf("mutation %d (%s): redundant -1/+1 of Sun Kiss (Red) — overlap suppression failed",
				i, m.Description)
		}
	}
}

// TestCardPairMutations_RespectsLegalFilter: a legal predicate that rejects a single pair
// variant suppresses only that variant's combos, not the whole pair. Sun Kiss (Yellow) gets
// rejected; the remaining 3 × 2 = 6 cross-products still emit per unique removal combo —
// matches how cardSwapMutations treats per-printing legality.
func TestCardPairMutations_RespectsLegalFilter(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	legal := func(c card.Card) bool { return c.ID() != card.SunKissYellow }
	muts := cardPairMutations(d, legal)
	for i, m := range muts {
		if strings.Contains(m.Description, "Sun Kiss (Yellow)") {
			t.Errorf("mutation %d (%s): added rejected Sun Kiss (Yellow)", i, m.Description)
		}
	}
	// 3 unique removal combos × 6 surviving cross-products = 18.
	const want = 18
	if len(muts) != want {
		t.Errorf("got %d mutations after rejecting Sun Kiss (Yellow), want %d", len(muts), want)
	}
}

// TestCardPairMutations_DeterministicOrdering: two back-to-back calls must produce the same
// mutation sequence. AllMutations consumers (the iterate-mode worker pool) rely on stable
// indexing for reproducibility under a fixed seed.
func TestCardPairMutations_DeterministicOrdering(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := cardPairMutations(d, nil)
	second := cardPairMutations(d, nil)
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

// TestFilterMaxCopiesViolations_StripsCapViolators: the post-filter must drop any mutation
// whose result deck holds more than maxCopies of a card. Built two synthetic mutations: one
// clean (deck 4 cards, all distinct) and one violator (5 copies of Moon Wish (Red) at
// maxCopies=2). Filter keeps the clean one, drops the violator.
func TestFilterMaxCopiesViolations_StripsCapViolators(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	b := cards.Get(card.ArcanicSpikeRed)
	mw := cards.Get(card.MoonWishRed)
	clean := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, b, mw, mw})
	violator := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
		[]card.Card{mw, mw, mw, mw, mw})

	muts := []Mutation{
		{Deck: clean, Description: "clean"},
		{Deck: violator, Description: "violator"},
	}
	out := filterMaxCopiesViolations(muts, 2)
	if len(out) != 1 {
		t.Fatalf("got %d mutations after filter, want 1 (only the clean one survives)", len(out))
	}
	if out[0].Description != "clean" {
		t.Errorf("survivor description = %q, want %q", out[0].Description, "clean")
	}
}

// TestRespectsMaxCopies_ShortCircuits: respectsMaxCopies returns false immediately when a
// count exceeds the cap, without scanning the full slice. Sentinel for the inner-loop fast
// path in filterMaxCopiesViolations.
func TestRespectsMaxCopies_ShortCircuits(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	cs := []card.Card{a, a, a}
	if respectsMaxCopies(cs, 2) {
		t.Error("3 copies at maxCopies=2 should fail respectsMaxCopies")
	}
	if !respectsMaxCopies(cs, 3) {
		t.Error("3 copies at maxCopies=3 should pass respectsMaxCopies")
	}
}
