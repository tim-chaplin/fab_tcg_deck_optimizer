package deck

// Composition-level formatting helpers for Deck. Lives in deck/ so the helpers can read
// d.cards directly — printing and marshalling never need the runtime deck order, only the
// composition multiset / display-name surface.

import (
	"fmt"
	"strings"
)

// NameCounts returns a DisplayName-keyed multiset of the deck's cards. Marshallers and
// diff tools work against the composition (not runtime order), so the multiset shape
// matches what they emit.
func (d *Deck) NameCounts() map[string]int {
	if d == nil {
		return nil
	}
	m := make(map[string]int, len(d.cards))
	for _, c := range d.cards {
		m[c.DisplayName()]++
	}
	return m
}

// NameCountsTransformed is NameCounts with each DisplayName run through transform before
// keying. Lets fabrary lowercase its [r]/[y]/[b] pitch suffix without re-walking the slice.
func (d *Deck) NameCountsTransformed(transform func(string) string) map[string]int {
	if d == nil {
		return nil
	}
	m := make(map[string]int, len(d.cards))
	for _, c := range d.cards {
		m[transform(c.DisplayName())]++
	}
	return m
}

// DisplayNames returns the deck's card DisplayNames in deck order. Used by printers
// that group-and-sort a flat name list (e.g. "Card list:" output).
func (d *Deck) DisplayNames() []string {
	if d == nil {
		return nil
	}
	out := make([]string, len(d.cards))
	for i, c := range d.cards {
		out[i] = c.DisplayName()
	}
	return out
}

// MaxNameLen returns the length of the longest DisplayName in the deck, or 0 when empty.
// Used to size fixed-width name columns in printed tables.
func (d *Deck) MaxNameLen() int {
	if d == nil {
		return 0
	}
	m := 0
	for _, c := range d.cards {
		if n := len(c.DisplayName()); n > m {
			m = n
		}
	}
	return m
}

// PitchCounts returns red/yellow/blue copy counts derived from the trailing
// "[R]"/"[Y]"/"[B]" suffix on each card's DisplayName. Cards without a recognised
// pitch tag contribute to no bucket.
func (d *Deck) PitchCounts() (red, yellow, blue int) {
	if d == nil {
		return 0, 0, 0
	}
	for _, c := range d.cards {
		switch displayPitchSuffix(c.DisplayName()) {
		case "R":
			red++
		case "Y":
			yellow++
		case "B":
			blue++
		}
	}
	return red, yellow, blue
}

// PitchCountsLine renders PitchCounts as the canonical "20 red / 8 yellow / 12 blue" line
// printers use under a deck's "Pitch:" header.
func (d *Deck) PitchCountsLine() string {
	red, yellow, blue := d.PitchCounts()
	return fmt.Sprintf("%d red / %d yellow / %d blue", red, yellow, blue)
}

// displayPitchSuffix extracts the "R"/"Y"/"B" character inside the trailing "[X]" tag of
// a DisplayName, or "" when the suffix doesn't match the pitch-tag shape.
func displayPitchSuffix(name string) string {
	if !strings.HasSuffix(name, "]") || len(name) < 4 {
		return ""
	}
	if name[len(name)-3] != '[' {
		return ""
	}
	return string(name[len(name)-2])
}
