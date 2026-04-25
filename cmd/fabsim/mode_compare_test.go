package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// captureStdout drains os.Stdout into a buffer for the duration of f and restores it after.
// Used by the compare-mode tests where the assertion is on the printed text rather than on
// any return value. Mirrors captureEvalOutput in mode_eval_test.go but only handles stdout
// since the compare-mode print helpers don't write to stderr.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = orig })

	var buf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); _, _ = io.Copy(&buf, r) }()

	f()
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe write: %v", err)
	}
	wg.Wait()
	return buf.String()
}

// TestPrintCardDelta_IncludesWeaponDifferences pins that printCardDelta surfaces a weapon
// swap as a -/+ pair the same way it does for a card swap. Identical card lists with
// different weapon loadouts used to render as "identical card lists" — that hid a real
// loadout difference from the reader.
func TestPrintCardDelta_IncludesWeaponDifferences(t *testing.T) {
	cs := []card.Card{cards.Get(card.ReadTheRunesRed)}
	d1 := deck.New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, cs)
	d2 := deck.New(hero.Viserai{}, []weapon.Weapon{weapon.ReapingBlade{}}, cs)

	out := captureStdout(t, func() { printCardDelta("d1", "d2", d1, d2) })

	if !strings.Contains(out, "-1 Nebula Blade") {
		t.Errorf("expected '-1 Nebula Blade' in delta; got:\n%s", out)
	}
	if !strings.Contains(out, "+1 Reaping Blade") {
		t.Errorf("expected '+1 Reaping Blade' in delta; got:\n%s", out)
	}
	if strings.Contains(out, "identical card and weapon lists") {
		t.Errorf("decks differ on weapons but printCardDelta declared them identical:\n%s", out)
	}
}

// TestPrintCardDelta_IdenticalLoadoutNotesBothCounts pins the empty-diff confirmation: when
// both decks have the same cards AND the same weapons, the line reports both totals so the
// reader knows the comparison covered weapons and not just cards.
func TestPrintCardDelta_IdenticalLoadoutNotesBothCounts(t *testing.T) {
	cs := []card.Card{cards.Get(card.ReadTheRunesRed), cards.Get(card.SnatchRed)}
	weps := []weapon.Weapon{weapon.NebulaBlade{}}
	d1 := deck.New(hero.Viserai{}, weps, cs)
	d2 := deck.New(hero.Viserai{}, weps, cs)

	out := captureStdout(t, func() { printCardDelta("d1", "d2", d1, d2) })

	if !strings.Contains(out, "identical card and weapon lists (2 cards, 1 weapons)") {
		t.Errorf("expected the empty-diff confirmation to report both counts; got:\n%s", out)
	}
}
