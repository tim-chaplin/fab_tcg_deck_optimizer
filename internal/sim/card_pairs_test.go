package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests that PairSwapMutations enumerates every (firstVariant, secondVariant) cross-product
// per distinct removed-ID combo, gated to implemented pairs.
func TestCardPairMutations_EnumeratesAllVariantCrossProducts(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	muts := PairSwapMutations(d, nil)
	const dedupedRemovalCombos = 3 // (a,a), (a,b), (b,b)
	implementedCombos := countImplementedPairCombos()
	want := implementedCombos * dedupedRemovalCombos
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (%d implemented variant combos × %d removal combos)",
			len(muts), want, implementedCombos, dedupedRemovalCombos)
	}

	// Every (firstID, secondID) cross-product from CardPairs[0] must appear at least once.
	type combo struct{ first, second ids.CardID }
	seen := map[combo]bool{}
	for _, m := range muts {
		for _, fID := range CardPairs[0].First {
			for _, sID := range CardPairs[0].Second {
				if strings.Contains(m.Description, "+1 "+GetCard(fID).Name()) &&
					strings.Contains(m.Description, "+1 "+GetCard(sID).Name()) {
					seen[combo{fID, sID}] = true
				}
			}
		}
	}
	wantCombos := len(CardPairs[0].First) * len(CardPairs[0].Second)
	if len(seen) != wantCombos {
		t.Errorf("variant cross-product coverage: saw %d distinct (first, second) pairs, want %d",
			len(seen), wantCombos)
	}
}

// countImplementedPairCombos returns the total number of (firstVariant, secondVariant)
// cross-product entries across CardPairs whose both halves are pool-eligible — exactly
// the combos PairAddAllowed lets through.
func countImplementedPairCombos() int {
	n := 0
	for _, p := range CardPairs {
		n += countImplementedInGroup(p.First) * countImplementedInGroup(p.Second)
	}
	return n
}

// countImplementedInGroup returns how many variants in g are pool-eligible — registered and
// free of the NotImplemented marker. Cards in notimplemented/ or unplayable/ are unregistered
// and don't count.
func countImplementedInGroup(g CardGroup) int {
	n := 0
	for _, id := range g {
		c := GetCard(id)
		if c == nil {
			continue
		}
		if _, unimplemented := c.(NotImplemented); !unimplemented {
			n++
		}
	}
	return n
}

// Tests that PairSwapMutations on a duplicate-ID deck removes BOTH copies of the duplicate
// and adds the new pair (index-based iteration, not unique-ID).
func TestCardPairMutations_RemovesBothCopiesOfDuplicate(t *testing.T) {
	hp := GetCard(ids.HocusPocusBlue)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{hp, hp})

	muts := PairSwapMutations(d, nil)
	// Exactly one removed-ID combo (HocusPocusBlue, HocusPocusBlue) × 9 variant combos.
	const want = 9
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (1 removal combo × 9 variant combos)",
			len(muts), want)
	}

	for i, m := range muts {
		if !strings.Contains(m.Description, "-1 Hocus Pocus [B], -1 Hocus Pocus [B]") {
			t.Errorf("mutation %d (%s): expected both copies of Hocus Pocus [B] removed",
				i, m.Description)
		}
		// Result deck has 2 cards (the new pair), zero HocusPocusBlue.
		if len(m.Deck.Cards) != 2 {
			t.Errorf("mutation %d (%s): card count %d, want 2", i, m.Description, len(m.Deck.Cards))
		}
		for _, c := range m.Deck.Cards {
			if c.ID() == ids.HocusPocusBlue {
				t.Errorf("mutation %d (%s): result deck still holds Hocus Pocus [B]",
					i, m.Description)
			}
		}
	}
}

// Tests that pair mutations still emit when one pair half is already present, as long as
// the add side picks a non-overlapping variant.
func TestCardPairMutations_FiresWhenOneHalfAlreadyPresent(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	sk := GetCard(ids.SunKissRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, a, sk})

	muts := PairSwapMutations(d, nil)
	if len(muts) == 0 {
		t.Fatal("expected pair mutations even with one half present")
	}
	sawDifferentSunKissVariantAdd := false
	for _, m := range muts {
		if strings.Contains(m.Description, "+1 Sun Kiss [Y]") ||
			strings.Contains(m.Description, "+1 Sun Kiss [B]") {
			sawDifferentSunKissVariantAdd = true
			break
		}
	}
	if !sawDifferentSunKissVariantAdd {
		t.Error("expected at least one mutation adding a non-Red Sun Kiss variant " +
			"when Red is already present")
	}
}

// Tests that PairSwapMutations is cap-blind — emits maxCopies-violating candidates that
// FilterMaxCopiesViolations strips downstream.
func TestCardPairMutations_GeneratesCapViolatingCandidates(t *testing.T) {
	skR := GetCard(ids.SunKissRed)
	a := GetCard(ids.ArcanicCrackleRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}},
		[]Card{skR, skR, a, a})

	// 3 unique removed-ID combos after dedupe: (skR, skR), (skR, a), (a, a). Overlap
	// suppression skips a combo when one of its removed IDs equals one of the add IDs;
	// SunKissRed is an add ID for 3 of the 9 (Moon Wish, Sun Kiss) cross-products. So:
	//   (skR, skR) and (skR, a) each emit 9 - 3 = 6 surviving combos.
	//   (a, a) emits all 9.
	// Total = 6 + 6 + 9 = 21.
	muts := PairSwapMutations(d, nil)
	const want = 21
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (cap-blind enumeration)", len(muts), want)
	}

	// At least one of those mutations must add Sun Kiss [R] again — pushing the count to 3
	// — which would violate maxCopies=2. PairSwapMutations does NOT enforce that; the post-
	// filter in AllMutations does.
	sawCapViolator := false
	for _, m := range muts {
		if strings.Contains(m.Description, "+1 Sun Kiss [R]") {
			counts := map[ids.CardID]int{}
			for _, c := range m.Deck.Cards {
				counts[c.ID()]++
			}
			if counts[ids.SunKissRed] > 2 {
				sawCapViolator = true
				break
			}
		}
	}
	if !sawCapViolator {
		t.Error("expected at least one cap-violating candidate from PairSwapMutations " +
			"(FilterMaxCopiesViolations is the responsible gate)")
	}
}

// Tests that PairSwapMutations works with unbalanced per-variant counts and the result
// decks stay at the original card count.
func TestCardPairMutations_HandlesUnbalancedHalfCounts(t *testing.T) {
	mwR := GetCard(ids.MoonWishRed)
	mwY := GetCard(ids.MoonWishYellow)
	mwB := GetCard(ids.MoonWishBlue)
	skR := GetCard(ids.SunKissRed)
	skY := GetCard(ids.SunKissYellow)
	skB := GetCard(ids.SunKissBlue)
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	cardsList := []Card{
		mwR, mwR, mwY, mwY, mwB,
		skR, skY, skB,
		a, a, a, b, b, b,
	}
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, cardsList)

	muts := PairSwapMutations(d, nil)
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
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})
	srcKey := CardMultisetKey(d.Cards)
	for i, m := range PairSwapMutations(d, nil) {
		if CardMultisetKey(m.Deck.Cards) == srcKey {
			t.Errorf("mutation %d (%s) produced a no-op (same multiset as source)", i, m.Description)
		}
	}
}

// Tests that PairSwapMutations suppresses combos that reduce to a single-slot swap (a -1/+1
// of the same card), which the single-slot generator already covers.
func TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps(t *testing.T) {
	skR := GetCard(ids.SunKissRed)
	a := GetCard(ids.ArcanicCrackleRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{skR, a, a, a})
	for i, m := range PairSwapMutations(d, nil) {
		if strings.Contains(m.Description, "-1 Sun Kiss [R]") &&
			strings.Contains(m.Description, "+1 Sun Kiss [R]") {
			t.Errorf("mutation %d (%s): redundant -1/+1 of Sun Kiss [R] — overlap suppression failed",
				i, m.Description)
		}
	}
}

// Tests that PairSwapMutations never adds a NotImplemented card, even when CardPairs lists
// pairings whose halves aren't all modelled.
func TestCardPairMutations_SkipsNotImplementedHalves(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	for i, m := range PairSwapMutations(d, nil) {
		for _, c := range m.Deck.Cards {
			if _, unimplemented := c.(NotImplemented); unimplemented {
				t.Errorf("mutation %d (%s) introduced NotImplemented card %s",
					i, m.Description, c.Name())
			}
		}
	}
}

// Tests that a legal predicate rejecting one pair variant suppresses only that variant's
// combos, not the whole pair.
func TestCardPairMutations_RespectsLegalFilter(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	legal := func(c Card) bool { return c.ID() != ids.SunKissYellow }
	muts := PairSwapMutations(d, legal)
	for i, m := range muts {
		if strings.Contains(m.Description, "Sun Kiss [Y]") {
			t.Errorf("mutation %d (%s): added rejected Sun Kiss [Y]", i, m.Description)
		}
	}
	// 3 unique removal combos × 6 surviving cross-products = 18.
	const want = 18
	if len(muts) != want {
		t.Errorf("got %d mutations after rejecting Sun Kiss [Y], want %d", len(muts), want)
	}
}

// Tests that PairSwapMutations is deterministic: two back-to-back calls produce the same
// mutation sequence.
func TestCardPairMutations_DeterministicOrdering(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	first := PairSwapMutations(d, nil)
	second := PairSwapMutations(d, nil)
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

// Tests that FilterMaxCopiesViolations drops any mutation whose result deck exceeds maxCopies
// of any card.
func TestFilterMaxCopiesViolations_StripsCapViolators(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	mw := GetCard(ids.MoonWishRed)
	clean := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, b, mw, mw})
	violator := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}},
		[]Card{mw, mw, mw, mw, mw})

	muts := []Mutation{
		{Deck: clean, Description: "clean"},
		{Deck: violator, Description: "violator"},
	}
	out := FilterMaxCopiesViolations(muts, 2)
	if len(out) != 1 {
		t.Fatalf("got %d mutations after filter, want 1 (only the clean one survives)", len(out))
	}
	if out[0].Description != "clean" {
		t.Errorf("survivor description = %q, want %q", out[0].Description, "clean")
	}
}

// TestRespectsMaxCopies_ShortCircuits: RespectsMaxCopies returns false immediately when a
// count exceeds the cap, without scanning the full slice. Sentinel for the inner-loop fast
// path in FilterMaxCopiesViolations.
func TestRespectsMaxCopies_ShortCircuits(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	cs := []Card{a, a, a}
	if RespectsMaxCopies(cs, 2) {
		t.Error("3 copies at maxCopies=2 should fail RespectsMaxCopies")
	}
	if !RespectsMaxCopies(cs, 3) {
		t.Error("3 copies at maxCopies=3 should pass RespectsMaxCopies")
	}
}
