package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

// runDiffCmd parses diff's flags (none today) and dispatches to runDiff. Both decks are
// positional args — diff never creates a deck, so there's no analogue of anneal's -deck
// checkpoint flag.
func runDiffCmd(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: fabsim diff <deck1> <deck2>")
	}
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 2 {
		die("diff: need exactly 2 positional deck names (got %d); try `fabsim diff <deck1> <deck2>`", fs.NArg())
	}
	runDiff(fs.Arg(0), fs.Arg(1))
}

// runDiff loads the two decks at mydecks/<name1>.json and mydecks/<name2>.json and prints the
// per-card count delta from deck1 to deck2, one line per changed card. Negative rows first,
// then positives; alphabetical within each group. Cards present in equal counts in both decks
// are omitted.
func runDiff(name1, name2 string) {
	p1 := resolveDeckPath(name1)
	p2 := resolveDeckPath(name2)
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
	if len(minuses) == 0 && len(pluses) == 0 {
		// Silent-on-match left the user guessing whether diff ran at all when the two decks had
		// matching card counts. Print an explicit confirmation so "no output" can't be confused
		// with "diff failed silently". Scoped to the card list to match what diff actually
		// compares — hero and weapon differences are not surfaced here.
		fmt.Printf("%s and %s have identical card lists (%d cards)\n", name1, name2, len(d1.Cards))
		return
	}
	for _, l := range minuses {
		fmt.Println(l)
	}
	for _, l := range pluses {
		fmt.Println(l)
	}
}
