package deck

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Mutation-enumeration tests pinned against the package's hand-rolled fakes (fakeCard,
// fakeWeapon, fakeRegistry — declared in deck_test.go) plus the fake CardPair fixture
// below, so this file's dependencies stay scoped to internal/deck and registry/ids. No
// production registry, cards, weapons, heroes, or sim package. Every algorithmic property
// the production enumerator promises is exercised through these fakes; the production
// CardPairs registry is only relied on by AllMutations behind the scenes (and the "no
// fake card collides with a real pair half" fixture invariant ensures AllMutations
// contributes zero pair mutations to the count assertions below).

// fakeHero is the Hero (any) stub used by every test deck. A typed empty struct
// keeps interface equality stable across mutated decks.
type fakeHero struct{}

// Fake card IDs are picked near the top of the uint16 range so AllMutations' iteration
// over the real CardPairs registry never finds a fake ID in any pair half — the pair
// generator emits zero mutations for these decks, and the AllMutations count assertions
// don't have to predict the production pair contribution.
const (
	fakeC1 ids.CardID = 60001
	fakeC2 ids.CardID = 60002
	fakeC3 ids.CardID = 60003
	fakeC4 ids.CardID = 60004
	fakeC5 ids.CardID = 60005

	// Stand-in IDs for the two halves of a fake CardPair. Three "moon-wish" variants and
	// three "sun-kiss" variants give a 3 × 3 = 9 cross-product to exercise pair enumeration
	// without touching the production CardPairs registry.
	fakeMW1 ids.CardID = 61001
	fakeMW2 ids.CardID = 61002
	fakeMW3 ids.CardID = 61003
	fakeSK1 ids.CardID = 61011
	fakeSK2 ids.CardID = 61012
	fakeSK3 ids.CardID = 61013
)

// makeFakeCard builds a fakeCard whose DisplayName follows a "card-N" template so
// mutation descriptions stay readable in test failure messages.
func makeFakeCard(id ids.CardID) fakeCard {
	return fakeCard{id: id, display: fmt.Sprintf("card-%d", id)}
}

// allMutationsFixture is the shared roster for AllMutations tests: 5 distinct fake cards
// + 2 1H weapons. Built once per test via this helper so the (loadouts, pool) sizes stay
// obvious in the wantCount arithmetic.
func allMutationsFixture() *fakeRegistry {
	return &fakeRegistry{
		cards: []Card{
			makeFakeCard(fakeC1),
			makeFakeCard(fakeC2),
			makeFakeCard(fakeC3),
			makeFakeCard(fakeC4),
			makeFakeCard(fakeC5),
		},
		weapons: []Weapon{
			fakeWeapon{name: "weapon-A", hands: 1},
			fakeWeapon{name: "weapon-B", hands: 1},
		},
	}
}

// pairFixture is the shared roster for pair-mutation tests: every fake pair half plus a
// neutral "filler" card so deck-side cards aren't accidentally pair-eligible.
func pairFixture() *fakeRegistry {
	return &fakeRegistry{
		cards: []Card{
			makeFakeCard(fakeC1),
			makeFakeCard(fakeC2),
			makeFakeCard(fakeMW1),
			makeFakeCard(fakeMW2),
			makeFakeCard(fakeMW3),
			makeFakeCard(fakeSK1),
			makeFakeCard(fakeSK2),
			makeFakeCard(fakeSK3),
		},
		weapons: []Weapon{
			fakeWeapon{name: "weapon-A", hands: 1},
		},
	}
}

// fakeMoonSunPair is the canonical (3-variant × 3-variant) test pair driving every
// PairSwapMutations test below. Enumeration ordering is group-slice order, so the IDs
// here directly determine the expected mutation sequence.
var fakeMoonSunPair = []CardPair{
	{
		First:  CardGroup{fakeMW1, fakeMW2, fakeMW3},
		Second: CardGroup{fakeSK1, fakeSK2, fakeSK3},
	},
}

// TestAllMutations_CountsAndShape pins the total mutation count against a deck of
// (weapons=[A, B], cards=[c1, c1, c2, c2]) over a 2-weapon / 5-card roster:
//   - WeaponLoadouts emits {(A,A), (A,B), (B,B)} = 3 loadouts; the deck holds (A, B), so
//     2 alternative loadouts remain.
//   - singleSwapMutations emits 2 unique removals × (5 pool − self − other-at-cap) =
//     2 × 3 = 6 candidates at maxCopies=2.
//   - pairSwapMutations emits 0 because no fake-card ID appears in CardPairs.
//
// Total: 8 mutations. Every result deck keeps its 4-card shape and shared hero, with a
// non-empty description.
func TestAllMutations_CountsAndShape(t *testing.T) {
	reg := allMutationsFixture()
	d := New(fakeHero{}, []Weapon{reg.weapons[0], reg.weapons[1]},
		[]Card{reg.cards[0], reg.cards[0], reg.cards[1], reg.cards[1]})

	muts := AllMutations(d, 2, reg)

	loadouts := weaponLoadouts(reg.LegalWeapons())
	wantWeaponMuts := len(loadouts) - 1
	wantCardMuts := 2 * (len(reg.LegalCards()) - 2)
	want := wantWeaponMuts + wantCardMuts

	if len(muts) != want {
		t.Fatalf("got %d mutations, want %d (%d weapon + %d card; pair=0 by fixture design)",
			len(muts), want, wantWeaponMuts, wantCardMuts)
	}
	for i, m := range muts {
		if len(m.Deck.cards) != 4 {
			t.Errorf("mutation %d: card count %d, want 4", i, len(m.Deck.cards))
		}
		if m.Deck.Hero != d.Hero {
			t.Errorf("mutation %d: hero changed", i)
		}
		if m.Description == "" {
			t.Errorf("mutation %d: empty description", i)
		}
	}
}

// Tests single-card-swap semantics: odd-count distributions are legal and raising
// maxCopies opens adds to in-deck cards below the cap.
func TestAllMutations_OddCountsAllowed(t *testing.T) {
	reg := allMutationsFixture()
	d := New(fakeHero{}, []Weapon{reg.weapons[0], reg.weapons[1]},
		[]Card{reg.cards[0], reg.cards[0], reg.cards[1], reg.cards[1]})

	// At maxCopies=3, each of the 2 in-deck cards (c1, c2) is below the cap, so
	// "remove c1, add c2" (and the mirror) become legal. That's 2 more card mutations than
	// the maxCopies=2 case.
	mutsLow := AllMutations(d, 2, reg)
	mutsHigh := AllMutations(d, 3, reg)
	if len(mutsHigh)-len(mutsLow) != 2 {
		t.Errorf("maxCopies=3 should produce exactly 2 more mutations than maxCopies=2; got diff=%d",
			len(mutsHigh)-len(mutsLow))
	}

	// Every mutation at maxCopies=3 still has exactly 4 cards; some should have an odd
	// count of one card — single-card swaps can leave odd-count slots when maxCopies
	// allows it.
	sawOdd := false
	for _, m := range mutsHigh {
		if len(m.Deck.cards) != 4 {
			t.Errorf("card count %d, want 4", len(m.Deck.cards))
		}
		counts := map[ids.CardID]int{}
		for _, c := range m.Deck.cards {
			counts[c.ID()]++
		}
		for _, n := range counts {
			if n%2 == 1 {
				sawOdd = true
			}
			if n > 3 {
				t.Errorf("card count %d exceeds maxCopies=3: %v", n, m.Description)
			}
		}
	}
	if !sawOdd {
		t.Errorf("expected at least one mutation with an odd-count card; single-card swaps always produce odd counts from a 2/2 starting deck")
	}
}

// TestAllMutations_PreservesSideboard pins that every derived mutation inherits the
// source deck's Sideboard verbatim. Without this guarantee an anneal round would silently
// drop the user's hand-managed sideboard as soon as it accepted a mutation and wrote the
// deck back.
func TestAllMutations_PreservesSideboard(t *testing.T) {
	reg := allMutationsFixture()
	d := New(fakeHero{}, []Weapon{reg.weapons[0], reg.weapons[1]},
		[]Card{reg.cards[0], reg.cards[0], reg.cards[1], reg.cards[1]})
	d.Sideboard = []string{"sb-x", "sb-y", "sb-y"}

	muts := AllMutations(d, 2, reg)
	if len(muts) == 0 {
		t.Fatal("expected at least one mutation")
	}

	wantCounts := map[string]int{"sb-x": 1, "sb-y": 2}
	for i, m := range muts {
		got := map[string]int{}
		for _, name := range m.Deck.Sideboard {
			got[name]++
		}
		for name, want := range wantCounts {
			if got[name] != want {
				t.Errorf("mutation %d (%s): sideboard count for %s = %d, want %d",
					i, m.Description, name, got[name], want)
				break
			}
		}
	}
}

// Two back-to-back AllMutations calls on the same deck produce identical mutation
// sequences — same length, same per-index weapon loadout / description / card-ID order.
func TestAllMutations_Deterministic(t *testing.T) {
	reg := allMutationsFixture()
	d := New(fakeHero{}, []Weapon{reg.weapons[0], reg.weapons[1]},
		[]Card{reg.cards[0], reg.cards[0], reg.cards[1], reg.cards[1]})

	first := AllMutations(d, 2, reg)
	second := AllMutations(d, 2, reg)

	if len(first) != len(second) {
		t.Fatalf("mutation counts differ between calls: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if weaponKey(first[i].Deck.Weapons) != weaponKey(second[i].Deck.Weapons) {
			t.Errorf("mutation %d weapons differ between calls", i)
		}
		if first[i].Description != second[i].Description {
			t.Errorf("mutation %d descriptions differ: %q vs %q",
				i, first[i].Description, second[i].Description)
		}
		for j, c := range first[i].Deck.cards {
			if c.ID() != second[i].Deck.cards[j].ID() {
				t.Errorf("mutation %d card[%d] differs between calls: %v vs %v",
					i, j, c.ID(), second[i].Deck.cards[j].ID())
			}
		}
	}
}

// Every emitted mutation has a different Fingerprint than the source deck. Defensive
// against an enumerator regression that lets a no-op (same multiset, same loadout) slip
// through.
func TestAllMutations_NoDuplicateOfSource(t *testing.T) {
	reg := allMutationsFixture()
	d := New(fakeHero{}, []Weapon{reg.weapons[0], reg.weapons[1]},
		[]Card{reg.cards[0], reg.cards[0], reg.cards[0], reg.cards[0]})
	srcKey := d.Fingerprint()
	for i, m := range AllMutations(d, 2, reg) {
		if m.Deck.Fingerprint() == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
		}
	}
}

// Tests that PairSwapMutations enumerates every (firstVariant, secondVariant)
// cross-product per distinct removed-ID combo.
func TestPairSwapMutations_EnumeratesAllVariantCrossProducts(t *testing.T) {
	reg := pairFixture()
	a, b := makeFakeCard(fakeC1), makeFakeCard(fakeC2)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{a, a, b, b})

	muts := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
	const dedupedRemovalCombos = 3 // (a,a), (a,b), (b,b)
	wantCombos := len(fakeMoonSunPair[0].First) * len(fakeMoonSunPair[0].Second)
	want := wantCombos * dedupedRemovalCombos
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (%d cross-products × %d removal combos)",
			len(muts), want, wantCombos, dedupedRemovalCombos)
	}

	// Every (firstID, secondID) cross-product from fakeMoonSunPair[0] must appear at
	// least once.
	type combo struct{ first, second ids.CardID }
	seen := map[combo]bool{}
	for _, m := range muts {
		for _, fID := range fakeMoonSunPair[0].First {
			for _, sID := range fakeMoonSunPair[0].Second {
				if strings.Contains(m.Description, "+1 "+makeFakeCard(fID).DisplayName()) &&
					strings.Contains(m.Description, "+1 "+makeFakeCard(sID).DisplayName()) {
					seen[combo{fID, sID}] = true
				}
			}
		}
	}
	if len(seen) != wantCombos {
		t.Errorf("variant cross-product coverage: saw %d distinct (first, second) pairs, want %d",
			len(seen), wantCombos)
	}
}

// Tests that PairSwapMutations on a duplicate-ID deck removes BOTH copies of the
// duplicate and adds the new pair (index-based iteration, not unique-ID).
func TestPairSwapMutations_RemovesBothCopiesOfDuplicate(t *testing.T) {
	reg := pairFixture()
	dup := makeFakeCard(fakeC1)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]}, []Card{dup, dup})

	muts := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
	// Exactly one removed-ID combo (c1, c1) × (3 × 3) = 9 variant combos.
	const want = 9
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (1 removal combo × 9 variant combos)",
			len(muts), want)
	}

	for i, m := range muts {
		want := fmt.Sprintf("-1 %s, -1 %s", dup.DisplayName(), dup.DisplayName())
		if !strings.Contains(m.Description, want) {
			t.Errorf("mutation %d (%s): expected both copies of %s removed",
				i, m.Description, dup.DisplayName())
		}
		// Result deck has 2 cards (the new pair), zero c1.
		if len(m.Deck.cards) != 2 {
			t.Errorf("mutation %d (%s): card count %d, want 2", i, m.Description, len(m.Deck.cards))
		}
		for _, c := range m.Deck.cards {
			if c.ID() == dup.ID() {
				t.Errorf("mutation %d (%s): result deck still holds %s",
					i, m.Description, dup.DisplayName())
			}
		}
	}
}

// Tests that pair mutations still emit when one pair half is already present, as long as
// the add side picks a non-overlapping variant.
func TestPairSwapMutations_FiresWhenOneHalfAlreadyPresent(t *testing.T) {
	reg := pairFixture()
	a := makeFakeCard(fakeC1)
	sk1 := makeFakeCard(fakeSK1)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{a, a, a, sk1})

	muts := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
	if len(muts) == 0 {
		t.Fatal("expected pair mutations even with one half present")
	}
	sawDifferentSecondVariantAdd := false
	for _, m := range muts {
		if strings.Contains(m.Description, "+1 "+(makeFakeCard(fakeSK2)).DisplayName()) ||
			strings.Contains(m.Description, "+1 "+(makeFakeCard(fakeSK3)).DisplayName()) {
			sawDifferentSecondVariantAdd = true
			break
		}
	}
	if !sawDifferentSecondVariantAdd {
		t.Error("expected at least one mutation adding a non-sk1 variant when sk1 is already present")
	}
}

// Tests that PairSwapMutations is cap-blind — emits maxCopies-violating candidates that
// FilterMaxCopiesViolations strips downstream.
func TestPairSwapMutations_GeneratesCapViolatingCandidates(t *testing.T) {
	reg := pairFixture()
	sk1 := makeFakeCard(fakeSK1)
	a := makeFakeCard(fakeC1)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{sk1, sk1, a, a})

	// 3 unique removed-ID combos after dedupe: (sk1, sk1), (sk1, a), (a, a). Overlap
	// suppression skips a combo when one of its removed IDs equals one of the add IDs;
	// sk1 is an add ID for 3 of the 9 (mw, sk) cross-products. So:
	//   (sk1, sk1) and (sk1, a) each emit 9 - 3 = 6 surviving combos.
	//   (a, a) emits all 9.
	// Total = 6 + 6 + 9 = 21.
	muts := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
	const want = 21
	if len(muts) != want {
		t.Fatalf("got %d pair mutations, want %d (cap-blind enumeration)", len(muts), want)
	}

	// At least one of those mutations must add sk1 again — pushing the count to 3 — which
	// would violate maxCopies=2. PairSwapMutations does NOT enforce that; the post-filter
	// in AllMutations does.
	sawCapViolator := false
	for _, m := range muts {
		if strings.Contains(m.Description, "+1 "+sk1.DisplayName()) {
			counts := map[ids.CardID]int{}
			for _, c := range m.Deck.cards {
				counts[c.ID()]++
			}
			if counts[sk1.ID()] > 2 {
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
func TestPairSwapMutations_HandlesUnbalancedHalfCounts(t *testing.T) {
	reg := pairFixture()
	mw1 := makeFakeCard(fakeMW1)
	mw2 := makeFakeCard(fakeMW2)
	mw3 := makeFakeCard(fakeMW3)
	sk1 := makeFakeCard(fakeSK1)
	sk2 := makeFakeCard(fakeSK2)
	sk3 := makeFakeCard(fakeSK3)
	a, b := makeFakeCard(fakeC1), makeFakeCard(fakeC2)
	cardsList := []Card{
		mw1, mw1, mw2, mw2, mw3,
		sk1, sk2, sk3,
		a, a, a, b, b, b,
	}
	d := New(fakeHero{}, []Weapon{reg.weapons[0]}, cardsList)

	muts := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
	if len(muts) == 0 {
		t.Fatal("expected pair mutations on unbalanced deck")
	}
	for i, m := range muts {
		if len(m.Deck.cards) != len(cardsList) {
			t.Errorf("mutation %d (%s): card count %d, want %d (size must stay stable)",
				i, m.Description, len(m.Deck.cards), len(cardsList))
		}
	}
}

// Every emitted pair mutation produces a deck with a different card multiset than the
// source. Defensive against a future bug where the overlap-suppression check misses a
// path that ends up at the source composition.
func TestPairSwapMutations_ResultDifferentFromSource(t *testing.T) {
	reg := pairFixture()
	a, b := makeFakeCard(fakeC1), makeFakeCard(fakeC2)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{a, a, b, b})
	srcKey := cardMultisetKey(d.cards)
	for i, m := range pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg)) {
		if cardMultisetKey(m.Deck.cards) == srcKey {
			t.Errorf("mutation %d (%s) produced a no-op (same multiset as source)", i, m.Description)
		}
	}
}

// Tests that PairSwapMutations suppresses combos that reduce to a single-slot swap (a
// -1/+1 of the same card), which the single-slot generator already covers.
func TestPairSwapMutations_OverlapSuppressionSkipsRedundantSwaps(t *testing.T) {
	reg := pairFixture()
	sk1 := makeFakeCard(fakeSK1)
	a := makeFakeCard(fakeC1)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{sk1, a, a, a})
	addStr := "+1 " + sk1.DisplayName()
	rmStr := "-1 " + sk1.DisplayName()
	for i, m := range pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg)) {
		if strings.Contains(m.Description, rmStr) && strings.Contains(m.Description, addStr) {
			t.Errorf("mutation %d (%s): redundant -1/+1 of %s — overlap suppression failed",
				i, m.Description, sk1.DisplayName())
		}
	}
}

// Tests that PairSwapMutations skips pair halves that aren't in the registry's legal
// pool. Builds a fixture where one half (sk2) is omitted from the registry and asserts
// no emitted mutation includes that ID on the add side.
func TestPairSwapMutations_SkipsPoolMissingHalves(t *testing.T) {
	missing := fakeSK2
	reg := &fakeRegistry{
		cards: []Card{
			makeFakeCard(fakeC1),
			makeFakeCard(fakeC2),
			makeFakeCard(fakeMW1),
			makeFakeCard(fakeMW2),
			makeFakeCard(fakeMW3),
			makeFakeCard(fakeSK1),
			makeFakeCard(fakeSK3),
			// fakeSK2 deliberately omitted
		},
		weapons: []Weapon{fakeWeapon{name: "weapon-A", hands: 1}},
	}
	a, b := makeFakeCard(fakeC1), makeFakeCard(fakeC2)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{a, a, b, b})
	missingDisplay := makeFakeCard(missing).DisplayName()
	for i, m := range pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg)) {
		if strings.Contains(m.Description, missingDisplay) {
			t.Errorf("mutation %d (%s): added pool-missing half %s",
				i, m.Description, missingDisplay)
		}
	}
}

// Tests that PairSwapMutations is deterministic: two back-to-back calls produce the same
// mutation sequence.
func TestPairSwapMutations_DeterministicOrdering(t *testing.T) {
	reg := pairFixture()
	a, b := makeFakeCard(fakeC1), makeFakeCard(fakeC2)
	d := New(fakeHero{}, []Weapon{reg.weapons[0]},
		[]Card{a, a, b, b})

	first := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
	second := pairSwapMutations(d, fakeMoonSunPair, buildLegalByID(reg))
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

// Tests that FilterMaxCopiesViolations drops any mutation whose result deck exceeds
// maxCopies of any card.
func TestFilterMaxCopiesViolations_StripsCapViolators(t *testing.T) {
	a := makeFakeCard(fakeC1)
	b := makeFakeCard(fakeC2)
	mw := makeFakeCard(fakeMW1)
	clean := New(fakeHero{}, []Weapon{fakeWeapon{name: "w", hands: 1}},
		[]Card{a, b, mw, mw})
	violator := New(fakeHero{}, []Weapon{fakeWeapon{name: "w", hands: 1}},
		[]Card{mw, mw, mw, mw, mw})

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

// RespectsMaxCopies returns false immediately when a count exceeds the cap, without
// scanning the full slice. Sentinel for the inner-loop fast path in
// FilterMaxCopiesViolations.
func TestRespectsMaxCopies_ShortCircuits(t *testing.T) {
	a := makeFakeCard(fakeC1)
	cs := []Card{a, a, a}
	if respectsMaxCopies(cs, 2) {
		t.Error("3 copies at maxCopies=2 should fail RespectsMaxCopies")
	}
	if !respectsMaxCopies(cs, 3) {
		t.Error("3 copies at maxCopies=3 should pass RespectsMaxCopies")
	}
}
