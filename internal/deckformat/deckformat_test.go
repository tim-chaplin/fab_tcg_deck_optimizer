package deckformat

import (
	"os"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in      string
		want    Format
		wantErr bool
	}{
		{"silver_age", SilverAge, false},
		{"", "", true},          // empty string is not a format — every run must specify one
		{"classic", "", true},   // unknown format name
		{"SilverAge", "", true}, // case-sensitive — compared against exact flag tokens
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("Parse(%q) = %q, nil; want error", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q): unexpected error %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("Parse(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestSilverAgeBanlistParity pins the banlist file against the code tags: every card in
// data_sources/silver_age_banlist.txt that has an implementation must tag every variant with
// NotSilverAgeLegal. Names are matched after stripping the color suffix and normalising
// curly apostrophes (U+2019) to ASCII.
func TestSilverAgeBanlistParity(t *testing.T) {
	// Banlist lives at the repo root; tests run from this package dir, so two levels up.
	data, err := os.ReadFile("../../data_sources/silver_age_banlist.txt")
	if err != nil {
		t.Fatalf("read banlist: %v", err)
	}

	banned := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		name := strings.TrimSpace(line)
		if name == "" || strings.HasSuffix(name, ":") {
			continue
		}
		banned[normalizeName(name)] = true
	}
	if len(banned) == 0 {
		t.Fatal("banlist file parsed to zero entries")
	}

	for _, id := range registry.AllCards() {
		c := registry.GetCard(id)
		base := normalizeName(stripVariantSuffix(c.Name()))
		if !banned[base] {
			continue
		}
		if _, ok := c.(sim.NotSilverAgeLegal); !ok {
			t.Errorf("%s is on the Silver Age banlist but isn't tagged with NotSilverAgeLegal", c.Name())
		}
	}
}

// stripVariantSuffix drops the trailing color marker so variants collapse to one entry.
func stripVariantSuffix(name string) string {
	for _, suffix := range []string{" [R]", " [Y]", " [B]"} {
		if strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(name, suffix)
		}
	}
	return name
}

// normalizeName replaces Unicode curly apostrophes with ASCII so typographic quotes in the
// banlist file match Name()'s plain quotes.
func normalizeName(s string) string {
	return strings.ReplaceAll(s, "\u2019", "'")
}

// TestIsLegal uses one card with NotSilverAgeLegal and one without to confirm Silver Age
// rejects the banned one and accepts the other.
func TestIsLegal(t *testing.T) {
	banned := notimpl.BelittleRed{}
	legal := cards.NimblismRed{}

	// Sanity: the marker is present on the banned card and absent on the legal one. Guards
	// against accidentally dropping the tag.
	if _, ok := sim.Card(banned).(sim.NotSilverAgeLegal); !ok {
		t.Fatal("BelittleRed: missing NotSilverAgeLegal marker")
	}
	if _, ok := sim.Card(legal).(sim.NotSilverAgeLegal); ok {
		t.Fatal("NimblismRed: has NotSilverAgeLegal marker but shouldn't")
	}

	cases := []struct {
		f    Format
		c    sim.Card
		want bool
	}{
		{SilverAge, banned, false},
		{SilverAge, legal, true},
	}
	for _, tc := range cases {
		if got := tc.f.IsLegal(tc.c); got != tc.want {
			t.Errorf("Format(%q).IsLegal(%s) = %v, want %v", tc.f, tc.c.Name(), got, tc.want)
		}
	}
}
