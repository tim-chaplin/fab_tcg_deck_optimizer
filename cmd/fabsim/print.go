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
	red, yellow, blue := pitchCounts(d.Cards)
	fmt.Printf("Hero:    %s\n", d.Hero.Name())
	fmt.Printf("Weapons: %s\n", weaponNames(d.Weapons))
	fmt.Printf("Pitch:   %d red / %d yellow / %d blue\n", red, yellow, blue)
	fmt.Println()
	fmt.Printf("Mean value: %.3f  (%s shuffles)\n", s.Mean(), commaInt(s.Runs))
	fmt.Printf("  Cycle 1 mean: %.3f\n", s.FirstCycle.Mean())
	fmt.Printf("  Cycle 2 mean: %.3f\n", s.SecondCycle.Mean())
}

// pitchCounts tallies red/yellow/blue copies by Pitch() value so the summary's "Pitch:" line
// stays a single expression. Cards with pitch outside 1-3 contribute to no bucket.
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
		printHistogram(d)
	}
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

// Fixed dimensions of the hand-value histogram chart body. Chosen to fit comfortably in an
// 80-column terminal alongside the y-axis label column and a little left margin.
const (
	histWidth  = 60
	histHeight = 12
)

// printHistogram renders Stats.Histogram as an ASCII bar chart. The chart body is always
// histWidth × histHeight characters regardless of how many distinct hand values the deck
// produced — sparse data stretches across the width, dense data bins into it — so the
// rendered output has predictable size and the axis labels alone carry the scale. No-ops on
// an unscored deck. Called by printBestDeck after the per-card stats block.
func printHistogram(d *deck.Deck) {
	minV := d.Stats.Min()
	maxV := d.Stats.Max()
	counts, peak := buildHistogramColumns(d.Stats.Histogram, minV, maxV, histWidth)
	if peak == 0 {
		return
	}
	bars := scaleBarHeights(counts, peak, histHeight)
	yLabelW := len(strconv.Itoa(peak))
	// bodyIndent is the number of spaces before the first bar column so every row (y-label,
	// axis baseline, tick labels, title) lines up: 1 lead + yLabelW label + " |" (2 chars).
	bodyIndent := strings.Repeat(" ", 1+yLabelW+2)

	fmt.Println()
	fmt.Printf("Hand-value distribution (%s hands):\n\n", commaInt(d.Stats.Hands))
	for row := 0; row < histHeight; row++ {
		// Top-down: row 0 is the peak, row histHeight-1 is one slot above the axis. A bar of
		// height h fills the bottom h rows, so this row is filled when histHeight-row <= h.
		rowFromBottom := histHeight - row
		var label string
		if row == 0 {
			label = fmt.Sprintf(" %*d", yLabelW, peak)
		} else {
			label = strings.Repeat(" ", yLabelW+1)
		}
		fmt.Print(label + " |")
		for col := 0; col < histWidth; col++ {
			if bars[col] >= rowFromBottom {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	// X-axis baseline and endpoint ticks. The baseline's "+" sits directly under the y-axis
	// bars, so the tick-label and title rows both live under the chart body (bodyIndent).
	fmt.Printf(" %*d +%s\n", yLabelW, 0, strings.Repeat("-", histWidth))
	fmt.Println(bodyIndent + xAxisTickRow(minV, maxV, histWidth))
	fmt.Println(bodyIndent + centerLabel("hand value", histWidth))
}

// buildHistogramColumns bins the raw histogram map into width fixed-width columns spanning
// [minV, maxV] inclusive. Two regimes:
//   - range ≤ width: each integer value gets one or more columns (the chart stretches).
//   - range > width: each column covers a contiguous integer range (the chart compresses).
//
// Returns the per-column count and the peak value so callers can scale bar heights.
func buildHistogramColumns(hist map[int]int, minV, maxV, width int) (counts []int, peak int) {
	counts = make([]int, width)
	rng := maxV - minV + 1
	if rng <= 0 {
		return counts, 0
	}
	if rng <= width {
		// Stretch: column col maps to a single integer value; adjacent columns may share one.
		for col := 0; col < width; col++ {
			v := minV + col*rng/width
			counts[col] = hist[v]
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

// xAxisTickRow returns a width-character string with minV left-justified and maxV
// right-justified. Collapses to "minV..maxV" when the two labels can't both fit.
func xAxisTickRow(minV, maxV, width int) string {
	lo := strconv.Itoa(minV)
	hi := strconv.Itoa(maxV)
	if len(lo)+len(hi)+1 > width {
		return lo + ".." + hi
	}
	return lo + strings.Repeat(" ", width-len(lo)-len(hi)) + hi
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
