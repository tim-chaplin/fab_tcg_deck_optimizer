package main

import (
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestRunEval_ReevaluateRewritesFile pins the -reevaluate flag's contract: passing it
// forces a fresh writeDeck so the on-disk stats (and any new best-turn fields introduced
// since the deck was last saved) catch up to the current binary's output. Without the
// flag, eval is read-only when the deck doesn't need sanitization.
func TestRunEval_ReevaluateRewritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")

	// Seed a deck with a stale Stats.Avg the fresh sim is guaranteed to overwrite. 40
	// random Viserai cards gives the sim enough to produce non-zero Value.
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	d.Stats.TotalValue = 0
	d.Stats.Hands = 0
	if err := writeDeck(d, path); err != nil {
		t.Fatalf("seed writeDeck: %v", err)
	}

	beforeJSON, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seeded json: %v", err)
	}

	// Silence runEval's stdout/stderr so the test output stays clean.
	origStdout, origStderr := os.Stdout, os.Stderr
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open devnull: %v", err)
	}
	defer devNull.Close()
	os.Stdout, os.Stderr = devNull, devNull
	t.Cleanup(func() { os.Stdout, os.Stderr = origStdout, origStderr })

	// First pass: reevaluate=false. File must NOT be rewritten (byte-identical to seed).
	runEval(path, 50, 0, 2, 1, fmtpkg.SilverAge, false)
	afterRead, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after eval: %v", err)
	}
	if !bytes.Equal(afterRead, beforeJSON) {
		t.Errorf("eval without -reevaluate should leave the file unchanged; len before=%d after=%d",
			len(beforeJSON), len(afterRead))
	}

	// Second pass: reevaluate=true. File MUST be rewritten — the sim produced fresh
	// stats (Hands>0) where the seed had Hands=0.
	runEval(path, 50, 0, 2, 1, fmtpkg.SilverAge, true)
	afterReevaluate, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after reevaluate: %v", err)
	}
	if bytes.Equal(afterReevaluate, beforeJSON) {
		t.Errorf("eval with -reevaluate should have rewritten the file, but bytes match the seed (len=%d)",
			len(beforeJSON))
	}

	// The sibling .txt must also be refreshed (writeDeck writes both).
	txtPath := fabraryPathFor(path)
	if _, err := os.Stat(txtPath); err != nil {
		t.Errorf("reevaluate should have written the sibling .txt (%s): %v", txtPath, err)
	}
}

