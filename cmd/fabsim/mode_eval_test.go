package main

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// captureEvalOutput redirects os.Stdout / os.Stderr into pipes, drains them concurrently
// (so runEval's writes don't block on a full pipe buffer when printBestDeck emits the
// full card list + per-card stats), runs f, and returns the captured bytes. Stdout and
// stderr are restored via t.Cleanup before the function returns.
func captureEvalOutput(t *testing.T, f func()) (stdout, stderr string) {
	t.Helper()
	origStdout, origStderr := os.Stdout, os.Stderr
	rStdout, wStdout, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	rStderr, wStderr, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stderr: %v", err)
	}
	os.Stdout, os.Stderr = wStdout, wStderr
	t.Cleanup(func() { os.Stdout, os.Stderr = origStdout, origStderr })

	var outBuf, errBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); _, _ = io.Copy(&outBuf, rStdout) }()
	go func() { defer wg.Done(); _, _ = io.Copy(&errBuf, rStderr) }()

	f()

	if err := wStdout.Close(); err != nil {
		t.Fatalf("close stdout write side: %v", err)
	}
	if err := wStderr.Close(); err != nil {
		t.Fatalf("close stderr write side: %v", err)
	}
	wg.Wait()
	return outBuf.String(), errBuf.String()
}

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

	// First pass: reevaluate=false. File must NOT be rewritten (byte-identical to seed).
	// captureEvalOutput drains stdout/stderr so the test isn't cluttered with eval's prints
	// and runEval doesn't block on a full pipe buffer.
	_, _ = captureEvalOutput(t, func() {
		runEval(path, 50, 0, 2, 1, fmtpkg.SilverAge, false, true)
	})
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
	_, _ = captureEvalOutput(t, func() {
		runEval(path, 50, 0, 2, 1, fmtpkg.SilverAge, true, true)
	})
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

// TestRunEval_DefaultPrintsFullDump captures runEval's stdout in the default (brief=false)
// mode and asserts the full printBestDeck shape: summary, card list, peak-Value best turn,
// per-card stats. This is the output a user sees when running `fabsim eval <deck> -incoming
// N` without flags.
func TestRunEval_DefaultPrintsFullDump(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")

	// Viserai seed so the "carryover runechants" header tag is part of the expected output.
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	if err := writeDeck(d, path); err != nil {
		t.Fatalf("seed writeDeck: %v", err)
	}

	stdout, stderr := captureEvalOutput(t, func() {
		runEval(path, 100, 0, 2, 1, fmtpkg.SilverAge, false, false)
	})
	if !strings.Contains(stdout, "Best turn played") {
		t.Errorf("eval output missing 'Best turn played' header:\n%s", stdout)
	}
	if !strings.Contains(stdout, "carryover runechants") {
		t.Errorf("Viserai eval output missing 'carryover runechants' tag:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Card list:") {
		t.Errorf("default eval output missing the card list:\n%s", stdout)
	}
	if stderr != "" {
		t.Errorf("read-only eval should not write to stderr; got:\n%s", stderr)
	}
}

// TestRunEval_BriefSkipsBestTurnAndCardList pins the -brief flag: the score summary is the
// only output. Useful for scripted re-scoring where the extra blocks would just be noise.
func TestRunEval_BriefSkipsBestTurnAndCardList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")

	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	if err := writeDeck(d, path); err != nil {
		t.Fatalf("seed writeDeck: %v", err)
	}

	stdout, stderr := captureEvalOutput(t, func() {
		runEval(path, 100, 0, 2, 1, fmtpkg.SilverAge, false, true)
	})
	if !strings.Contains(stdout, "Best deck (min") {
		t.Errorf("brief eval output missing the score header:\n%s", stdout)
	}
	if strings.Contains(stdout, "Best turn played") {
		t.Errorf("brief eval should suppress the best-turn block; got:\n%s", stdout)
	}
	if strings.Contains(stdout, "Card list:") {
		t.Errorf("brief eval should suppress the card list; got:\n%s", stdout)
	}
	if stderr != "" {
		t.Errorf("read-only brief eval should not write to stderr; got:\n%s", stderr)
	}
}
