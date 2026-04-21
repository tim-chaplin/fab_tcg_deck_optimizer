package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
)

// runDiff loads the two decks at mydecks/<name1>.json and mydecks/<name2>.json and prints the
// per-card count delta from deck1 to deck2, one line per changed card. Negative rows first,
// then positives; alphabetical within each group. Cards present in equal counts in both decks
// are omitted.
func runDiff(name1, name2 string) {
	p1, err := mydecks.Path(name1)
	if err != nil {
		die("%v", err)
	}
	p2, err := mydecks.Path(name2)
	if err != nil {
		die("%v", err)
	}
	d1, _ := loadExisting(p1)
	if d1 == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", p1)
		os.Exit(1)
	}
	d2, _ := loadExisting(p2)
	if d2 == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", p2)
		os.Exit(1)
	}

	counts1 := map[string]int{}
	for _, c := range d1.Cards {
		counts1[c.Name()]++
	}
	counts2 := map[string]int{}
	for _, c := range d2.Cards {
		counts2[c.Name()]++
	}
	names := make(map[string]struct{}, len(counts1)+len(counts2))
	for n := range counts1 {
		names[n] = struct{}{}
	}
	for n := range counts2 {
		names[n] = struct{}{}
	}
	sorted := make([]string, 0, len(names))
	for n := range names {
		sorted = append(sorted, n)
	}
	sort.Strings(sorted)

	var minuses, pluses []string
	for _, n := range sorted {
		d := counts2[n] - counts1[n]
		switch {
		case d < 0:
			minuses = append(minuses, fmt.Sprintf("%d %s", d, n))
		case d > 0:
			pluses = append(pluses, fmt.Sprintf("+%d %s", d, n))
		}
	}
	for _, l := range minuses {
		fmt.Println(l)
	}
	for _, l := range pluses {
		fmt.Println(l)
	}
}
