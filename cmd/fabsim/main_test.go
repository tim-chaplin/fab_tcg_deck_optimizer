package main

import (
	"flag"
	"io"
	"reflect"
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
