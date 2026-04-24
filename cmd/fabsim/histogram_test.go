package main

import "testing"

// TestBuildHistogramColumns_Stretch verifies the "range <= width" regime: when the histogram
// spans fewer distinct integer values than the chart is wide, each value stretches across
// multiple adjacent columns so the chart still fills the full width.
func TestBuildHistogramColumns_Stretch(t *testing.T) {
	// Three distinct values (min=0, max=2, range=3) across a 6-wide chart → every value gets
	// exactly two columns and all six cols pull real data.
	hist := map[int]int{0: 5, 1: 10, 2: 3}
	counts, peak := buildHistogramColumns(hist, 0, 2, 6)
	want := []int{5, 5, 10, 10, 3, 3}
	for i, got := range counts {
		if got != want[i] {
			t.Errorf("counts[%d] = %d, want %d (full=%v)", i, got, want[i], counts)
		}
	}
	if peak != 10 {
		t.Errorf("peak = %d, want 10", peak)
	}
}

// TestBuildHistogramColumns_Compress verifies the "range > width" regime: when the data spans
// more distinct integers than the chart is wide, adjacent values aggregate into the same
// column so the chart stays a fixed width without dropping data.
func TestBuildHistogramColumns_Compress(t *testing.T) {
	// Values 0..9 across a 5-wide chart — each column aggregates two adjacent integers.
	hist := map[int]int{0: 1, 1: 2, 2: 3, 3: 4, 4: 5, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10}
	counts, peak := buildHistogramColumns(hist, 0, 9, 5)
	want := []int{1 + 2, 3 + 4, 5 + 6, 7 + 8, 9 + 10}
	for i, got := range counts {
		if got != want[i] {
			t.Errorf("counts[%d] = %d, want %d", i, got, want[i])
		}
	}
	if peak != 19 {
		t.Errorf("peak = %d, want 19 (sum of 9+10)", peak)
	}
}

// TestBuildHistogramColumns_EmptyRange covers the pre-evaluation guard: a deck with no hands
// (min/max both zero, hist empty) gets a zero-filled column slice and a zero peak, so the
// caller can short-circuit cleanly instead of divide-by-zero on peak.
func TestBuildHistogramColumns_EmptyRange(t *testing.T) {
	counts, peak := buildHistogramColumns(map[int]int{}, 0, 0, 5)
	if peak != 0 {
		t.Errorf("peak = %d, want 0 on empty histogram", peak)
	}
	// Single-value range (min==max) still returns width-length slice; all entries are hist[min].
	for i, c := range counts {
		if c != 0 {
			t.Errorf("counts[%d] = %d, want 0", i, c)
		}
	}
}

// TestScaleBarHeights_ZeroStaysZero_NonZeroRoundsUp pins the visibility-floor contract: any
// column with a non-zero count renders at least one row tall, even if its count would
// proportionally round down to zero next to a much larger peak, so long-tail buckets don't
// silently disappear.
func TestScaleBarHeights_ZeroStaysZero_NonZeroRoundsUp(t *testing.T) {
	counts := []int{0, 1, 50, 100}
	bars := scaleBarHeights(counts, 100, 10)
	// Zero count → zero rows.
	if bars[0] != 0 {
		t.Errorf("bars[0] = %d, want 0 (zero count should stay hidden)", bars[0])
	}
	// 1/100 of a 10-row chart is 0 rows proportionally; the floor forces it to 1.
	if bars[1] != 1 {
		t.Errorf("bars[1] = %d, want 1 (tiny count should round up to visibility floor)", bars[1])
	}
	// 50/100 → 5 rows.
	if bars[2] != 5 {
		t.Errorf("bars[2] = %d, want 5 (50%% of peak)", bars[2])
	}
	// Peak itself → full height.
	if bars[3] != 10 {
		t.Errorf("bars[3] = %d, want 10 (peak fills the chart)", bars[3])
	}
}

// TestXAxisTicks_FullFiveOnWideRange verifies the happy path: a large-enough range across
// the full chart width produces five ticks (min, lower quartile, midpoint, upper quartile,
// max) in left-to-right order. min and max always anchor the ends.
func TestXAxisTicks_FullFiveOnWideRange(t *testing.T) {
	ticks := xAxisTicks(0, 59, 60)
	if len(ticks) != 5 {
		t.Fatalf("len = %d, want 5 (min + 3 quartiles + max)", len(ticks))
	}
	wantCols := []int{0, 14, 29, 44, 59}
	for i, tk := range ticks {
		if tk.col != wantCols[i] {
			t.Errorf("ticks[%d].col = %d, want %d", i, tk.col, wantCols[i])
		}
	}
	if ticks[0].value != 0 || ticks[len(ticks)-1].value != 59 {
		t.Errorf("ends = (%d, %d), want (0, 59)", ticks[0].value, ticks[len(ticks)-1].value)
	}
}

// TestXAxisTicks_DedupesNarrowRange pins the narrow-range contract: when the data spans only
// a handful of distinct integers, interior quartile ticks whose value duplicates a neighbour
// are dropped so the axis never prints the same label twice.
func TestXAxisTicks_DedupesNarrowRange(t *testing.T) {
	// min=7, max=9 across width 60: quartile cols 14/29/44 all map to values in {7, 8, 9}.
	// min(7) and max(9) are reserved; the only novel quartile value is 8 at col 29.
	ticks := xAxisTicks(7, 9, 60)
	values := map[int]int{}
	for _, tk := range ticks {
		values[tk.value]++
	}
	for v, n := range values {
		if n != 1 {
			t.Errorf("value %d appears %d times, want 1 (dedup should drop repeats)", v, n)
		}
	}
	if _, ok := values[7]; !ok {
		t.Error("min=7 missing from ticks")
	}
	if _, ok := values[9]; !ok {
		t.Error("max=9 missing from ticks")
	}
}

// TestXAxisTickRow_LayoutCentresLabels pins the label placement: each tick label is centred on
// its column (with edge ticks clipped inward to fit), and when labels would collide the
// interior tick is dropped so min and max always render.
func TestXAxisTickRow_LayoutCentresLabels(t *testing.T) {
	ticks := xAxisTicks(0, 59, 60)
	got := xAxisTickRow(ticks, 60)
	if len(got) != 60 {
		t.Fatalf("len = %d, want 60 (chart width)", len(got))
	}
	if got[0] != '0' {
		t.Errorf("first char = %q, want '0' (min left-clipped to col 0)", got[0])
	}
	// "59" is right-clipped so the 9 lands at col 59 and the 5 at col 58.
	if got[58:60] != "59" {
		t.Errorf("last two chars = %q, want \"59\" (max right-clipped to last col)", got[58:60])
	}
}

// TestXAxisBaseline_MarksTickPositions covers the bottom-axis rendering: the leading "+"
// sits under the y-axis, the body is dashes, and each interior tick gets an additional "+"
// anchor aligned under its label.
func TestXAxisBaseline_MarksTickPositions(t *testing.T) {
	ticks := xAxisTicks(0, 59, 60)
	base := xAxisBaseline(ticks, 60)
	if len(base) != 61 {
		t.Fatalf("len = %d, want 61 (leading + plus 60 cols)", len(base))
	}
	if base[0] != '+' {
		t.Errorf("base[0] = %q, want '+'", base[0])
	}
	// Interior ticks at cols 14, 29, 44 produce "+" one position to the right (buf[col+1]).
	for _, col := range []int{14, 29, 44} {
		if base[col+1] != '+' {
			t.Errorf("base[%d] = %q, want '+' (tick anchor)", col+1, base[col+1])
		}
	}
	// Rightmost tick at col 59 lands at buf[60].
	if base[60] != '+' {
		t.Errorf("base[60] = %q, want '+' (max-tick anchor)", base[60])
	}
}

// TestYAxisTickLabels_FourRowsOnTallPeak verifies that a peak comfortably above the height
// emits the expected four-row label set: row 0 at the peak plus three interior quartile
// rows at 3/4, 1/2, 1/4 of the peak.
func TestYAxisTickLabels_FourRowsOnTallPeak(t *testing.T) {
	ticks := yAxisTickLabels(1200, 12)
	want := map[int]int{0: 1200, 3: 900, 6: 600, 9: 300}
	if len(ticks) != len(want) {
		t.Fatalf("len = %d, want %d (peak + 3 quartiles)", len(ticks), len(want))
	}
	for row, v := range want {
		if ticks[row] != v {
			t.Errorf("ticks[%d] = %d, want %d", row, ticks[row], v)
		}
	}
}

// TestYAxisTickLabels_CollapsesTinyPeak pins the dedup contract: when the peak is small
// enough that multiple quartile rows would report the same integer value, only the first
// occurrence is kept so the axis never prints duplicates.
func TestYAxisTickLabels_CollapsesTinyPeak(t *testing.T) {
	// Peak=1 over height=12: all interior quartiles compute to 0, so only row 0 survives.
	ticks := yAxisTickLabels(1, 12)
	if len(ticks) != 1 {
		t.Errorf("len = %d, want 1 (tiny peak should collapse interior ticks)", len(ticks))
	}
	if ticks[0] != 1 {
		t.Errorf("ticks[0] = %d, want 1", ticks[0])
	}
}

// TestCenterLabel centres a label within the chart width so the title reads under the centre
// of the bar area. Labels longer than the width pass through unchanged rather than overflow
// left.
func TestCenterLabel(t *testing.T) {
	if got := centerLabel("hi", 6); got != "  hi" {
		t.Errorf("centerLabel(\"hi\", 6) = %q, want %q", got, "  hi")
	}
	// Odd leftover space biases left (integer division).
	if got := centerLabel("hi", 5); got != " hi" {
		t.Errorf("centerLabel(\"hi\", 5) = %q, want %q", got, " hi")
	}
	// Label wider than the chart passes through verbatim.
	if got := centerLabel("overflowing title", 5); got != "overflowing title" {
		t.Errorf("long label should pass through unchanged; got %q", got)
	}
}
