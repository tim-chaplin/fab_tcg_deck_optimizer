package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// TestCardPairMutations_EnumeratesAllVariantCrossProducts: with neither pair half present,
// the generator emits a candidate per (firstVariant, secondVariant) cross-product per
// distinct removed-ID combo from the deck. For a [a, a, b, b] deck the unique removed-ID
// combos are {(a,a), (a,b), (b,b)} = 3; with 9 cross-products per implemented pair that's
// 3 × 9 = 27 mutations per implemented pair.
//
// Pairs whose halves carry NotImplemented don't contribute — PairAddAllowed gates them
// out — so the expected total scales with the count of fully-implemented pairs, not the
// raw len(CardPairs).
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
// the combos PairAddAllowed lets through. Tests that compute expected mutation counts
// use this so future churn (a card dropping its NotImplemented marker, a card moving
// between subpackages) doesn't silently make the assertion stale.
func countImplementedPairCombos() int {
	n := 0
	for _, p := range CardPairs {
		n += countImplementedInGroup(p.First) * countImplementedInGroup(p.Second)
	}
	return n
}

// countImplementedInGroup returns how many variants in g are pool-eligible: registered
// (GetCard returns non-nil) and free of the NotImplemented marker. Unregistered IDs
// belong to cards in internal/cards/notimplemented/ or internal/cards/unplayable/, both
// of which PairAddAllowed rejects, so they don't count.
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

// TestCardPairMutations_RemovesBothCopiesOfDuplicate is the pilot for the index-based
// generator: a 2-card deck of [HocusPocusBlue, HocusPocusBlue] must yield mutations that
// remove BOTH copies of HocusPocusBlue and add a Moon Wish / Sun Kiss pair. A unique-ID
// generator would skip this (only one unique ID, no inter-ID pair exists); index-based
// iteration over (0, 1) reaches it directly.
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

// TestCardPairMutations_FiresWhenOneHalfAlreadyPresent: pair mutations fire whenever the
// removal-pair / add-pair don't overlap on a card ID. With one half partially present, the
// climber can still grow the OTHER variant of that side as a pair-shape mutation rather than
// two sequential single-slot swaps.
//
// Same-ID overlap suppression (e.g. -1 SunKissRed + +1 SunKissRed reducing to a single-slot)
// is the orthogonal optimisation tested in
// TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps; here we check that
// non-overlapping variant combinations still emit despite Sun Kiss [R] being a removal
// candidate.
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

// TestCardPairMutations_GeneratesCapViolatingCandidates pins the cap-blind contract:
// PairSwapMutations enumerates every (i, j) × (firstVariant, secondVariant) tuple that
// survives overlap suppression — even ones whose result deck would violate maxCopies.
// FilterMaxCopiesViolations is the gate that strips violators downstream.
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

// TestCardPairMutations_HandlesUnbalancedHalfCounts: the generator should work with arbitrary
// per-variant counts of each half. Drives this with a deck holding 5 Moon Wish (across
// variants) and 3 Sun Kiss (across variants) — a realistic mid-climb state. cap-blind
// enumeration emits candidates regardless of saturation; the resulting decks must remain at
// the original card count.
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

// TestCardPairMutations_OverlapSuppressionSkipsRedundantSwaps: when a removal target is
// itself a pair member, the resulting mutation reduces to a single-slot swap (the matching
// pair member's count is unchanged after -1 +1). Single-slot already covers that, so the
// pair generator skips those combos. Drives this with a deck containing Sun Kiss [R] as a
// removal candidate and verifies no mutation removes and re-adds Sun Kiss [R].
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

// TestCardPairMutations_SkipsNotImplementedHalves: pair mutations never name a card that
// carries NotImplemented as one of the +1 adds. CardPairs registers pairings whose
// halves aren't all modelled yet (e.g. Belittle / Minnowism, Amulet of Havencall / Rally
// the Rearguard); those entries shouldn't leak NotImplemented printings into the search
// pool just because they're listed.
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

// TestCardPairMutations_RespectsLegalFilter: a legal predicate that rejects a single pair
// variant suppresses only that variant's combos, not the whole pair. Sun Kiss [Y] gets
// rejected; the remaining 3 × 2 = 6 cross-products still emit per unique removal combo —
// matches how singleSwapMutations treats per-printing legality.
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

// TestCardPairMutations_DeterministicOrdering: two back-to-back calls must produce the same
// mutation sequence. AllMutations consumers (the iterate-mode worker pool) rely on stable
// indexing for reproducibility under a fixed seed.
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

// TestFilterMaxCopiesViolations_StripsCapViolators: the post-filter must drop any mutation
// whose result deck holds more than maxCopies of a card. Built two synthetic mutations: one
// clean (deck 4 cards, all distinct) and one violator (5 copies of Moon Wish [R] at
// maxCopies=2). Filter keeps the clean one, drops the violator.
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
