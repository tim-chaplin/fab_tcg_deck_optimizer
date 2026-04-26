package deck

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestRandom_FilterExcludesRejected confirms the legal predicate is actually applied to the
// candidate pool: a filter that blocks Plunder Run (all variants) should produce decks that
// never contain any Plunder Run printing, even across many samples.
func TestRandom_FilterExcludesRejected(t *testing.T) {
	bannedIDs := map[card.ID]bool{
		card.PlunderRunRed:    true,
		card.PlunderRunYellow: true,
		card.PlunderRunBlue:   true,
	}
	legal := func(c card.Card) bool { return !bannedIDs[c.ID()] }
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := Random(hero.Viserai{}, 40, 2, rng, legal)
		for j, c := range d.Cards {
			if bannedIDs[c.ID()] {
				t.Errorf("sample %d: card[%d] = %s was in the banlist", i, j, c.Name())
			}
		}
	}
}

// TestLegalPool_SkipsNotImplemented pins the contract that cards tagged with
// card.NotImplemented never land in the search pool, with or without a legal predicate. The
// property holds regardless of which registered cards carry the tag today, so it doubles as
// a regression guard against accidental weakening of the filter.
func TestLegalPool_SkipsNotImplemented(t *testing.T) {
	for _, pred := range []func(card.Card) bool{nil, func(card.Card) bool { return true }} {
		for _, id := range legalPool(pred) {
			c := cards.Get(id)
			if _, ok := c.(card.NotImplemented); ok {
				t.Errorf("legalPool included NotImplemented card %s", c.Name())
			}
		}
	}
}

// TestRandom_ExcludesNotImplemented confirms no sampled random deck contains a card tagged
// with card.NotImplemented. Drives Random over many seeds so a leak into the pool would show
// up as a failure even when the tagged set is small relative to the full pool.
func TestRandom_ExcludesNotImplemented(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := Random(hero.Viserai{}, 40, 2, rng, nil)
		for j, c := range d.Cards {
			if _, ok := c.(card.NotImplemented); ok {
				t.Errorf("sample %d card[%d] = %s implements NotImplemented", i, j, c.Name())
			}
		}
	}
}

// TestLegalPool_ExcludesTaggedCardsByID gives TestLegalPool_SkipsNotImplemented teeth: it
// picks a concrete registered card we know currently carries the NotImplemented marker
// (Strike Gold [R], gold-token rider) and asserts it's absent from legalPool's output.
// Without at least one real tagged card the property test is vacuous, so this guards against
// a regression where the marker interface itself silently breaks. Self-retires if Strike Gold
// ever loses the tag (gold-token economy gets modelled) so maintenance is only a delete.
func TestLegalPool_ExcludesTaggedCardsByID(t *testing.T) {
	if _, ok := cards.Get(card.StrikeGoldRed).(card.NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	for _, id := range legalPool(nil) {
		if id == card.StrikeGoldRed {
			t.Fatalf("legalPool included Strike Gold [R] despite its NotImplemented tag")
		}
	}
}

// TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal drives the sanitizer
// against a deck that starts with two NotImplemented copies in it (Strike Gold Red is a real
// tagged card). After sanitization the deck must: (a) have zero NotImplemented cards, (b)
// be the same size, (c) respect maxCopies across the post-sanitize distribution, (d) report
// exactly two swaps, each naming the original tagged card.
func TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal(t *testing.T) {
	if _, ok := cards.Get(card.StrikeGoldRed).(card.NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	tagged := cards.Get(card.StrikeGoldRed)
	safe := cards.Get(card.ArcanicCrackleRed)
	if _, t2 := safe.(card.NotImplemented); t2 {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented keeper for this test")
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
		[]card.Card{safe, safe, tagged, tagged})

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil)

	if len(replaced) != 2 {
		t.Errorf("replaced %d slots, want 2", len(replaced))
	}
	if len(d.Cards) != 4 {
		t.Errorf("card count after sanitize = %d, want 4", len(d.Cards))
	}
	for i, c := range d.Cards {
		if _, ok := c.(card.NotImplemented); ok {
			t.Errorf("card[%d] = %s still implements NotImplemented", i, c.Name())
		}
	}
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
		if counts[c.ID()] > 2 {
			t.Errorf("%s appears %d times, exceeds maxCopies=2", c.Name(), counts[c.ID()])
		}
	}
	for _, r := range replaced {
		if r.From.ID() != card.StrikeGoldRed {
			t.Errorf("replacement From = %s, want Strike Gold [R]", r.From.Name())
		}
		if _, ok := r.To.(card.NotImplemented); ok {
			t.Errorf("replacement To = %s implements NotImplemented", r.To.Name())
		}
	}
}

// TestSanitizeNotImplemented_NoOpOnCleanDeck confirms the sanitizer is an identity operation
// when the deck already has no NotImplemented cards: no replacements, no mutations to
// Cards.
func TestSanitizeNotImplemented_NoOpOnCleanDeck(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	if _, tagged := a.(card.NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel")
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	before := append([]card.Card(nil), d.Cards...)

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil)

	if len(replaced) != 0 {
		t.Errorf("replacements on clean deck = %d, want 0", len(replaced))
	}
	for i, c := range d.Cards {
		if c.ID() != before[i].ID() {
			t.Errorf("card[%d] mutated: %s → %s", i, before[i].Name(), c.Name())
		}
	}
}

// TestAllMutations_ExcludesNotImplementedAdditions confirms no single-slot mutation can
// introduce a NotImplemented card. Starting deck must contain only implemented cards so any
// NotImplemented copy in a mutation output must have come from the add pool;
// ArcanicCrackleRed is the chosen sentinel (no NotImplemented marker).
func TestAllMutations_ExcludesNotImplementedAdditions(t *testing.T) {
	a := cards.Get(card.ArcanicCrackleRed)
	if _, tagged := a.(card.NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel for this test")
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	for _, m := range AllMutations(d, 2, nil) {
		for _, c := range m.Deck.Cards {
			if _, ok := c.(card.NotImplemented); ok {
				t.Errorf("%s introduced NotImplemented card %s", m.Description, c.Name())
			}
		}
	}
}

// TestAllMutations_FilterExcludesRejectedAdditions confirms banned cards never appear as
// swap-in candidates. A banned card already in the deck IS still a valid removal target — the
// hill climb must be able to swap it out — so we assert that the starting deck's banned card is
// never in the post-mutation card list either (which would require it to have been added back).
func TestAllMutations_FilterExcludesRejectedAdditions(t *testing.T) {
	bannedIDs := map[card.ID]bool{
		card.PlunderRunRed: true,
	}
	legal := func(c card.Card) bool { return !bannedIDs[c.ID()] }

	pr := cards.Get(card.PlunderRunRed)
	other := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{pr, pr, other, other})

	for i, m := range AllMutations(d, 2, legal) {
		bannedIn := 0
		for _, c := range m.Deck.Cards {
			if bannedIDs[c.ID()] {
				bannedIn++
			}
		}
		// The starting deck has 2 copies of Plunder Run [R]. A mutation that removes one leaves
		// 1; a mutation that removes the other leaves 1; a weapon-only mutation leaves all 2. No
		// mutation should ADD another copy.
		if bannedIn > 2 {
			t.Errorf("mutation %d (%s): has %d banned copies, want <=2 (no additions allowed)",
				i, m.Description, bannedIn)
		}
	}
}
