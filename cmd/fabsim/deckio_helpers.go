package main

// Deck load / save plumbing: atomic JSON + fabrary .txt writes, must-load variants that die
// on failure, path resolution, and the NotImplemented-sanitization wrapper that prints the
// per-swap warning lines. Every subcommand routes deck persistence through these helpers so
// the on-disk format and error messages stay uniform.

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
)

// loadExisting reads and deserializes the deck at path. Returns (nil, 0, nil) when the file
// doesn't exist — the caller treats that as "no previous best, generate a fresh deck."
// Returns (nil, 0, err) when the file exists but can't be read or parsed: callers must NOT
// treat that as "missing" or they'd silently overwrite a corrupt file with a random deck
// (looping wrapper scripts would clobber a converged deck after a Ctrl-C mid-write).
func loadExisting(path string) (*deck.Deck, float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("read %s: %w", path, err)
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		return nil, 0, fmt.Errorf("parse %s: %w (file exists but isn't a valid deck — "+
			"refusing to silently overwrite; inspect the file and delete it manually if you "+
			"want a fresh start)", path, err)
	}
	return d, d.Stats.Mean(), nil
}

// writeDeck persists d as JSON at path plus a sibling fabrary-format .txt ("x.json" →
// "x.txt") so the saved deck is ready to paste into fabrary.net without a second export step.
//
// ApplyDefaults runs first so both files carry the hardcoded default equipment / sideboard
// loadout the user runs on every deck. Persisting the defaults into the JSON (not just the
// .txt) means a reloaded deck already has them in Equipment / Sideboard without another
// round trip.
//
// Both files are written atomically via writeFileAtomic: data lands in <path>.tmp first,
// then os.Rename swaps it into place, so a Ctrl-C mid-write can never leave the destination
// empty or partially written.
func writeDeck(d *deck.Deck, path string) error {
	d.ApplyDefaults()
	data, err := deckio.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := writeFileAtomic(path, data); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	txtPath := fabraryPathFor(path)
	if err := writeFileAtomic(txtPath, []byte(fabrary.Marshal(d))); err != nil {
		return fmt.Errorf("write %s: %w", txtPath, err)
	}
	return nil
}

// writeFileAtomic writes data to a temp file in the same directory as path and renames it
// over path. The same-directory placement keeps the rename within one filesystem so it stays
// atomic. Removes the temp file on any error so a failed write doesn't leave junk behind.
func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	// Clean up the temp file on any failure path so crashed writes don't litter mydecks/ with
	// .tmp-* files. The rename success path makes this a no-op.
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// fabraryPathFor derives the sibling .txt path. A ".json" extension is swapped for ".txt";
// anything else gets ".txt" appended so non-JSON paths can't be overwritten.
func fabraryPathFor(jsonPath string) string {
	if ext := filepath.Ext(jsonPath); ext == ".json" {
		return strings.TrimSuffix(jsonPath, ext) + ".txt"
	}
	return jsonPath + ".txt"
}

// sanitizeLoadedDeck swaps every card.NotImplemented copy in d for a random legal
// replacement, prints a warning summary on stderr when any swap was made, and returns the
// ordered list of swaps. maxCopies caps post-sanitize copies per printing; legal restricts
// the replacement pool (typically the run's format predicate). Returns nil when the deck
// was already clean — callers can use that to skip the forced-reevaluation branch.
//
// The sanitizer mutates d.Cards in place. Callers that care about the pre-sanitize score
// for a delta warning should capture it before calling this.
func sanitizeLoadedDeck(d *deck.Deck, maxCopies int, rng *rand.Rand, legal func(card.Card) bool) []deck.NotImplementedReplacement {
	replaced := d.SanitizeNotImplemented(maxCopies, rng, legal)
	if len(replaced) == 0 {
		return nil
	}
	fmt.Fprintf(os.Stderr, "warning: loaded deck contained %d NotImplemented card(s); replacing with legal substitutes:\n", len(replaced))
	for _, r := range replaced {
		fmt.Fprintf(os.Stderr, "  -1 %s, +1 %s\n", r.From.Name(), r.To.Name())
	}
	return replaced
}

// mustLoadDeck loads the deck at path or dies. For subcommands that always operate on an
// existing deck (eval, diff), both "missing" and "corrupt" are fatal. anneal handles the
// distinction itself: "missing" is a valid input ("no deck yet, generate one") while
// "corrupt" needs the loud refusal to overwrite.
func mustLoadDeck(path string) *deck.Deck {
	d, _, err := loadExisting(path)
	if err != nil {
		die("%v", err)
	}
	if d == nil {
		die("could not load deck from %s (file not found)", path)
	}
	return d
}

// resolveDeckPath is the positional-arg counterpart to anneal's -deck flag. Subcommands that
// always operate on an existing deck (eval, diff) accept the deck name as a positional arg
// and resolve it to mydecks/<name>.json via mydecks.Path.
func resolveDeckPath(name string) string {
	p, err := mydecks.Path(name)
	if err != nil {
		die("%v", err)
	}
	return p
}
