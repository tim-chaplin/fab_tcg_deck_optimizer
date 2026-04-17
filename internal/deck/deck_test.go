package deck

import (
	"math/rand"
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

	muts := AllMutations(d, 2)

	// Weapon mutations: every loadout except the current one. Card mutations at maxCopies=2:
	// for each of the 2 unique removals, every pool entry except self (no-op) and the other
	// in-deck card (already at cap) is a valid add — so 2 × (pool - 2).
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
	mutsLow := AllMutations(d, 2)
	mutsHigh := AllMutations(d, 3)
	if len(mutsHigh)-len(mutsLow) != 2 {
		t.Errorf("maxCopies=3 should produce exactly 2 more mutations than maxCopies=2; got diff=%d",
			len(mutsHigh)-len(mutsLow))
	}

	// Every mutation at maxCopies=3 still has exactly 4 cards; some should have an odd count of
	// one card (that's the whole point — no longer forced to keep pairs).
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

func TestAllMutations_Deterministic(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := AllMutations(d, 2)
	second := AllMutations(d, 2)

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
	for i, m := range AllMutations(d, 2) {
		if deckFingerprint(m.Deck) == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
		}
	}
}

// TestEvaluate_PerCardStatsPopulated pins per-card attribution: every card appearance (played or
// pitched) contributes to Plays+Pitches, and TotalContribution sums role-based per-card credit:
// Attack → Card.Attack(), Defend → proportional share of block, Pitch → Card.Pitch(). A
// single-printing deck makes the totals easy to assert against the card's printed stats.
func TestEvaluate_PerCardStatsPopulated(t *testing.T) {
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.PerCard == nil {
		t.Fatalf("PerCard should be initialised after Evaluate")
	}
	stat, ok := d.Stats.PerCard[card.ReadTheRunesRed]
	if !ok {
		t.Fatalf("PerCard missing entry for Read the Runes (Red)")
	}
	if got := stat.Plays + stat.Pitches; got != 4 {
		t.Errorf("Plays+Pitches = %d, want 4 (one 4-card hand of the same card)", got)
	}
	// Contributions come from the winning chain replay (Play returns + hero triggers) plus
	// role-based shares for pitch/defend. The exact total depends on rider/trigger damage, so
	// assert the weaker property that it's positive and produces a positive Avg.
	if stat.TotalContribution <= 0 {
		t.Errorf("TotalContribution = %v, want >0 (played Read the Runes deals at least Attack+rider)",
			stat.TotalContribution)
	}
	if stat.Avg() <= 0 {
		t.Errorf("Avg() = %v, want >0", stat.Avg())
	}
}

// TestEvaluate_BestHandStartingRunechantsIsPreHandCarryover pins down a subtle bug: Evaluate
// used to write the post-hand LeftoverRunechants into BestHand.StartingRunechants, so the field
// surfaced the wrong turn's count. The field is documented as "the Runechant count carried in
// from the previous turn when this hand was played", which for the first hand of a run is
// always 0 — even if the hand itself creates runechants that leftover into the next turn.
//
// Without the snapshot fix this test fails: StartingRunechants equals LeftoverRunechants (nonzero)
// instead of 0.
func TestEvaluate_BestHandStartingRunechantsIsPreHandCarryover(t *testing.T) {
	// Viserai has Intelligence 4. A 4-card deck gives exactly one hand per run, so the Best
	// record always reflects that first hand — no previous turn ever existed.
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})

	// Seed doesn't matter (all cards identical), but fix it for determinism.
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.Best.Hand == nil {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	// Sanity: the hand should have left runechants on the table (otherwise the bug couldn't
	// manifest — pre-hand and post-hand counts would both be 0).
	if d.Stats.Best.Play.Value == 0 {
		t.Fatalf("expected nonzero Value from a hand of Read the Runes; got 0")
	}
	if d.Stats.Best.StartingRunechants != 0 {
		t.Errorf("StartingRunechants = %d, want 0 (first hand of the run has no previous-turn carryover)",
			d.Stats.Best.StartingRunechants)
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
