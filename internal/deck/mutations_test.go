package deck

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

func TestAllMutations_CountsAndShape(t *testing.T) {
	// Build a tiny deck: 2 unique cards × 2 copies = 4 cards, plus one weapon.
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	muts := AllMutations(d, 2, nil)

	// Weapon mutations: every loadout except the current one. Card mutations at maxCopies=2:
	// for each of the 2 unique removals, every pool entry except self (no-op) and the other
	// in-deck card (already at cap) is a valid add — so 2 × (pool - 2). Use legalPool(nil)
	// instead of cards.Deckable() directly so the expected count tracks AllMutations's own
	// filtering (NotImplemented cards are skipped).
	loadouts := weaponLoadouts(weapon.All)
	pool := legalPool(nil)
	wantWeaponMuts := len(loadouts) - 1
	wantCardMuts := 2 * (len(pool) - 2)
	want := wantWeaponMuts + wantCardMuts

	if len(muts) != want {
		t.Fatalf("got %d mutations, want %d (%d weapon + %d card)",
			len(muts), want, wantWeaponMuts, wantCardMuts)
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
// maxCopies should open up adds to cards already in the deck that are below the cap.
func TestAllMutations_OddCountsAllowed(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

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
		counts := map[card.ID]int{}
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

// TestAllMutations_OrdersByAscendingAvg pins the iterate-friendly ordering: the removed card in
// the first card-mutation batch should be the one with the lowest per-card Avg in the current
// deck's stats. Run with both (lower-ID is low-avg) and (higher-ID is low-avg) so the test fails
// if the implementation accidentally sorts only by card.ID and gets a free pass from whichever
// direction happens to align.
func TestAllMutations_OrdersByAscendingAvg(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)  // lower card.ID
	b := cards.Get(card.ArcanicSpikeRed) // higher card.ID

	cases := []struct {
		name        string
		lowAvgCard  card.Card
		highAvgCard card.Card
	}{
		{"low-avg card has lower ID", a, b},
		{"low-avg card has higher ID", b, a},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
				[]card.Card{a, a, b, b})
			// Both cards see the same Plays so only Avg (= TotalContribution / Plays) drives the
			// ordering — no path for card.ID to sneak in via a sub-ordering rule.
			d.Stats.PerCard = map[card.ID]CardPlayStats{
				tc.lowAvgCard.ID():  {Plays: 10, TotalContribution: 10}, // Avg 1.0
				tc.highAvgCard.ID(): {Plays: 10, TotalContribution: 80}, // Avg 8.0
			}

			muts := AllMutations(d, 2, nil)
			// Skip the weapon-mutation block (len(loadouts)-1 entries, one per alternative loadout).
			firstCardMut := muts[len(weaponLoadouts(weapon.All))-1]
			wantPrefix := "-1 " + tc.lowAvgCard.Name() + ","
			if !strings.HasPrefix(firstCardMut.Description, wantPrefix) {
				t.Errorf("first card mutation removed wrong card\n  got:  %q\n  want prefix: %q",
					firstCardMut.Description, wantPrefix)
			}
		})
	}
}

// TestAllMutations_PreservesSideboard pins that every derived Mutation inherits the source
// deck's Sideboard verbatim. Without this guarantee an anneal round would silently drop the
// user's hand-managed sideboard as soon as it accepted a mutation and wrote the deck back.
func TestAllMutations_PreservesSideboard(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})
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
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := AllMutations(d, 2, nil)
	second := AllMutations(d, 2, nil)

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
		for j, c := range first[i].Deck.Cards {
			if c.ID() != second[i].Deck.Cards[j].ID() {
				t.Errorf("mutation %d card[%d] differs between calls: %v vs %v",
					i, j, c.ID(), second[i].Deck.Cards[j].ID())
			}
		}
	}
}

func TestAllMutations_NoDuplicateOfSource(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	srcKey := deckFingerprint(d)
	for i, m := range AllMutations(d, 2, nil) {
		if deckFingerprint(m.Deck) == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
		}
	}
}
