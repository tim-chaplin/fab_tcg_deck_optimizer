package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
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
// swap as a -/+ pair so a loadout diff that lives only in the weapon list isn't silently
// collapsed into the "identical card lists" branch.
func TestPrintCardDelta_IncludesWeaponDifferences(t *testing.T) {
	cs := []sim.Card{registry.GetCard(ids.ReadTheRunesRed)}
	d1 := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.NebulaBlade{}}, cs)
	d2 := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, cs)

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

// TestPrintCardDelta_WeaponsLeadEachBlock pins the within-block order: weapons sit at the
// top of the minus block and at the top of the plus block, ahead of any card lines, so the
// loadout-defining piece is the first thing the reader sees in each direction.
func TestPrintCardDelta_WeaponsLeadEachBlock(t *testing.T) {
	read := registry.GetCard(ids.ReadTheRunesRed)
	snatch := registry.GetCard(ids.SnatchRed)
	d1 := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.NebulaBlade{}}, []sim.Card{read, read})
	d2 := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, []sim.Card{snatch, snatch})

	out := captureStdout(t, func() { printCardDelta("d1", "d2", d1, d2) })

	// Walk the body lines in printed order and assert: weapon minus precedes card minus,
	// then weapon plus precedes card plus.
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	wantOrder := []string{
		"  -1 Nebula Blade",
		"  -2 Read the Runes [R]",
		"  +1 Reaping Blade",
		"  +2 Snatch [R]",
	}
	got := make([]string, 0, len(wantOrder))
	for _, l := range lines {
		if strings.HasPrefix(l, "  -") || strings.HasPrefix(l, "  +") {
			got = append(got, l)
		}
	}
	if len(got) != len(wantOrder) {
		t.Fatalf("delta line count: got %d want %d; full output:\n%s", len(got), len(wantOrder), out)
	}
	for i, want := range wantOrder {
		if got[i] != want {
			t.Errorf("line %d: got %q want %q; full output:\n%s", i, got[i], want, out)
		}
	}
}

// TestPrintCardDelta_IdenticalLoadoutNotesBothCounts pins the empty-diff confirmation: when
// both decks have the same cards AND the same weapons, the line reports both totals so the
// reader knows the comparison covered weapons and not just cards.
func TestPrintCardDelta_IdenticalLoadoutNotesBothCounts(t *testing.T) {
	cs := []sim.Card{registry.GetCard(ids.ReadTheRunesRed), registry.GetCard(ids.SnatchRed)}
	weps := []sim.Weapon{weapons.NebulaBlade{}}
	d1 := sim.New(heroes.Viserai{}, weps, cs)
	d2 := sim.New(heroes.Viserai{}, weps, cs)

	out := captureStdout(t, func() { printCardDelta("d1", "d2", d1, d2) })

	if !strings.Contains(out, "identical card and weapon lists (2 cards, 1 weapons)") {
		t.Errorf("expected the empty-diff confirmation to report both counts; got:\n%s", out)
	}
}
