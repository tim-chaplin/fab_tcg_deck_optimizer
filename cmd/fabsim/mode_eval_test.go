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

// TestRunEval_DefaultRewritesFile pins the default eval contract: the sim always runs and
// the on-disk .json / .txt are overwritten with the fresh stats. That way a saved deck's
// avg stays in sync with the current binary's modelling without the user having to
// remember an opt-in flag.
func TestRunEval_DefaultRewritesFile(t *testing.T) {
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

	// Default eval (printOnly=false) MUST rewrite: the sim produced fresh stats (Hands>0)
	// where the seed had Hands=0. captureEvalOutput drains stdout/stderr so the test isn't
	// cluttered with eval's prints and runEval doesn't block on a full pipe buffer.
	_, _ = captureEvalOutput(t, func() {
		runEval(path, 50, 0, 2, 1, fmtpkg.SilverAge, false, true)
	})
	afterDefault, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after eval: %v", err)
	}
	if bytes.Equal(afterDefault, beforeJSON) {
		t.Errorf("default eval should have rewritten the file, but bytes match the seed (len=%d)",
			len(beforeJSON))
	}

	// The sibling .txt must also be refreshed (writeDeck writes both).
	txtPath := fabraryPathFor(path)
	if _, err := os.Stat(txtPath); err != nil {
		t.Errorf("default eval should have written the sibling .txt (%s): %v", txtPath, err)
	}
}

// TestRunEval_PrintOnlyLeavesFileUnchanged pins -print-only's contract: the sim is skipped
// and neither the .json nor the .txt is touched. Used for a quick inspection of saved stats
// without burning shuffles or mutating the file.
func TestRunEval_PrintOnlyLeavesFileUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")

	// Seed with real stats so the printed summary has something to show. A re-sim would
	// overwrite Hands>0 with a different number; -print-only must preserve the seeded file
	// byte-for-byte.
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	d.Evaluate(20, 0, rng)
	if err := writeDeck(d, path); err != nil {
		t.Fatalf("seed writeDeck: %v", err)
	}
	beforeJSON, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seeded json: %v", err)
	}

	_, stderr := captureEvalOutput(t, func() {
		runEval(path, 50, 0, 2, 1, fmtpkg.SilverAge, true, true)
	})
	afterRead, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after print-only: %v", err)
	}
	if !bytes.Equal(afterRead, beforeJSON) {
		t.Errorf("-print-only should leave the file unchanged; len before=%d after=%d",
			len(beforeJSON), len(afterRead))
	}
	// -print-only must not emit the delta-and-rewrite banner the simulate path uses. Unrelated
	// stderr is still allowed (e.g. a sanitize warning) so the test only pins the absence of
	// the rewrite-line marker.
	if strings.Contains(stderr, "rewriting") {
		t.Errorf("-print-only should not log a rewrite; got stderr:\n%s", stderr)
	}
}

// TestRunEval_DefaultPrintsFullDump captures runEval's stdout in the default (brief=false)
// mode and asserts the full printBestDeck shape: summary, card list, peak-Value best turn,
// per-card stats. This is the output a user sees when running `fabsim eval <deck> -incoming
// N` without flags. Default eval also re-simulates and rewrites, so stderr carries the
// delta-and-rewrite line.
func TestRunEval_DefaultPrintsFullDump(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")

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
	if !strings.Contains(stdout, "Card list:") {
		t.Errorf("default eval output missing the card list:\n%s", stdout)
	}
	if !strings.Contains(stderr, "rewriting") {
		t.Errorf("default eval should announce the rewrite on stderr; got:\n%s", stderr)
	}
}

// TestRunEval_BriefSkipsBestTurnAndCardList pins the -brief flag: the score summary is the
// only stdout. Default eval also re-simulates and rewrites, so stderr still carries the
// rewrite line — brief only suppresses stdout's card-list/best-turn blocks.
func TestRunEval_BriefSkipsBestTurnAndCardList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")

	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	if err := writeDeck(d, path); err != nil {
		t.Fatalf("seed writeDeck: %v", err)
	}

	stdout, _ := captureEvalOutput(t, func() {
		runEval(path, 100, 0, 2, 1, fmtpkg.SilverAge, false, true)
	})
	if !strings.Contains(stdout, "Mean value:") {
		t.Errorf("brief eval output missing the 'Mean value:' stats line:\n%s", stdout)
	}
	if strings.Contains(stdout, "Best turn played") {
		t.Errorf("brief eval should suppress the best-turn block; got:\n%s", stdout)
	}
	if strings.Contains(stdout, "Card list:") {
		t.Errorf("brief eval should suppress the card list; got:\n%s", stdout)
	}
}
