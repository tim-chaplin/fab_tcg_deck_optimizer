package format

import (
	"os"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
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

// TestSilverAgeBanlistParity pins the authoritative banlist file against the code tags: every
// card name in data_sources/silver_age_banlist.txt that has an implementation must have every one
// of its variants tagged with NotSilverAgeLegal. Catches drift when someone edits the banlist
// without updating code, or implements an already-banned card without tagging it.
//
// Names are matched after stripping the " (Red)" / " (Yellow)" / " (Blue)" color suffix and
// normalising Unicode curly apostrophes (U+2019) to ASCII — the banlist file uses curly quotes
// in places where Name() returns straight ones.
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

	for _, id := range cards.All() {
		c := cards.Get(id)
		base := normalizeName(stripVariantSuffix(c.Name()))
		if !banned[base] {
			continue
		}
		if _, ok := c.(card.NotSilverAgeLegal); !ok {
			t.Errorf("%s is on the Silver Age banlist but isn't tagged with NotSilverAgeLegal", c.Name())
		}
	}
}

// stripVariantSuffix drops the trailing color marker so variants collapse to one entry.
func stripVariantSuffix(name string) string {
	for _, suffix := range []string{" (Red)", " (Yellow)", " (Blue)"} {
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

// TestIsLegal uses two concrete cards: Plunder Run (Red) which implements NotSilverAgeLegal, and
// Nimblism (Red) which doesn't. Silver Age rejects the banned one and accepts the other.
func TestIsLegal(t *testing.T) {
	banned := generic.PlunderRunRed{}
	legal := generic.NimblismRed{}

	// Sanity: the marker is present on the banned card and absent on the legal one. Guards
	// against accidentally dropping the tag.
	if _, ok := card.Card(banned).(card.NotSilverAgeLegal); !ok {
		t.Fatal("PlunderRunRed: missing NotSilverAgeLegal marker")
	}
	if _, ok := card.Card(legal).(card.NotSilverAgeLegal); ok {
		t.Fatal("NimblismRed: has NotSilverAgeLegal marker but shouldn't")
	}

	cases := []struct {
		f    Format
		c    card.Card
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
