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
	fmt.Printf("  Cycle 1 mean: %s\n", formatMean(s.FirstCycle.Mean()))
	fmt.Printf("  Cycle 2 mean: %s\n", formatMean(s.SecondCycle.Mean()))
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

// formatMean returns a 3-decimal rendering of a mean value.
func formatMean(mean float64) string {
	return fmt.Sprintf("%.3f", mean)
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

// printBestDeck dumps the full deck report: summary, card list, best turn played, the
// hand-value histogram, and the per-card value table pairing each card's role-based avg
// contribution with its correlational marginal hand-value lift. Sections silently skip
// themselves when the underlying Stats slice/map is empty so unscored decks still render
// the parts that do exist.
func printBestDeck(d *deck.Deck) {
	printDeckSummary(d)
	fmt.Println()
	printCardList(d)
	printBestTurn(d)
	if len(d.Stats.Histogram) > 0 {
		printHistogram(d, histogramTitle(d), naturalHistogramScale(d))
	}
	if len(d.Stats.PerCardMarginal) > 0 {
		printCardValues(d)
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

// printCardValues renders one row per unique card with two complementary signals:
//
//   - avg value: role-based contribution per appearance — attack power on attacks,
//     proportional prevented-damage share on defends, Pitch value on pitches. Captures
//     what the card typically does in the turn it's played.
//   - marginal +/-: mean turn value when the card sits in the dealt hand or arsenal-in
//     slot, minus the mean turn value when it's absent. Picks up within-turn indirect
//     lift (card draw, runechant generation, mid-turn triggers) the role-based avg misses.
//
// Sorted by marginal descending so suspected above-curve cards surface at the top and the
// drags sit at the bottom — the spread is a smell test for buggy implementations or
// oversimplified mechanics. See deck.Stats.PerCardMarginal for the cross-turn caveat that
// limits this view's reach for next-turn-payoff cards.
func printCardValues(d *deck.Deck) {
	type row struct {
		name      string
		avg       float64
		margin    float64
		hasAvg    bool
		hasMargin bool
	}
	rows := make([]row, 0, len(d.Stats.PerCardMarginal))
	for id, m := range d.Stats.PerCardMarginal {
		r := row{name: cards.Get(id).Name()}
		if play, ok := d.Stats.PerCard[id]; ok && (play.Plays+play.Pitches) > 0 {
			r.avg = play.Avg()
			r.hasAvg = true
		}
		if m.PresentHands > 0 && m.AbsentHands > 0 {
			r.margin = m.Marginal()
			r.hasMargin = true
		}
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].margin != rows[j].margin {
			return rows[i].margin > rows[j].margin
		}
		return rows[i].name < rows[j].name
	})

	fmt.Println()
	fmt.Println("Card value (marginal = mean turn value with vs without the card in hand or arsenal; avg = role-based contribution per appearance):")
	nameW := maxNameLen(d.Cards)
	for _, r := range rows {
		fmt.Printf("  %-*s  marginal %s  avg %s\n",
			nameW, r.name, formatCardMargin(r.margin, r.hasMargin), formatCardValue(r.avg, r.hasAvg))
	}
}

// formatCardValue renders a non-negative role-based avg in fixed width, falling back to a
// blank-aligned dash when the card never tallied a Play or Pitch (e.g. a card only ever
// held / arsenaled). The dash keeps the column aligned without printing a misleading 0.000.
func formatCardValue(v float64, has bool) string {
	if !has {
		return "    -"
	}
	return fmt.Sprintf("%5.2f", v)
}

// formatCardMargin renders the signed marginal value with an explicit sign. Cards present
// in every hand (or never present) have no comparison and render as a blank-aligned dash so
// the column stays read-able without an artificial 0.000.
func formatCardMargin(v float64, has bool) string {
	if !has {
		return "     -"
	}
	return fmt.Sprintf("%+6.2f", v)
}

// Dimensions of the hand-value histogram chart body. histWidth is the chart width in the
// compress regime; histHeight is the row count; histStretchSlot is the per-bar column budget
// (1 bar + slot-1 gap chars) in the stretch regime — every bar gets exactly slot cols of
// horizontal space so the spacing reads uniformly. histMaxStretchWidth caps how wide a
// stretched chart is allowed to grow before the chart compresses to histWidth and bins
// instead, so absurd hand-value spreads can't balloon the chart past terminal-friendly
// widths.
const (
	histWidth           = 60
	histHeight          = 12
	histStretchSlot     = 3
	histMaxStretchWidth = 120
)

// histogramScale fixes the axis ranges a chart renders against: x-axis spans minV..maxV
// inclusive and y-axis tops out at peak.
type histogramScale struct {
	minV, maxV, peak int
}

// naturalHistogramScale returns the histogramScale derived from d's own min/max/peak — what to
// pass to printHistogram for a single-deck chart at its native resolution.
func naturalHistogramScale(d *deck.Deck) histogramScale {
	minV := d.Stats.Min()
	maxV := d.Stats.Max()
	_, peak := buildHistogramColumns(d.Stats.Histogram, minV, maxV, histChartWidth(maxV-minV+1))
	return histogramScale{minV, maxV, peak}
}

// unionHistogramScale returns the smallest scale that fits both decks' data, so the two
// charts can be rendered with matching x and y axes. The peak is computed under the union
// range and width because binning shifts when the range widens.
func unionHistogramScale(d1, d2 *deck.Deck) histogramScale {
	minV := min(d1.Stats.Min(), d2.Stats.Min())
	maxV := max(d1.Stats.Max(), d2.Stats.Max())
	width := histChartWidth(maxV - minV + 1)
	_, peak1 := buildHistogramColumns(d1.Stats.Histogram, minV, maxV, width)
	_, peak2 := buildHistogramColumns(d2.Stats.Histogram, minV, maxV, width)
	return histogramScale{minV, maxV, max(peak1, peak2)}
}

// printHistogram renders Stats.Histogram as an ASCII bar chart under the supplied title line.
// The chart body is always histWidth x histHeight characters regardless of how many distinct
// hand values the deck produced — sparse data stretches across the width, dense data bins into
// it — so the rendered output has predictable size and the axis labels alone carry the scale.
// title is printed verbatim above the chart so the caller can identify the deck the chart is
// for. The scale fixes both axes; values outside scale.minV..scale.maxV contribute nothing,
// and bars scale against scale.peak rather than this deck's natural peak. No-ops when
// scale.peak == 0.
func printHistogram(d *deck.Deck, title string, scale histogramScale) {
	if scale.peak == 0 {
		return
	}
	width := histChartWidth(scale.maxV - scale.minV + 1)
	counts, _ := buildHistogramColumns(d.Stats.Histogram, scale.minV, scale.maxV, width)
	bars := scaleBarHeights(counts, scale.peak, histHeight)
	yLabelW := len(strconv.Itoa(scale.peak))
	// bodyIndent is the number of spaces before the first bar column so every row (y-label,
	// axis baseline, tick labels, title) lines up: 1 lead + yLabelW label + " |" (2 chars).
	bodyIndent := strings.Repeat(" ", 1+yLabelW+2)
	yTicks := yAxisTickLabels(scale.peak, histHeight)
	xTicks := xAxisTicks(scale.minV, scale.maxV, width)

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

// histChartWidth picks the chart body width for a given integer range. The chart stretches
// to (rng-1)*histStretchSlot+1 cols whenever that fits histMaxStretchWidth so every bar gets
// exactly histStretchSlot cols of horizontal space and the spacing stays uniform. Ranges
// that would exceed histMaxStretchWidth at slot=3 fall back to histWidth and the compress
// regime aggregates adjacent values into bins.
func histChartWidth(rng int) int {
	if rng <= 1 {
		return histWidth
	}
	if w := stretchWidth(rng); w <= histMaxStretchWidth {
		return w
	}
	return histWidth
}

// stretchWidth returns the chart width that fits rng integer values at slot=3: each bar
// occupies its leftmost col with histStretchSlot-1 spaces of gap, and the rightmost bar at
// (rng-1)*slot lives at the final col. The +1 closes the right edge under the last bar.
func stretchWidth(rng int) int {
	return (rng-1)*histStretchSlot + 1
}

// buildHistogramColumns bins the raw histogram map into width fixed-width columns spanning
// [minV, maxV] inclusive. Three regimes:
//   - rng == 1: the lone bar is centred, so a single-value chart isn't pinned to col 0.
//   - width fits stretchWidth(rng): each integer value gets its own one-character bar at
//     col (v-minV)*histStretchSlot, with histStretchSlot-1 empty cols of gap between bars.
//     Spacing is uniform regardless of rng.
//   - otherwise: compress. Each column aggregates the contiguous integer range that maps to
//     it; bars are necessarily adjacent since every column carries data.
//
// Returns the per-column count and the peak value so callers can scale bar heights.
func buildHistogramColumns(hist map[int]int, minV, maxV, width int) (counts []int, peak int) {
	counts = make([]int, width)
	rng := maxV - minV + 1
	if rng <= 0 {
		return counts, 0
	}
	switch {
	case rng == 1:
		counts[(width-1)/2] = hist[minV]
	case width >= stretchWidth(rng):
		for v := minV; v <= maxV; v++ {
			counts[(v-minV)*histStretchSlot] = hist[v]
		}
	default:
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
//   - Degenerate rng<=1 (single value): the lone bar is centred, so the column is width/2.
//   - Stretch regime (width fits stretchWidth(rng)): col = (v-minV)*histStretchSlot, giving
//     uniform slot-spaced positions independent of width.
//   - Compress regime (width can't fit stretch): bars are contiguous; v maps to the start of
//     its bin.
func colForValue(v, minV, maxV, width int) int {
	rng := maxV - minV + 1
	if rng <= 1 {
		return (width - 1) / 2
	}
	if width >= stretchWidth(rng) {
		return (v - minV) * histStretchSlot
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
