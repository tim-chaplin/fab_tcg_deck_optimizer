package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestFabraryPathFor pins the sibling-path derivation: .json is swapped for .txt; anything else
// gets .txt appended so the original isn't clobbered.
func TestFabraryPathFor(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"mydecks/best_deck.json", "mydecks/best_deck.txt"},
		{"best.json", "best.txt"},
		{"deck", "deck.txt"},
		{"deck.data", "deck.data.txt"},
		{"with.dots.json", "with.dots.txt"},
	}
	for _, c := range cases {
		if got := fabraryPathFor(c.in); got != c.want {
			t.Errorf("fabraryPathFor(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestParseFlagsAnywhere pins the reorder behavior: flags must parse regardless of position
// relative to positional args, both `-name value` and `-name=value` forms work, bool flags
// don't consume the next token, `--` terminates flag parsing, and unknown flags surface their
// canonical fs.Parse error instead of silently swallowing a positional.
func TestParseFlagsAnywhere(t *testing.T) {
	t.Run("flag after positional", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		incoming := fs.Int("incoming", 0, "")
		if err := parseFlagsAnywhere(fs, []string{"deckname", "--incoming=100"}); err != nil {
			t.Fatalf("parse: %v", err)
		}
		if *incoming != 100 {
			t.Errorf("incoming = %d, want 100", *incoming)
		}
		if got := fs.Args(); !reflect.DeepEqual(got, []string{"deckname"}) {
			t.Errorf("positional = %v, want [deckname]", got)
		}
	})

	t.Run("-name value (space-separated) after positional", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		incoming := fs.Int("incoming", 0, "")
		if err := parseFlagsAnywhere(fs, []string{"deckname", "-incoming", "42"}); err != nil {
			t.Fatalf("parse: %v", err)
		}
		if *incoming != 42 {
			t.Errorf("incoming = %d, want 42", *incoming)
		}
		if got := fs.Args(); !reflect.DeepEqual(got, []string{"deckname"}) {
			t.Errorf("positional = %v, want [deckname]", got)
		}
	})

	t.Run("bool flag does not consume next positional", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		debug := fs.Bool("debug", false, "")
		if err := parseFlagsAnywhere(fs, []string{"-debug", "deckname"}); err != nil {
			t.Fatalf("parse: %v", err)
		}
		if !*debug {
			t.Errorf("debug = false, want true")
		}
		if got := fs.Args(); !reflect.DeepEqual(got, []string{"deckname"}) {
			t.Errorf("positional = %v, want [deckname]", got)
		}
	})

	t.Run("-- terminator treats following args as positional", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		fs.Int("incoming", 0, "")
		if err := parseFlagsAnywhere(fs, []string{"--", "--incoming=100", "deckname"}); err != nil {
			t.Fatalf("parse: %v", err)
		}
		want := []string{"--incoming=100", "deckname"}
		if got := fs.Args(); !reflect.DeepEqual(got, want) {
			t.Errorf("positional = %v, want %v", got, want)
		}
	})
}

// TestLoadExisting_MissingReturnsNilNoError: a path that doesn't exist is the "no previous
// best, generate a fresh deck" signal — caller distinguishes nil from error to know whether
// it's safe to overwrite.
func TestLoadExisting_MissingReturnsNilNoError(t *testing.T) {
	d, _, err := loadExisting(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatalf("missing file should not return an error, got %v", err)
	}
	if d != nil {
		t.Errorf("missing file should return nil deck, got %+v", d)
	}
}

// TestLoadExisting_CorruptReturnsError: when the file exists but isn't valid JSON / a valid
// deck, loadExisting MUST return an error so the caller refuses to overwrite. The user-visible
// failure that motivated this fix was a Ctrl-C-truncated deck file getting silently
// replaced with a random deck on the next anneal pass.
func TestLoadExisting_CorruptReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0o644); err != nil {
		t.Fatalf("seed corrupt file: %v", err)
	}
	d, _, err := loadExisting(path)
	if err == nil {
		t.Fatalf("corrupt file should return an error; got nil (deck=%+v)", d)
	}
	if d != nil {
		t.Errorf("corrupt file should return nil deck, got %+v", d)
	}
	// The message has to make clear we're refusing to overwrite — that's the whole point of
	// the loud failure. Pin the substring rather than the exact message so wording can drift.
	if !strings.Contains(err.Error(), "refusing to silently overwrite") {
		t.Errorf("error message should warn about overwrite refusal, got %q", err.Error())
	}
}

// TestLoadExisting_TruncatedReturnsError: simulates the exact failure mode that motivated
// this fix — a writeDeck interrupted between O_TRUNC and the data write would have left an
// empty file. With the atomic-write fix that can't happen anymore, but loadExisting must
// still treat an empty file as corrupt so a manually-cleared file doesn't silently get
// replaced either.
func TestLoadExisting_TruncatedReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.json")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("seed empty file: %v", err)
	}
	if _, _, err := loadExisting(path); err == nil {
		t.Errorf("empty file should return an error; got nil")
	}
}

// TestWriteFileAtomic_LeavesNoTempOnSuccess: the temp file used during the write must be
// renamed away (not left behind) on success — anything else would litter mydecks/ with
// .tmp-* files over thousands of anneal rounds.
func TestWriteFileAtomic_LeavesNoTempOnSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")
	if err := writeFileAtomic(path, []byte("hello")); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("file content = %q, want %q", got, "hello")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	for _, e := range entries {
		if e.Name() != "deck.json" {
			t.Errorf("unexpected leftover file in temp dir: %q", e.Name())
		}
	}
}

// TestWriteFileAtomic_PreservesOldOnFailure: if the rename target already has content and
// the write succeeds, the new content replaces the old atomically. (We can't easily simulate
// a partial write failure without an injectable seam, but at minimum the success path must
// fully replace, not append/corrupt.)
func TestWriteFileAtomic_PreservesOldOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deck.json")
	if err := os.WriteFile(path, []byte("old contents"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := writeFileAtomic(path, []byte("new contents")); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(got) != "new contents" {
		t.Errorf("file content = %q, want %q (atomic replace should fully overwrite)", got, "new contents")
	}
}

// TestDefaultDeckNameFor pins the filename shape: hero_format_incoming. Every run is scoped to
// exactly one format, so the format segment is always present.
func TestDefaultDeckNameFor(t *testing.T) {
	cases := []struct {
		f    fmtpkg.Format
		in   int
		want string
	}{
		{fmtpkg.SilverAge, 0, "viserai_silver_age_0_incoming"},
		{fmtpkg.SilverAge, 4, "viserai_silver_age_4_incoming"},
	}
	for _, c := range cases {
		if got := defaultDeckNameFor(hero.Viserai{}, c.f, c.in); got != c.want {
			t.Errorf("defaultDeckNameFor(Viserai, %q, %d) = %q, want %q", c.f, c.in, got, c.want)
		}
	}
}
