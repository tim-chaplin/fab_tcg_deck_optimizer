package main

import (
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckformat"
)

// runCompareCmd parses compare's flags and dispatches to runCompare. Both decks are positional
// args; simulation knobs (-deep-shuffles, -incoming, -seed, -format, -max-copies) are accepted
// so both decks can be re-scored under matched conditions before comparison. -incoming is
// required because comparing decks scored against different incoming-damage assumptions would
// be apples to oranges.
func runCompareCmd(args []string) {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: fabsim compare <deck1> <deck2> [flags]")
		fmt.Fprintln(fs.Output())
		fmt.Fprintln(fs.Output(), "Flags:")
		fs.PrintDefaults()
	}
	deepShuffles := fs.Int("deep-shuffles", 10000, "shuffles per deck used to re-score each deck before comparing")
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required — both decks are re-scored against this value)")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed (each deck builds an independent RNG from this seed)")
	formatFlag := fs.String("format", string(deckformat.SilverAge), "constructed format predicate applied to replacement picks when a loaded deck contains NotImplemented cards")
	maxCopies := fs.Int("max-copies", defaultMaxCopies, "maximum copies of any single card printing per deck, applied when replacing NotImplemented cards in a loaded deck")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 2 {
		die("compare: need exactly 2 positional deck names (got %d); try `fabsim compare <deck1> <deck2>`", fs.NArg())
	}
	requireFlag(fs, "compare", "incoming")
	fmtValue, err := deckformat.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}
	runCompare(fs.Arg(0), fs.Arg(1), *deepShuffles, *incoming, *maxCopies, *seed, fmtValue)
}

// runCompare re-evaluates both decks under identical (deepShuffles, incoming) settings so the
// printed stats are apples to apples regardless of what conditions either deck happened to be
// scored under last, then prints a stat-by-stat side-by-side comparison: pitch counts, mean
// hand value, per-cycle means, the hand-value histograms, and finally the per-card count
// delta. The header line at the top of the output records the (deepShuffles, incoming)
// settings so the per-section rows don't have to repeat them.
func runCompare(name1, name2 string, deepShuffles, incoming, maxCopies int, seed int64, fmtValue deckformat.Format) {
	d1 := evaluateAndPersist(resolveDeckPath(name1), deepShuffles, incoming, maxCopies, seed, fmtValue)
	d2 := evaluateAndPersist(resolveDeckPath(name2), deepShuffles, incoming, maxCopies, seed, fmtValue)
	s1, s2 := d1.Stats, d2.Stats

	fmt.Printf("compare: -deep-shuffles=%s -incoming=%d\n", commaInt(deepShuffles), incoming)
	fmt.Println()

	printSideBySideStats(name1, name2, []statSection{
		{"Pitch values", pitchCountsLine(d1.Cards), pitchCountsLine(d2.Cards)},
		{"Mean hand value", formatMean(s1.Mean()), formatMean(s2.Mean())},
		{"Cycle 1 mean", formatMean(s1.FirstCycle.Mean()), formatMean(s2.FirstCycle.Mean())},
		{"Cycle 2 mean", formatMean(s1.SecondCycle.Mean()), formatMean(s2.SecondCycle.Mean())},
	})

	if len(s1.Histogram) > 0 && len(s2.Histogram) > 0 {
		// Both decks render against the same axis ranges so values, bar widths, and tick
		// labels line up between the two charts and a side-by-side read is meaningful.
		scale := unionHistogramScale(d1, d2)
		fmt.Println()
		fmt.Println("Hand-value distributions:")
		printHistogram(d1, fmt.Sprintf("  %s:", name1), scale)
		printHistogram(d2, fmt.Sprintf("  %s:", name2), scale)
	}

	fmt.Println()
	printCardDelta(name1, name2, d1, d2)
}

// printCardDelta writes the per-loadout count delta between d1 and d2 — negative rows first,
// then positives; weapons lead each block (alphabetical) followed by cards (alphabetical),
// so the loadout-defining piece sits at the top of each list. Entries present in equal
// counts in both decks are omitted. Covers both the deck's Cards and its Weapons since both
// are sim-relevant loadout choices a comparison reader needs to see. Hero, Equipment, and
// Sideboard are out of scope. When the two loadouts match exactly an explicit confirmation
// line replaces the empty body so silence can't be mistaken for a failure.
func printCardDelta(name1, name2 string, d1, d2 *deck.Deck) {
	counts1 := loadoutCounts(d1)
	counts2 := loadoutCounts(d2)
	weaponNames := loadoutWeaponNames(d1, d2)

	type entry struct {
		line     string
		name     string
		isWeapon bool
	}
	allNames := make(map[string]struct{}, len(counts1)+len(counts2))
	for n := range counts1 {
		allNames[n] = struct{}{}
	}
	for n := range counts2 {
		allNames[n] = struct{}{}
	}

	var minuses, pluses []entry
	for n := range allNames {
		delta := counts2[n] - counts1[n]
		if delta == 0 {
			continue
		}
		_, isWeapon := weaponNames[n]
		e := entry{name: n, isWeapon: isWeapon}
		if delta < 0 {
			e.line = fmt.Sprintf("%d %s", delta, n)
			minuses = append(minuses, e)
		} else {
			e.line = fmt.Sprintf("+%d %s", delta, n)
			pluses = append(pluses, e)
		}
	}
	sortLoadoutEntries := func(es []entry) {
		sort.Slice(es, func(i, j int) bool {
			if es[i].isWeapon != es[j].isWeapon {
				return es[i].isWeapon
			}
			return es[i].name < es[j].name
		})
	}
	sortLoadoutEntries(minuses)
	sortLoadoutEntries(pluses)

	fmt.Println("Card / weapon differences:")
	if len(minuses) == 0 && len(pluses) == 0 {
		fmt.Printf("  %s and %s have identical card and weapon lists (%d cards, %d weapons)\n",
			name1, name2, len(d1.Cards), len(d1.Weapons))
		return
	}
	for _, e := range minuses {
		fmt.Println("  " + e.line)
	}
	for _, e := range pluses {
		fmt.Println("  " + e.line)
	}
}

// loadoutCounts tallies the deck's cards and weapons by display name in a single map.
// Weapon names don't collide with card names in the current registry, so a flat
// name-keyed map cleanly captures both lists for diffing. card.DisplayName keeps pitch
// printings as distinct entries so a "-1 Aether Slash [R], +1 Aether Slash [Y]" diff is
// legible.
func loadoutCounts(d *deck.Deck) map[string]int {
	out := make(map[string]int, len(d.Cards)+len(d.Weapons))
	for _, c := range d.Cards {
		out[card.DisplayName(c)]++
	}
	for _, w := range d.Weapons {
		out[card.DisplayName(w)]++
	}
	return out
}

// loadoutWeaponNames returns the set of weapon display names appearing in either deck,
// used to flag which diff entries should sort first within their +/- block.
func loadoutWeaponNames(decks ...*deck.Deck) map[string]struct{} {
	out := map[string]struct{}{}
	for _, d := range decks {
		for _, w := range d.Weapons {
			out[card.DisplayName(w)] = struct{}{}
		}
	}
	return out
}
