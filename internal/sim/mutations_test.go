package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

func TestAllMutations_CountsAndShape(t *testing.T) {
	// Build a tiny deck: 2 unique cards × 2 copies = 4 cards, plus one weapon. Both starter
	// cards must be implemented (NOT carrying NotImplemented) so the LegalPool / removal-
	// counting math below holds. ArcanicCrackleRed and ArcanicSpikeRed are stable picks.
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	muts := AllMutations(d, 2, nil)

	// Weapon mutations: every loadout except the current one. Card mutations at maxCopies=2:
	// for each of the 2 unique removals, every pool entry except self (no-op) and the other
	// in-deck card (already at cap) is a valid add — so 2 × (pool - 2). Pair mutations: for
	// each registered cardPair whose halves are both absent, C(min(uniques, K), 2) candidates
	// — with 2 unique deck IDs that's 1 per absent pair. Both halves absent ⇒ all pairs
	// contribute. Use LegalPool(nil) so the count tracks AllMutations's own filtering
	// (NotImplemented cards are skipped).
	loadouts := WeaponLoadouts(AllWeapons)
	pool := LegalPool(nil)
	wantWeaponMuts := len(loadouts) - 1
	wantCardMuts := 2 * (len(pool) - 2)
	wantPairMuts := expectedPairMutCount(d, 2)
	want := wantWeaponMuts + wantCardMuts + wantPairMuts

	if len(muts) != want {
		t.Fatalf("got %d mutations, want %d (%d weapon + %d card + %d pair)",
			len(muts), want, wantWeaponMuts, wantCardMuts, wantPairMuts)
	}
	for i, m := range muts {
		if len(m.Deck.Cards) != 4 {
			t.Errorf("mutation %d: card count %d, want 4", i, len(m.Deck.Cards))
		}
		if m.Deck.Hero.Name() != d.Hero.Name() {
			t.Errorf("mutation %d: hero changed", i)
		}
		if m.Description == "" {
			t.Errorf("mutation %d: empty description", i)
		}
	}
}

// TestAllMutations_OddCountsAllowed exercises the single-card-swap semantics: a mutation may leave
// the deck with an odd number of any given printing (e.g. 1×A + 3×B at maxCopies=3), and raising
// maxCopies should open up adds to cards already in the deck that are below the cap. Both
// starter cards must be implemented so the diff math holds (NotImplemented removals are absent
// from the addID pool, which suppresses the "swap to in-deck other" mutation pair the diff
// expects).
func TestAllMutations_OddCountsAllowed(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	// At maxCopies=3, each of the 2 in-deck cards (a, b) is below the cap, so "remove a, add b"
	// (and the mirror) become legal. That's 2 more card mutations than the maxCopies=2 case.
	mutsLow := AllMutations(d, 2, nil)
	mutsHigh := AllMutations(d, 3, nil)
	if len(mutsHigh)-len(mutsLow) != 2 {
		t.Errorf("maxCopies=3 should produce exactly 2 more mutations than maxCopies=2; got diff=%d",
			len(mutsHigh)-len(mutsLow))
	}

	// Every mutation at maxCopies=3 still has exactly 4 cards; some should have an odd count of
	// one card — single-card swaps can leave odd-count slots when maxCopies allows it.
	sawOdd := false
	for _, m := range mutsHigh {
		if len(m.Deck.Cards) != 4 {
			t.Errorf("card count %d, want 4", len(m.Deck.Cards))
		}
		counts := map[ids.CardID]int{}
		for _, c := range m.Deck.Cards {
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

// TestAllMutations_PreservesSideboard pins that every derived Mutation inherits the source
// deck's Sideboard verbatim. Without this guarantee an anneal round would silently drop the
// user's hand-managed sideboard as soon as it accepted a mutation and wrote the deck back.
func TestAllMutations_PreservesSideboard(t *testing.T) {
	a := GetCard(ids.AetherSlashRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})
	d.Sideboard = []string{a.Name(), b.Name(), b.Name()}

	muts := AllMutations(d, 2, nil)
	if len(muts) == 0 {
		t.Fatal("expected at least one mutation")
	}

	wantCounts := map[string]int{a.Name(): 1, b.Name(): 2}
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

func TestAllMutations_Deterministic(t *testing.T) {
	a := GetCard(ids.AetherSlashRed)
	b := GetCard(ids.ArcanicSpikeRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, b, b})

	first := AllMutations(d, 2, nil)
	second := AllMutations(d, 2, nil)

	if len(first) != len(second) {
		t.Fatalf("mutation counts differ between calls: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if WeaponKey(first[i].Deck.Weapons) != WeaponKey(second[i].Deck.Weapons) {
			t.Errorf("mutation %d weapons differ between calls", i)
		}
		if first[i].Description != second[i].Description {
			t.Errorf("mutation %d descriptions differ: %q vs %q",
				i, first[i].Description, second[i].Description)
		}
		for j, c := range first[i].Deck.Cards {
			if c.ID() != second[i].Deck.Cards[j].ID() {
				t.Errorf("mutation %d card[%d] differs between calls: %v vs %v",
					i, j, c.ID(), second[i].Deck.Cards[j].ID())
			}
		}
	}
}

func TestAllMutations_NoDuplicateOfSource(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	d := New(heroes.Viserai{}, []Weapon{weapons.NebulaBlade{}}, []Card{a, a, a, a})
	srcKey := DeckFingerprint(d)
	for i, m := range AllMutations(d, 2, nil) {
		if DeckFingerprint(m.Deck) == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
		}
	}
}

// expectedPairMutCount mirrors pairSwapMutations's emission rule for a given deck so the
// CountsAndShape test can predict the pair-mutation contribution without re-implementing the
// generator. For each registered pair × variant cross-product, counts the distinct
// (sorted-removed-IDs) combos that survive overlap suppression. Combos rejected by the
// shared maxCopies post-filter are dropped at the end.
func expectedPairMutCount(d *Deck, maxCopies int) int {
	// Build the deduped removed-ID combo set the index-based generator collapses into.
	type idPair struct{ a, b ids.CardID }
	combos := map[idPair]bool{}
	for i := 0; i < len(d.Cards); i++ {
		for j := i + 1; j < len(d.Cards); j++ {
			a, b := SortedIDPair(d.Cards[i].ID(), d.Cards[j].ID())
			combos[idPair{a, b}] = true
		}
	}

	srcCounts := map[ids.CardID]int{}
	for _, c := range d.Cards {
		srcCounts[c.ID()]++
	}

	total := 0
	for _, p := range CardPairs {
		for _, fID := range p.First {
			if !PairAddAllowed(GetCard(fID), nil) {
				continue
			}
			for _, sID := range p.Second {
				if !PairAddAllowed(GetCard(sID), nil) {
					continue
				}
				for combo := range combos {
					if combo.a == fID || combo.a == sID ||
						combo.b == fID || combo.b == sID {
						continue
					}
					// Apply the shared maxCopies post-filter: simulate the swap and reject
					// when any ID exceeds the cap in the result deck.
					if !swapRespectsMaxCopies(srcCounts, combo.a, combo.b, fID, sID, maxCopies) {
						continue
					}
					total++
				}
			}
		}
	}
	return total
}

// swapRespectsMaxCopies checks whether removing one of each removeA / removeB and adding
// one of each addA / addB leaves all card counts at or below maxCopies. Mirrors
// RespectsMaxCopies's contract on the pre-mutation count map without rebuilding the deck.
func swapRespectsMaxCopies(srcCounts map[ids.CardID]int, removeA, removeB, addA, addB ids.CardID, maxCopies int) bool {
	delta := map[ids.CardID]int{}
	delta[removeA]--
	delta[removeB]--
	delta[addA]++
	delta[addB]++
	for id, change := range delta {
		if srcCounts[id]+change > maxCopies {
			return false
		}
	}
	return true
}
