package deck

import (
	"math/rand"
	"reflect"
	"slices"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// fakeWeapon is a tiny Weapon implementation for the package's tests. Hand count drives
// loadout legality; name disambiguates dual-wielded slots.
type fakeWeapon struct {
	name  string
	hands int
}

func (w fakeWeapon) Name() string { return w.name }
func (w fakeWeapon) Hands() int   { return w.hands }

// fakeCard is a tiny Card implementation. ID is the per-printing key the copy-budget
// arithmetic operates on; DisplayName is the canonical sideboard-merge key. The other
// methods are zero-value stubs the deck never touches — deck only reads ID, Name, and
// DisplayName, but now that Card is aliased to the rich card.Card the full method set has
// to satisfy the interface.
type fakeCard struct {
	id      ids.CardID
	display string
}

func (c fakeCard) ID() ids.CardID                                     { return c.id }
func (c fakeCard) Name() string                                       { return c.display }
func (c fakeCard) DisplayName() string                                { return c.display }
func (c fakeCard) Cost() int                                          { return 0 }
func (c fakeCard) Pitch() int                                         { return 0 }
func (c fakeCard) Attack() int                                        { return 0 }
func (c fakeCard) Defense() int                                       { return 0 }
func (c fakeCard) Types(card.GameEngine) card.TypeSet                 { return 0 }
func (c fakeCard) GoAgain(card.GameEngine) bool                       { return false }
func (c fakeCard) Play(card.GameEngine, card.Logger, *card.CardState) {}

// fakeRegistry is the in-memory registry the tests build against.
type fakeRegistry struct {
	cards   []Card
	weapons []Weapon
}

func (r *fakeRegistry) LegalCards() []Card     { return slices.Clone(r.cards) }
func (r *fakeRegistry) LegalWeapons() []Weapon { return slices.Clone(r.weapons) }

func newFakeRegistry() *fakeRegistry {
	return &fakeRegistry{
		cards: []Card{
			fakeCard{id: 1, display: "Aether Slash [R]"},
			fakeCard{id: 2, display: "Aether Slash [Y]"},
			fakeCard{id: 3, display: "Read the Runes [R]"},
			fakeCard{id: 4, display: "Hocus Pocus [B]"},
		},
		weapons: []Weapon{fakeWeapon{name: "Reaping Blade", hands: 2}, fakeWeapon{name: "Scepter of Pain", hands: 1}},
	}
}

// Tests that New panics on a 3-weapon loadout (max 2) and on a 2H + 1H loadout (both must
// be 1H when 2 are equipped).
func TestNew_RejectsIllegalWeaponLoadouts(t *testing.T) {
	cases := []struct {
		name    string
		weapons []Weapon
	}{
		{"three weapons", []Weapon{
			fakeWeapon{name: "a", hands: 1},
			fakeWeapon{name: "b", hands: 1},
			fakeWeapon{name: "c", hands: 1},
		}},
		{"2H plus 1H", []Weapon{
			fakeWeapon{name: "two", hands: 2},
			fakeWeapon{name: "one", hands: 1},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("New(%v) didn't panic", tc.weapons)
				}
			}()
			New(nil, tc.weapons, nil)
		})
	}
}

// testDefaults is the loadout the ApplyDefaults tests merge against. The names are made
// up — generic deck logic doesn't care about specific card identity.
var testDefaults = Defaults{
	Equipment: []string{"Helm A", "Boots A", "Robes A"},
	Sideboard: []SideboardDefault{
		{Name: "Side Card", Count: 1},
		{Name: "Read the Runes [R]", Count: 2},
	},
}

// Tests that ApplyDefaults adds each Equipment entry once and is idempotent on the second
// pass.
func TestApplyDefaults_EquipmentIdempotent(t *testing.T) {
	d := New(nil, nil, nil)
	d.ApplyDefaults(testDefaults)
	want := []string{"Helm A", "Boots A", "Robes A"}
	if !reflect.DeepEqual(d.Equipment, want) {
		t.Errorf("after first ApplyDefaults: Equipment = %v, want %v", d.Equipment, want)
	}
	d.ApplyDefaults(testDefaults)
	if !reflect.DeepEqual(d.Equipment, want) {
		t.Errorf("after second ApplyDefaults: Equipment = %v, want %v (should be a no-op)", d.Equipment, want)
	}
}

// Tests that ApplyDefaults adds each Sideboard entry up to its target count and is
// idempotent on the second pass.
func TestApplyDefaults_SideboardIdempotent(t *testing.T) {
	d := New(nil, nil, nil)
	d.ApplyDefaults(testDefaults)
	want := []string{"Side Card", "Read the Runes [R]", "Read the Runes [R]"}
	if !reflect.DeepEqual(d.Sideboard, want) {
		t.Errorf("after first ApplyDefaults: Sideboard = %v, want %v", d.Sideboard, want)
	}
	d.ApplyDefaults(testDefaults)
	if !reflect.DeepEqual(d.Sideboard, want) {
		t.Errorf("after second ApplyDefaults: Sideboard = %v, want %v (should be a no-op)", d.Sideboard, want)
	}
}

// Tests that the sideboard merge respects SideboardCopyCap against main-deck copies —
// when the main deck already holds the per-card cap, the default sideboard add is skipped.
func TestApplyDefaults_SideboardClampsAgainstMainDeck(t *testing.T) {
	red := fakeCard{id: 3, display: "Read the Runes [R]"}
	d := New(nil, nil, []Card{red, red})
	d.ApplyDefaults(testDefaults)

	for _, name := range d.Sideboard {
		if name == "Read the Runes [R]" {
			t.Errorf("Sideboard added Read the Runes [R] despite main deck already holding %d copies (cap is %d)",
				2, SideboardCopyCap)
		}
	}
}

// Tests that a card present in the main deck is trimmed out of an existing sideboard, so the two
// sections never hold a combined count over the cap (the both-sections duplication bug).
func TestApplyDefaults_SideboardTrimsCopiesNowInMain(t *testing.T) {
	red := fakeCard{id: 3, display: "Read the Runes [R]"}
	d := New(nil, nil, []Card{red, red}) // main holds the cap
	d.Sideboard = []string{"Read the Runes [R]", "Read the Runes [R]", "Side Card"}
	d.ApplyDefaults(testDefaults)

	got := map[string]int{}
	for _, name := range d.Sideboard {
		got[name]++
	}
	if got["Read the Runes [R]"] != 0 {
		t.Errorf("Read the Runes [R] left in sideboard %dx while the main holds the cap (%d); want 0",
			got["Read the Runes [R]"], SideboardCopyCap)
	}
	if got["Side Card"] != 1 {
		t.Errorf("Side Card = %dx, want 1 (a non-main sideboard card is untouched)", got["Side Card"])
	}
}

// Tests the partial trim: with one copy in the main deck, the sideboard keeps only (cap - 1).
func TestApplyDefaults_SideboardTrimsToRemainingRoom(t *testing.T) {
	red := fakeCard{id: 3, display: "Read the Runes [R]"}
	d := New(nil, nil, []Card{red}) // main holds one copy
	d.Sideboard = []string{"Read the Runes [R]", "Read the Runes [R]"}
	d.ApplyDefaults(testDefaults)

	n := 0
	for _, name := range d.Sideboard {
		if name == "Read the Runes [R]" {
			n++
		}
	}
	if want := SideboardCopyCap - 1; n != want {
		t.Errorf("Read the Runes [R] in sideboard %dx, want %d (cap %d - 1 in main)",
			n, want, SideboardCopyCap)
	}
}

// Tests that Random builds a deck with the requested size, picks weapons from the
// legal pool, and never exceeds maxCopies for any single printing.
func TestRandom_BuildsLegalDeckWithinCopyBudget(t *testing.T) {
	reg := newFakeRegistry()
	rng := rand.New(rand.NewSource(42))

	d := Random(nil, 8, 2, rng, reg)
	if len(d.cards) != 8 {
		t.Errorf("deck size = %d, want 8", len(d.cards))
	}
	counts := map[ids.CardID]int{}
	for _, c := range d.cards {
		counts[c.ID()]++
	}
	for id, n := range counts {
		if n > 2 {
			t.Errorf("card id %d appears %d times, want ≤ 2 (maxCopies)", id, n)
		}
	}
	if len(d.Weapons) == 0 || len(d.Weapons) > 2 {
		t.Errorf("Weapons count = %d, want 1 or 2", len(d.Weapons))
	}
}

// Tests Shuffle's two-part contract: post-shuffle multiset equals the pre-shuffle multiset,
// and at least one position changes between two seeds (proving it isn't a no-op).
func TestShuffle_RandomisesCardsInPlace(t *testing.T) {
	master := New(nil, nil, []Card{
		fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}, fakeCard{id: 4},
		fakeCard{id: 5}, fakeCard{id: 6}, fakeCard{id: 7}, fakeCard{id: 8},
	})
	beforeCounts := map[ids.CardID]int{}
	for _, c := range master.cards {
		beforeCounts[c.ID()]++
	}
	d := master.Copy()
	d.Shuffle(rand.New(rand.NewSource(1)))
	afterCounts := map[ids.CardID]int{}
	for _, c := range d.cards {
		afterCounts[c.ID()]++
	}
	if !reflect.DeepEqual(beforeCounts, afterCounts) {
		t.Errorf("multiset changed across Shuffle: before=%v after=%v", beforeCounts, afterCounts)
	}

	d2 := master.Copy()
	d2.Shuffle(rand.New(rand.NewSource(2)))
	if reflect.DeepEqual(d.cards, d2.cards) {
		t.Errorf("two seeds produced identical orderings: %v", d.cards)
	}
}

// TestCopy_IsolatesShuffleFromMaster pins the master / copy split: after Copy + Shuffle on
// the copy, the master's Cards order is unchanged. Single-test cover for the contract the
// per-shuffle eval loop relies on (master shared across goroutines, each worker mutates
// its own Copy).
func TestCopy_IsolatesShuffleFromMaster(t *testing.T) {
	master := New(nil, nil, []Card{
		fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}, fakeCard{id: 4},
	})
	before := append([]Card(nil), master.cards...)
	worker := master.Copy()
	worker.Shuffle(rand.New(rand.NewSource(99)))
	if !reflect.DeepEqual(master.cards, before) {
		t.Errorf("master mutated by worker.Shuffle: got %v, want %v", master.cards, before)
	}
}

// TestDraw_ConsumesFromTop pins the contract: Draw(n) returns the top n cards in deck
// order and advances the deck so Size = original − n. Asserts against the first card via
// PeekTop so the test doesn't reach for the (test-only) Cards backing slice.
func TestDraw_ConsumesFromTop(t *testing.T) {
	d := New(nil, nil, []Card{
		fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}, fakeCard{id: 4},
	})
	top := d.PeekTop()
	drawn := d.Draw(2)
	if len(drawn) != 2 || drawn[0] != top {
		t.Errorf("Draw(2) = %v, want first card to match PeekTop %v", drawn, top)
	}
	if d.Size() != 2 {
		t.Errorf("Size after Draw(2) of a 4-card deck = %d, want 2", d.Size())
	}
}

// PeekTop returns nil when the deck is empty, and the top card otherwise. Drawing the
// whole deck and Peeking should give nil; PutBottom-ing one and re-Peeking gives that
// card back.
func TestPeekTop_ReturnsTopOrNil(t *testing.T) {
	d := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}})
	if got := d.PeekTop(); got == nil || got.ID() != 1 {
		t.Errorf("PeekTop on full deck = %v, want card id 1", got)
	}
	d.Draw(2)
	if got := d.PeekTop(); got != nil {
		t.Errorf("PeekTop on empty deck = %v, want nil", got)
	}
	d.PutBottom([]Card{fakeCard{id: 9}})
	if got := d.PeekTop(); got == nil || got.ID() != 9 {
		t.Errorf("PeekTop after PutBottom on empty deck = %v, want card id 9", got)
	}
}

// TestDraw_PanicsOnOverdraw confirms Draw(n) panics when n exceeds Size instead of
// returning a partial slice. Production callers check Size first.
func TestDraw_PanicsOnOverdraw(t *testing.T) {
	d := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}})
	defer func() {
		if recover() == nil {
			t.Errorf("Draw(3) on a 2-card deck didn't panic")
		}
	}()
	d.Draw(3)
}

// TestPutBottom_AppendsToBottom confirms PutBottom puts cards at the bottom of the deck
// in the order passed, preserving the existing top.
func TestPutBottom_AppendsToBottom(t *testing.T) {
	d := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}})
	d.PutBottom([]Card{fakeCard{id: 99}})
	if d.Size() != 4 {
		t.Fatalf("Size after PutBottom = %d, want 4", d.Size())
	}
	wantOrder := []ids.CardID{1, 2, 3, 99}
	for i, want := range wantOrder {
		got := d.Draw(1)[0]
		if got.ID() != want {
			t.Errorf("PutBottom card[%d] = %d, want %d", i, got.ID(), want)
		}
	}
}

// TestPutTop_PrependsToTop confirms PutTop puts cards at the top of the deck in the order
// passed (first card of the slice is the new top), preserving the existing bottom.
func TestPutTop_PrependsToTop(t *testing.T) {
	d := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}})
	d.PutTop([]Card{fakeCard{id: 99}, fakeCard{id: 98}})
	if d.Size() != 5 {
		t.Fatalf("Size after PutTop = %d, want 5", d.Size())
	}
	wantOrder := []ids.CardID{99, 98, 1, 2, 3}
	for i, want := range wantOrder {
		got := d.Draw(1)[0]
		if got.ID() != want {
			t.Errorf("PutTop card[%d] = %d, want %d", i, got.ID(), want)
		}
	}
}

// TestTutor_RemovesHighestScoring confirms Tutor scans the whole deck, removes the card
// with the highest positive score, and returns it. Cards scoring zero or negative are
// ignored; a deck where every card scores zero returns (nil, false) without removing.
func TestTutor_RemovesHighestScoring(t *testing.T) {
	d := New(nil, nil, []Card{
		fakeCard{id: 1}, fakeCard{id: 5}, fakeCard{id: 3}, fakeCard{id: 7},
	})
	score := func(c Card) int { return int(c.ID()) }
	got, ok := d.Tutor(score)
	if !ok || got.ID() != 7 {
		t.Errorf("Tutor returned %v ok=%v, want id 7 ok=true", got, ok)
	}
	if d.Size() != 3 {
		t.Errorf("Size after Tutor = %d, want 3", d.Size())
	}
	for _, c := range d.cards {
		if c.ID() == 7 {
			t.Errorf("Tutor left id 7 in the deck: %v", d.cards)
		}
	}
}

// TestTutor_NoMatchReturnsFalse confirms Tutor leaves the deck untouched when no card
// scores > 0.
func TestTutor_NoMatchReturnsFalse(t *testing.T) {
	d := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}})
	beforeSize := d.Size()
	got, ok := d.Tutor(func(Card) int { return 0 })
	if ok || got != nil {
		t.Errorf("Tutor on no-match returned %v ok=%v, want nil false", got, ok)
	}
	if d.Size() != beforeSize {
		t.Errorf("Size mutated by failed Tutor: %d, want %d", d.Size(), beforeSize)
	}
}

// TestShallowCopy_ShufflePanics confirms the safety net on the ShallowCopy optimization:
// a card calling Shuffle mid-turn on a per-permutation deck would silently corrupt every
// peer permutation sharing the same slice backing, so ShallowCopy-produced wrappers must
// panic on Shuffle. Any attack-turn runner test that exercises a card whose Play() shuffles
// the deck mid-turn will trip this panic.
func TestShallowCopy_ShufflePanics(t *testing.T) {
	master := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}})
	shallow := master.ShallowCopy()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Shuffle on a ShallowCopy-produced wrapper did not panic")
		}
	}()
	shallow.Shuffle(rand.New(rand.NewSource(1)))
}

// TestDeck_ShuffleOnMasterStillWorks confirms the panic is gated on the ShallowCopy
// marker — Shuffle on a fresh / Copy()'d deck behaves normally.
func TestDeck_ShuffleOnMasterStillWorks(t *testing.T) {
	d := New(nil, nil, []Card{fakeCard{id: 1}, fakeCard{id: 2}, fakeCard{id: 3}})
	d.Shuffle(rand.New(rand.NewSource(1))) // must not panic
	if d.Size() != 3 {
		t.Errorf("post-Shuffle size = %d, want 3", d.Size())
	}
}

// TestDeck_WithoutAndCount checks Count tallies copies and Without drops every copy of an id while
// preserving order, carrying Sideboard / Equipment, and leaving the original deck untouched.
func TestDeck_WithoutAndCount(t *testing.T) {
	d := New(nil, nil, []Card{
		fakeCard{id: 1, display: "A"},
		fakeCard{id: 2, display: "B"},
		fakeCard{id: 1, display: "A"},
		fakeCard{id: 3, display: "C"},
	})
	d.Sideboard = []string{"SB"}
	d.Equipment = []string{"EQ"}

	if got := d.Count(1); got != 2 {
		t.Errorf("Count(1) = %d, want 2", got)
	}
	if got := d.Count(3); got != 1 {
		t.Errorf("Count(3) = %d, want 1", got)
	}
	if got := d.Count(99); got != 0 {
		t.Errorf("Count(99) = %d, want 0 (absent card)", got)
	}

	w := d.Without(1)
	if w.Size() != 2 {
		t.Errorf("Without(1).Size() = %d, want 2 (both copies of 1 gone)", w.Size())
	}
	if w.Count(1) != 0 {
		t.Errorf("Without(1) kept %d copies of card 1", w.Count(1))
	}
	if order := []ids.CardID{w.cards[0].ID(), w.cards[1].ID()}; !slices.Equal(order, []ids.CardID{2, 3}) {
		t.Errorf("Without(1) surviving order = %v, want [2 3]", order)
	}
	if len(w.Sideboard) != 1 || w.Sideboard[0] != "SB" || len(w.Equipment) != 1 || w.Equipment[0] != "EQ" {
		t.Errorf("Without dropped sideboard/equipment: %v / %v", w.Sideboard, w.Equipment)
	}
	if d.Size() != 4 {
		t.Errorf("Without mutated the original deck; Size() = %d, want 4", d.Size())
	}
}
