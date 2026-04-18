package main

import (
	"testing"

	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestFabraryPathFor pins the sibling-path derivation. A .json extension is replaced; anything
// else gets .txt appended instead of clobbered so the JSON can't accidentally be overwritten.
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
