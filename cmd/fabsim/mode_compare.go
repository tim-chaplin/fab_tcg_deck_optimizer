package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// runCompareCmd parses compare's flags (none today) and dispatches to runCompare. Both decks
// are positional args — compare never creates a deck, so there's no analogue of anneal's
// -deck checkpoint flag.
func runCompareCmd(args []string) {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: fabsim compare <deck1> <deck2>")
	}
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 2 {
		die("compare: need exactly 2 positional deck names (got %d); try `fabsim compare <deck1> <deck2>`", fs.NArg())
	}
	runCompare(fs.Arg(0), fs.Arg(1))
}

// runCompare loads the two decks at mydecks/<name1>.json and mydecks/<name2>.json and prints
// a stat-by-stat side-by-side comparison: pitch counts, mean hand value, per-cycle means, the
// hand-value histograms, and finally the per-card count delta.
func runCompare(name1, name2 string) {
	d1 := mustLoadDeck(resolveDeckPath(name1))
	d2 := mustLoadDeck(resolveDeckPath(name2))
	s1, s2 := d1.Stats, d2.Stats

	printSideBySideStats(name1, name2, []statSection{
		{"Pitch values", pitchCountsLine(d1.Cards), pitchCountsLine(d2.Cards)},
		{"Mean hand value", meanValueLine(s1), meanValueLine(s2)},
		{"Cycle 1 mean", cycleMeanLine(s1.FirstCycle), cycleMeanLine(s2.FirstCycle)},
		{"Cycle 2 mean", cycleMeanLine(s1.SecondCycle), cycleMeanLine(s2.SecondCycle)},
	})

	if len(s1.Histogram) > 0 || len(s2.Histogram) > 0 {
		fmt.Println()
		fmt.Println("Hand-value distributions:")
		if len(s1.Histogram) > 0 {
			printHistogram(d1, fmt.Sprintf("  %s (%s hands):", name1, commaInt(s1.Hands)))
		}
		if len(s2.Histogram) > 0 {
			printHistogram(d2, fmt.Sprintf("  %s (%s hands):", name2, commaInt(s2.Hands)))
		}
	}

	fmt.Println()
	printCardDelta(name1, name2, d1, d2)
}

// printCardDelta writes the per-card count delta between d1 and d2 — negative rows first,
// then positives; alphabetical within each group; cards present in equal counts in both decks
// are omitted. When the two card lists match exactly an explicit confirmation line replaces
// the empty body so silence can't be mistaken for a failure. Scoped to the card list — hero
// and weapon differences are not surfaced here.
func printCardDelta(name1, name2 string, d1, d2 *deck.Deck) {
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
	fmt.Println("Card differences:")
	if len(minuses) == 0 && len(pluses) == 0 {
		fmt.Printf("  %s and %s have identical card lists (%d cards)\n", name1, name2, len(d1.Cards))
		return
	}
	for _, l := range minuses {
		fmt.Println("  " + l)
	}
	for _, l := range pluses {
		fmt.Println("  " + l)
	}
}
