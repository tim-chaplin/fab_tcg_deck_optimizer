package deck_test

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// TestRandom_FilterExcludesRejected confirms the legal predicate is actually applied to the
// candidate pool: a filter that blocks Plunder Run (all variants) should produce decks that
// never contain any Plunder Run printing, even across many samples.
func TestRandom_FilterExcludesRejected(t *testing.T) {
	bannedIDs := map[ids.CardID]bool{
		ids.PlunderRunRed:    true,
		ids.PlunderRunYellow: true,
		ids.PlunderRunBlue:   true,
	}
	legal := func(c deck.Card) bool { return !bannedIDs[c.ID()] }
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := deck.Random(heroes.Viserai{}, 40, 2, rng, legal, registry.Registry{})
		uniq, _ := d.UniqueIDs()
		for _, id := range uniq {
			if bannedIDs[id] {
				t.Errorf("sample %d: %s was in the banlist", i, registry.GetCard(id).DisplayName())
			}
		}
	}
}

// TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal exercises the SanitizeNotImplemented
// happy path: tagged slots get replaced, size is preserved, the per-printing copy budget is
// respected, and one swap entry is reported per replaced slot.
func TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal(t *testing.T) {
	if _, ok := registry.GetCard(ids.StrikeGoldRed).(registry.NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	tagged := registry.GetCard(ids.StrikeGoldRed)
	safe := registry.GetCard(ids.ArcanicCrackleRed)
	if _, t2 := safe.(registry.NotImplemented); t2 {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented keeper for this test")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}},
		[]deck.Card{safe, safe, tagged, tagged})

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil, registry.Registry{})

	if len(replaced) != 2 {
		t.Errorf("replaced %d slots, want 2", len(replaced))
	}
	if d.Size() != 4 {
		t.Errorf("card count after sanitize = %d, want 4", d.Size())
	}
	uniq, _ := d.UniqueIDs()
	for _, id := range uniq {
		if _, ok := registry.GetCard(id).(registry.NotImplemented); ok {
			t.Errorf("%s still implements NotImplemented", registry.GetCard(id).DisplayName())
		}
	}
	for name, count := range d.NameCounts() {
		if count > 2 {
			t.Errorf("%s appears %d times, exceeds maxCopies=2", name, count)
		}
	}
	for _, r := range replaced {
		if r.From.ID() != ids.StrikeGoldRed {
			t.Errorf("replacement From = %s, want Strike Gold [R]", r.From.DisplayName())
		}
		if _, ok := r.To.(registry.NotImplemented); ok {
			t.Errorf("replacement To = %s implements NotImplemented", r.To.DisplayName())
		}
	}
}

// TestSanitizeNotImplemented_NoOpOnCleanDeck confirms the sanitizer is an identity operation
// when the deck already has no NotImplemented cards: no replacements, no composition change.
func TestSanitizeNotImplemented_NoOpOnCleanDeck(t *testing.T) {
	a := registry.GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(registry.NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	beforeFingerprint := d.Fingerprint()

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil, registry.Registry{})

	if len(replaced) != 0 {
		t.Errorf("replacements on clean deck = %d, want 0", len(replaced))
	}
	if got := d.Fingerprint(); got != beforeFingerprint {
		t.Errorf("composition mutated: fingerprint %q → %q", beforeFingerprint, got)
	}
}

// TestAllMutations_FilterExcludesRejectedAdditions confirms banned cards never appear as
// swap-in candidates while remaining valid removal targets when present in the starting deck.
func TestAllMutations_FilterExcludesRejectedAdditions(t *testing.T) {
	bannedIDs := map[ids.CardID]bool{
		ids.CriticalStrikeRed: true,
	}
	legal := func(c deck.Card) bool { return !bannedIDs[c.ID()] }

	cs := registry.GetCard(ids.CriticalStrikeRed)
	other := registry.GetCard(ids.AetherSlashRed)
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{cs, cs, other, other})

	for i, m := range deck.AllMutations(d, 2, registry.Registry{}, legal) {
		bannedIn := 0
		counts := m.Deck.NameCounts()
		for id := range bannedIDs {
			bannedIn += counts[registry.GetCard(id).DisplayName()]
		}
		// The starting deck has 2 copies of Critical Strike [R]. A mutation that removes one
		// leaves 1; a mutation that removes the other leaves 1; a weapon-only mutation leaves
		// all 2. No mutation should ADD another copy.
		if bannedIn > 2 {
			t.Errorf("mutation %d (%s): has %d banned copies, want <=2 (no additions allowed)",
				i, m.Description, bannedIn)
		}
	}
}
