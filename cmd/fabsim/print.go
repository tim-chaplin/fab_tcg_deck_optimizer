package main

// Human-readable renderings shared by the subcommands: the compact deck summary, the full
// "best deck" report (summary + card list + best turn + per-card stats), and the small
// formatting helpers every subcommand reaches for (weapon-list join, comma-grouped integers,
// max-name-length column sizing).

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// printCardList writes the deck's card list in canonical "Card list:" form: one
// grouped-and-sorted count-and-name line per unique card. When the deck carries user-managed
// Equipment or Sideboard sections, trailing "Equipment:" / "Sideboard:" blocks list their
// contents in the same grouped form — empty sections are silently skipped so stock decks
// stay untouched.
func printCardList(d *deck.Deck) {
	fmt.Println("Card list:")
	printGroupedCards(d.Cards)
	if len(d.Equipment) > 0 {
		fmt.Println("Equipment:")
		printGroupedStrings(d.Equipment)
	}
	if len(d.Sideboard) > 0 {
		fmt.Println("Sideboard:")
		printGroupedStrings(d.Sideboard)
	}
}

// printGroupedCards writes one count-and-name line per unique card in cs, sorted by name.
// Shared between the main card list and the sideboard block so formatting stays consistent.
func printGroupedCards(cs []card.Card) {
	names := make([]string, len(cs))
	for i, c := range cs {
		names[i] = c.Name()
	}
	printGroupedStrings(names)
}

// printGroupedStrings is the string-slice counterpart of printGroupedCards — used by the
// Equipment section where entries are opaque names rather than registry cards.
func printGroupedStrings(ss []string) {
	counts := map[string]int{}
	for _, s := range ss {
		counts[s]++
	}
	names := make([]string, 0, len(counts))
	for n := range counts {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("  %dx %s\n", counts[n], n)
	}
}

// printDeckSummary prints the compact summary: a loadout block (hero, weapons, pitch colour
// counts) followed by a blank line and a stats block (overall mean, per-cycle means).
// printBestDeck wraps this with the card list, best-turn block, and per-card stats;
// `fabsim eval -brief` calls printDeckSummary directly so a scripted re-score gets just the
// numbers without the card-list scroll.
func printDeckSummary(d *deck.Deck) {
	s := d.Stats
	fmt.Printf("Hero:    %s\n", d.Hero.Name())
	fmt.Printf("Weapons: %s\n", weaponNames(d.Weapons))
	fmt.Printf("Pitch:   %s\n", pitchCountsLine(d.Cards))
	fmt.Println()
	fmt.Printf("Mean value: %s\n", meanValueLine(s))
	fmt.Printf("  Cycle 1 mean: %s\n", cycleMeanLine(s.FirstCycle))
	fmt.Printf("  Cycle 2 mean: %s\n", cycleMeanLine(s.SecondCycle))
}

// pitchCountsLine returns the "20 red / 8 yellow / 12 blue" rendering of the deck's pitch
// distribution.
func pitchCountsLine(cs []card.Card) string {
	red, yellow, blue := pitchCounts(cs)
	return fmt.Sprintf("%d red / %d yellow / %d blue", red, yellow, blue)
}

// meanValueLine returns "14.041 (10,000 shuffles)" — the deck's overall mean plus the run
// count that produced it.
func meanValueLine(s deck.Stats) string {
	return fmt.Sprintf("%.3f (%s shuffles)", s.Mean(), commaInt(s.Runs))
}

// cycleMeanLine returns the per-cycle mean as a 3-decimal string.
func cycleMeanLine(c deck.CycleStats) string {
	return fmt.Sprintf("%.3f", c.Mean())
}

// statSection is one row of compare's side-by-side stat layout: a label header and the two
// rendered values to slot under it (val1 sits next to name1, val2 next to name2).
type statSection struct {
	label, val1, val2 string
}

// printSideBySideStats stacks several labelled stat blocks, one per section:
//
//	<label>:
//	  <name1>: <val1>
//	  <name2>: <val2>
//
// Deck names are right-padded to the longer of the two so every section's value column lines
// up vertically. Sections are separated by a blank line.
func printSideBySideStats(name1, name2 string, sections []statSection) {
	nameW := len(name1)
	if len(name2) > nameW {
		nameW = len(name2)
	}
	for i, sec := range sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s:\n", sec.label)
		fmt.Printf("  %-*s  %s\n", nameW+1, name1+":", sec.val1)
		fmt.Printf("  %-*s  %s\n", nameW+1, name2+":", sec.val2)
	}
}

// pitchCounts tallies red/yellow/blue copies by Pitch() value. Cards with pitch outside 1-3
// contribute to no bucket.
func pitchCounts(cs []card.Card) (red, yellow, blue int) {
	for _, c := range cs {
		switch c.Pitch() {
		case 1:
			red++
		case 2:
			yellow++
		case 3:
			blue++
		}
	}
	return red, yellow, blue
}

func printBestDeck(d *deck.Deck) {
	printDeckSummary(d)
	fmt.Println()
	printCardList(d)
	printBestTurn(d)
	if len(d.Stats.PerCard) > 0 {
		printPerCardStats(d)
	}
	if len(d.Stats.Histogram) > 0 {
		printHistogram(d, histogramTitle(d))
	}
}

// histogramTitle returns the standard "Hand-value distribution (N hands):" header.
func histogramTitle(d *deck.Deck) string {
	return fmt.Sprintf("Hand-value distribution (%s hands):", commaInt(d.Stats.Hands))
}

// printBestTurn renders the persisted peak-Value turn — "Best turn played (value N):"
// header plus FormatBestTurn's sectioned play order — when Stats.Best holds one. No-ops on
// an unscored deck so callers don't have to guard. Shared by printBestDeck and runEval.
// Starting Runechants (Runeblade carryover state) pipe through to FormatBestTurn, which
// folds a non-zero count into its "Auras in play at start of turn:" line; zero is omitted
// entirely so only turns that actually started with pending tokens surface the detail.
func printBestTurn(d *deck.Deck) {
	b := d.Stats.Best
	if len(b.Summary.BestLine) == 0 {
		return
	}
	fmt.Println()
	fmt.Printf("Best turn played (value %d):\n", b.Summary.Value)
	fmt.Println(hand.FormatBestTurn(b.Summary, b.StartingRunechants))
}

// printPerCardStats renders per-card averages collected by deck.Evaluate: mean per-card
// contribution across hands the card appeared in. Contribution is role-based (attack power on
// attacks, proportional prevented-damage share on defends, Pitch on pitches), so the ranking
// reflects what each card typically does in its hand rather than the hand's total value.
func printPerCardStats(d *deck.Deck) {
	type row struct {
		name           string
		deckCount      int
		plays, pitches int
		avg            float64
	}
	deckCounts := map[card.ID]int{}
	for _, c := range d.Cards {
		deckCounts[c.ID()]++
	}
	rows := make([]row, 0, len(d.Stats.PerCard))
	for id, s := range d.Stats.PerCard {
		rows = append(rows, row{
			name:      cards.Get(id).Name(),
			deckCount: deckCounts[id],
			plays:     s.Plays,
			pitches:   s.Pitches,
			avg:       s.Avg(),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].avg != rows[j].avg {
			return rows[i].avg > rows[j].avg
		}
		ni, nj := rows[i].plays+rows[i].pitches, rows[j].plays+rows[j].pitches
		if ni != nj {
			return ni > nj
		}
		return rows[i].name < rows[j].name
	})

	fmt.Println()
	fmt.Println("Card value (avg contribution per appearance: attack=power, defend=share of block, pitch=resource):")
	for _, r := range rows {
		fmt.Printf("  %-*s avg %6.3f over %4d hands (%4d plays, %4d pitches, %dx in deck)\n",
			maxNameLen(d.Cards), r.name, r.avg, r.plays+r.pitches, r.plays, r.pitches, r.deckCount)
	}
}

// Dimensions of the hand-value histogram chart body. histWidth is the hard cap so the chart
// always fits under an 80-column terminal; histHeight is the row count; histStretchSlot is
// the target per-bar column budget (1 bar + 2 spaces) in the stretch regime — the chart
// shrinks below histWidth when the range fits at this spacing, and falls back to evenly
// distributing bars across histWidth when the range is too wide to honour the slot.
const (
	histWidth        = 60
	histHeight       = 12
	histStretchSlot  = 3
)

// printHistogram renders Stats.Histogram as an ASCII bar chart under the supplied title line.
// The chart body is always histWidth x histHeight characters regardless of how many distinct
// hand values the deck produced — sparse data stretches across the width, dense data bins into
// it — so the rendered output has predictable size and the axis labels alone carry the scale.
// title is printed verbatim above the chart so the caller can identify the deck the chart is
// for. No-ops on an unscored deck.
func printHistogram(d *deck.Deck, title string) {
	minV := d.Stats.Min()
	maxV := d.Stats.Max()
	width := histChartWidth(maxV - minV + 1)
	counts, peak := buildHistogramColumns(d.Stats.Histogram, minV, maxV, width)
	if peak == 0 {
		return
	}
	bars := scaleBarHeights(counts, peak, histHeight)
	yLabelW := len(strconv.Itoa(peak))
	// bodyIndent is the number of spaces before the first bar column so every row (y-label,
	// axis baseline, tick labels, title) lines up: 1 lead + yLabelW label + " |" (2 chars).
	bodyIndent := strings.Repeat(" ", 1+yLabelW+2)
	yTicks := yAxisTickLabels(peak, histHeight)
	xTicks := xAxisTicks(minV, maxV, width)

	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	for row := 0; row < histHeight; row++ {
		// Top-down: row 0 is the peak, row histHeight-1 is one slot above the axis. A bar of
		// height h fills the bottom h rows, so this row is filled when histHeight-row <= h.
		rowFromBottom := histHeight - row
		var label string
		if v, ok := yTicks[row]; ok {
			label = fmt.Sprintf(" %*d", yLabelW, v)
		} else {
			label = strings.Repeat(" ", yLabelW+1)
		}
		fmt.Print(label + " |")
		for col := 0; col < width; col++ {
			if bars[col] >= rowFromBottom {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	// Baseline: the "+" under the y-axis is the origin tick; additional "+" marks sit under
	// each interior x-axis tick so the scale is readable at a glance. The tick-label and
	// title rows both live under the chart body (bodyIndent).
	fmt.Printf(" %*d %s\n", yLabelW, 0, xAxisBaseline(xTicks, width))
	fmt.Println(bodyIndent + xAxisTickRow(xTicks, width))
	fmt.Println(bodyIndent + centerLabel("hand value", width))
}

// histChartWidth picks the chart body width for a given integer range. Short ranges shrink
// the chart to rng*histStretchSlot-1 cols so each bar gets a fixed (slot-1)-space gap
// instead of the airy spread you'd get if all bars were forced to span histWidth. Longer
// ranges that would exceed histWidth at this spacing clamp to histWidth; the compress
// regime uses the full histWidth for maximum resolution.
func histChartWidth(rng int) int {
	if rng <= 1 {
		return histWidth
	}
	if rng > histWidth {
		return histWidth
	}
	desired := rng*histStretchSlot - (histStretchSlot - 1)
	if desired > histWidth {
		return histWidth
	}
	return desired
}

// buildHistogramColumns bins the raw histogram map into width fixed-width columns spanning
// [minV, maxV] inclusive. Two regimes:
//   - range <= width: each integer value gets its own one-character bar placed at an evenly
//     distributed column; remaining columns stay zero, giving visual separation between bars
//     so adjacent values don't fuse into a solid block.
//   - range > width: each column aggregates a contiguous integer range (the chart compresses
//     and bars are necessarily contiguous since every column carries data).
//
// Returns the per-column count and the peak value so callers can scale bar heights.
func buildHistogramColumns(hist map[int]int, minV, maxV, width int) (counts []int, peak int) {
	counts = make([]int, width)
	rng := maxV - minV + 1
	if rng <= 0 {
		return counts, 0
	}
	if rng <= width {
		// Stretch: one bar per integer value, evenly distributed so the leftmost value lands
		// at col 0 and the rightmost at col width-1. A single-value range (rng=1) can't span
		// anything, so centre it.
		if rng == 1 {
			counts[(width-1)/2] = hist[minV]
		} else {
			for v := minV; v <= maxV; v++ {
				col := (v - minV) * (width - 1) / (rng - 1)
				counts[col] = hist[v]
			}
		}
	} else {
		// Compress: column col aggregates every value in [binLo, binHi).
		for col := 0; col < width; col++ {
			binLo := minV + col*rng/width
			binHi := minV + (col+1)*rng/width
			n := 0
			for v := binLo; v < binHi; v++ {
				n += hist[v]
			}
			counts[col] = n
		}
	}
	for _, c := range counts {
		if c > peak {
			peak = c
		}
	}
	return counts, peak
}

// scaleBarHeights converts per-column counts into bar heights in rows. Non-zero counts always
// round up to at least one row so tiny buckets stay visible next to a tall peak.
func scaleBarHeights(counts []int, peak, height int) []int {
	bars := make([]int, len(counts))
	for i, c := range counts {
		if c == 0 {
			continue
		}
		h := c * height / peak
		if h == 0 {
			h = 1
		}
		bars[i] = h
	}
	return bars
}

// yAxisTickLabels picks up to four rows on the y-axis to carry numeric labels: the peak (row
// 0) plus up-to-three interior quartile rows. Quartile labels are skipped when their
// integer-truncated value would duplicate a neighbour — e.g. a peak of 1 collapses all
// interior ticks to 0, so only the peak label survives. The map is keyed by row index so the
// render loop can look up per-row without re-computing.
func yAxisTickLabels(peak, height int) map[int]int {
	// Row 0 is the peak; rows height/4, height/2, 3*height/4 mark the interior quartiles. The
	// value at row r represents the top of that row, so v = peak * (height - r) / height.
	ticks := map[int]int{0: peak}
	candidates := []int{height / 4, height / 2, 3 * height / 4}
	seen := map[int]bool{peak: true, 0: true}
	for _, r := range candidates {
		if r == 0 {
			continue
		}
		v := peak * (height - r) / height
		if seen[v] {
			continue
		}
		seen[v] = true
		ticks[r] = v
	}
	return ticks
}

// xAxisTick marks a single x-axis position and the hand value it represents.
type xAxisTick struct {
	col   int
	value int
}

// xAxisTicks returns up to five x-axis ticks at the quartile VALUES of [minV, maxV] (min,
// lower quartile, midpoint, upper quartile, max), each mapped to its bar column via
// colForValue so every tick lands directly under the bar it labels. Duplicate values are
// deduped so a narrow range (e.g. a 3-integer spread) can't emit the same label twice.
// rng = maxV-minV+1 treats the span as an inclusive count of integer values, matching how
// buildHistogramColumns derives its per-column value mapping.
func xAxisTicks(minV, maxV, width int) []xAxisTick {
	if width <= 0 {
		return nil
	}
	rng := maxV - minV + 1
	if rng <= 0 {
		return nil
	}
	// Target values: min and max anchor the ends; three interior quartile values fill the
	// middle. Interior ticks are added only when their value is novel so narrow ranges
	// collapse to just min+max rather than repeating labels.
	placed := map[int]bool{}
	ticks := make([]xAxisTick, 0, 5)
	add := func(v int) {
		if placed[v] {
			return
		}
		placed[v] = true
		ticks = append(ticks, xAxisTick{col: colForValue(v, minV, maxV, width), value: v})
	}
	add(minV)
	add(maxV)
	add(minV + (rng-1)/4)
	add(minV + (rng-1)/2)
	add(minV + 3*(rng-1)/4)
	sort.Slice(ticks, func(i, j int) bool { return ticks[i].col < ticks[j].col })
	return ticks
}

// colForValue returns the chart column that renders the bar for integer value v. The formula
// mirrors buildHistogramColumns so tick labels always land directly under their bar:
//   - Stretch regime (rng <= width): bars are at evenly-distributed positions with minV at
//     col 0 and maxV at col width-1.
//   - Compress regime (rng > width): bars are contiguous; v maps to the start of its bin.
//   - Degenerate rng<=1 (single value): the lone bar is centred, so the column is width/2.
func colForValue(v, minV, maxV, width int) int {
	rng := maxV - minV + 1
	if rng <= 1 {
		return (width - 1) / 2
	}
	if rng <= width {
		return (v - minV) * (width - 1) / (rng - 1)
	}
	return (v - minV) * width / rng
}

// xAxisBaseline renders the chart's bottom axis: a leading "+" under the y-axis, then
// "-" across the full chart width with an additional "+" at every interior tick position so
// the tick labels below have visible anchors on the axis.
func xAxisBaseline(ticks []xAxisTick, width int) string {
	buf := make([]byte, width+1) // +1 for the leading "+"
	buf[0] = '+'
	for i := 1; i < len(buf); i++ {
		buf[i] = '-'
	}
	for _, t := range ticks {
		// Skip the leftmost tick; the leading "+" already marks it. Any tick at the last
		// column lands at buf[width], which is still inside the buffer.
		if t.col == 0 {
			continue
		}
		buf[t.col+1] = '+'
	}
	return string(buf)
}

// xAxisTickRow renders the tick-label row below the baseline. Each tick label is centred on
// its column; labels that would overflow the chart are clipped inward so no character lands
// outside [0, width). When two labels would collide, the later one is dropped — min and max
// are added first so interior labels are the ones that lose.
func xAxisTickRow(ticks []xAxisTick, width int) string {
	buf := make([]byte, width)
	for i := range buf {
		buf[i] = ' '
	}
	place := func(col int, label string) bool {
		start := col - len(label)/2
		if start < 0 {
			start = 0
		}
		if start+len(label) > width {
			start = width - len(label)
		}
		for i := 0; i < len(label); i++ {
			if buf[start+i] != ' ' {
				return false
			}
		}
		copy(buf[start:], label)
		return true
	}
	// min and max first (by extracting them from the sorted ticks list), then interior.
	var mn, mx *xAxisTick
	var interior []xAxisTick
	for i := range ticks {
		switch {
		case ticks[i].col == 0:
			mn = &ticks[i]
		case ticks[i].col == width-1:
			mx = &ticks[i]
		default:
			interior = append(interior, ticks[i])
		}
	}
	if mn != nil {
		place(mn.col, strconv.Itoa(mn.value))
	}
	if mx != nil {
		place(mx.col, strconv.Itoa(mx.value))
	}
	for _, t := range interior {
		place(t.col, strconv.Itoa(t.value))
	}
	return string(buf)
}

// centerLabel returns s padded with spaces so it is horizontally centered within width.
// Shorter-than-width strings are centred; equal-or-longer strings pass through unchanged.
func centerLabel(s string, width int) string {
	if len(s) >= width {
		return s
	}
	left := (width - len(s)) / 2
	return strings.Repeat(" ", left) + s
}

// maxNameLen returns the length of the longest Name() across cs, or 0 when empty. Used to
// size fixed-width card-name columns in printed tables.
func maxNameLen(cs []card.Card) int {
	m := 0
	for _, c := range cs {
		if n := len(c.Name()); n > m {
			m = n
		}
	}
	return m
}

// weaponNames joins the deck's weapon names with ", " for the summary's "Weapons:" line.
// A single-weapon loadout prints the name bare; an empty loadout prints "none" so the column
// stays filled rather than rendering as a trailing blank.
func weaponNames(ws []weapon.Weapon) string {
	if len(ws) == 0 {
		return "none"
	}
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	return strings.Join(names, ", ")
}

// commaInt renders n with ',' thousands separators (e.g. 10000 -> "10,000"). Used for the
// summary's shuffle count so six- and seven-digit totals stay legible at a glance.
func commaInt(n int) string {
	s := strconv.Itoa(n)
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	if len(s) <= 3 {
		return sign + s
	}
	var b strings.Builder
	b.Grow(len(s) + (len(s)-1)/3)
	head := len(s) % 3
	if head == 0 {
		head = 3
	}
	b.WriteString(s[:head])
	for i := head; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	return sign + b.String()
}
