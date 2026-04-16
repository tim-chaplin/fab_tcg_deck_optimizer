package deck

import (
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

	muts := AllMutations(d)

	// Weapon mutations: every loadout except NebulaBlade. Card mutations: 2 unique in deck ×
	// (pool size − 2 in-deck) replacements.
	loadouts := weaponLoadouts(cards.AllWeapons)
	pool := cards.Deckable()
	wantWeaponMuts := len(loadouts) - 1
	wantCardMuts := 2 * (len(pool) - 2)
	want := wantWeaponMuts + wantCardMuts

	if len(muts) != want {
		t.Fatalf("got %d mutations, want %d (%d weapon + %d card)",
			len(muts), want, wantWeaponMuts, wantCardMuts)
	}
	for i, m := range muts {
		if len(m.Cards) != 4 {
			t.Errorf("mutation %d: card count %d, want 4", i, len(m.Cards))
		}
		if m.Hero.Name() != d.Hero.Name() {
			t.Errorf("mutation %d: hero changed", i)
		}
	}
}

func TestAllMutations_Deterministic(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := AllMutations(d)
	second := AllMutations(d)

	if len(first) != len(second) {
		t.Fatalf("mutation counts differ between calls: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if weaponKey(first[i].Weapons) != weaponKey(second[i].Weapons) {
			t.Errorf("mutation %d weapons differ between calls", i)
		}
		for j, c := range first[i].Cards {
			if c.ID() != second[i].Cards[j].ID() {
				t.Errorf("mutation %d card[%d] differs between calls: %v vs %v",
					i, j, c.ID(), second[i].Cards[j].ID())
			}
		}
	}
}

func TestAllMutations_NoDuplicateOfSource(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	srcKey := deckFingerprint(d)
	for i, m := range AllMutations(d) {
		if deckFingerprint(m) == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
		}
	}
}

// deckFingerprint builds a comparable summary of a deck for equality checks in tests.
func deckFingerprint(d *Deck) string {
	s := weaponKey(d.Weapons) + "|"
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	// Stable ordering — iterate over all possible IDs in byID order isn't exposed, so use a
	// sorted slice of (id, count).
	type pair struct {
		id card.ID
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
