package deck

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// nameCounts returns a multiset of names for comparison in the tests below. Using maps
// keeps assertions order-independent since ApplyDefaults appends to the existing slice.
func nameCounts(ss []string) map[string]int {
	m := map[string]int{}
	for _, s := range ss {
		m[s]++
	}
	return m
}

// TestApplyDefaults_FillsEmptyEquipment pins the equipment-merge base case: an empty
// Equipment list ends up with every defaultEquipment entry present once.
func TestApplyDefaults_FillsEmptyEquipment(t *testing.T) {
	d := New(hero.Viserai{}, nil, nil)
	d.ApplyDefaults()

	got := nameCounts(d.Equipment)
	for _, name := range defaultEquipment {
		if got[name] != 1 {
			t.Errorf("equipment[%q] = %d, want 1 (got full equipment: %v)", name, got[name], d.Equipment)
		}
	}
}

// TestApplyDefaults_KeepsUserEquipmentExtras pins that ApplyDefaults never removes or
// clamps user-supplied equipment — extra copies stay put and unknown names pass through.
func TestApplyDefaults_KeepsUserEquipmentExtras(t *testing.T) {
	d := New(hero.Viserai{}, nil, nil)
	d.Equipment = []string{"Beckoning Haunt", "Beckoning Haunt", "Custom Equipment Piece"}
	d.ApplyDefaults()

	got := nameCounts(d.Equipment)
	if got["Beckoning Haunt"] != 2 {
		t.Errorf("Beckoning Haunt count = %d, want 2 (user's two copies preserved)", got["Beckoning Haunt"])
	}
	if got["Custom Equipment Piece"] != 1 {
		t.Errorf("Custom Equipment Piece dropped (count = %d)", got["Custom Equipment Piece"])
	}
	// The other defaults still land even though the user listed their own.
	if got["Blossom of Spring"] != 1 {
		t.Errorf("Blossom of Spring count = %d, want 1", got["Blossom of Spring"])
	}
}

// TestApplyDefaults_Idempotent pins that running ApplyDefaults twice is a no-op. An
// idempotent merge matters because writeDeck runs it on every save, and re-loading +
// re-saving shouldn't accumulate copies.
func TestApplyDefaults_Idempotent(t *testing.T) {
	d := New(hero.Viserai{}, nil, nil)
	d.ApplyDefaults()
	firstEquipment := append([]string(nil), d.Equipment...)
	firstSideboard := append([]string(nil), d.Sideboard...)
	d.ApplyDefaults()

	if !reflect.DeepEqual(nameCounts(d.Equipment), nameCounts(firstEquipment)) {
		t.Errorf("equipment changed on second call:\n  first:  %v\n  second: %v", firstEquipment, d.Equipment)
	}
	if !reflect.DeepEqual(nameCounts(d.Sideboard), nameCounts(firstSideboard)) {
		t.Errorf("sideboard changed on second call:\n  first:  %v\n  second: %v", firstSideboard, d.Sideboard)
	}
}

// TestApplyDefaults_SideboardRespectsCopyCap pins the main-deck interaction: when the main
// deck already holds copies of a default sideboard entry, ApplyDefaults tops the sideboard
// up only far enough to keep the combined total at or below sideboardCopyCap.
func TestApplyDefaults_SideboardRespectsCopyCap(t *testing.T) {
	readRunes := cards.Get(card.ReadTheRunesRed)

	// Main has 2 copies → sideboard should stay empty for this entry.
	d := New(hero.Viserai{}, nil, []card.Card{readRunes, readRunes})
	d.ApplyDefaults()
	if nameCounts(d.Sideboard)["Read the Runes (Red)"] != 0 {
		t.Errorf("Read the Runes already at cap in main deck; sideboard should not add any. Got sideboard: %v", d.Sideboard)
	}

	// Main has 1 copy → sideboard gets 1 (total 2).
	d = New(hero.Viserai{}, nil, []card.Card{readRunes})
	d.ApplyDefaults()
	if got := nameCounts(d.Sideboard)["Read the Runes (Red)"]; got != 1 {
		t.Errorf("sideboard Read the Runes count = %d, want 1 (topping up to 2 total)", got)
	}

	// Main has 0 copies → sideboard gets the entry's full target (2).
	d = New(hero.Viserai{}, nil, nil)
	d.ApplyDefaults()
	if got := nameCounts(d.Sideboard)["Read the Runes (Red)"]; got != 2 {
		t.Errorf("sideboard Read the Runes count = %d, want 2", got)
	}
}

// TestApplyDefaults_SideboardKeepsUserExtras pins that a user who explicitly added more
// than the default target keeps their copies — the merge only tops up, never clamps down.
func TestApplyDefaults_SideboardKeepsUserExtras(t *testing.T) {
	d := New(hero.Viserai{}, nil, nil)
	// Crown of Dichotomy's default target is 1; start with 3 user-supplied copies.
	d.Sideboard = []string{"Crown of Dichotomy", "Crown of Dichotomy", "Crown of Dichotomy"}
	d.ApplyDefaults()

	if got := nameCounts(d.Sideboard)["Crown of Dichotomy"]; got != 3 {
		t.Errorf("Crown of Dichotomy count = %d, want 3 (user's extras preserved)", got)
	}
}
