package main

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// Tests the stretch regime: each integer value places one bar at col (v-min)*slot, with
// slot-1 empty cols between so adjacent values don't fuse.
func TestBuildHistogramColumns_Stretch(t *testing.T) {
	// Three distinct values (min=0, max=2, range=3) at slot=3 → bars at cols 0, 3, 6 within
	// the natural 7-wide stretch chart.
	hist := map[int]int{0: 5, 1: 10, 2: 3}
	counts, peak := buildHistogramColumns(hist, 0, 2, 7)
	want := []int{5, 0, 0, 10, 0, 0, 3}
	for i, got := range counts {
		if got != want[i] {
			t.Errorf("counts[%d] = %d, want %d (full=%v)", i, got, want[i], counts)
		}
	}
	if peak != 10 {
		t.Errorf("peak = %d, want 10", peak)
	}
}

// Tests that a single-value histogram renders its lone bar centred rather than at col 0.
func TestBuildHistogramColumns_StretchSingleValueCentred(t *testing.T) {
	hist := map[int]int{42: 7}
	counts, peak := buildHistogramColumns(hist, 42, 42, 6)
	if peak != 7 {
		t.Errorf("peak = %d, want 7", peak)
	}
	// width=6, centre = (6-1)/2 = 2.
	for i, c := range counts {
		want := 0
		if i == 2 {
			want = 7
		}
		if c != want {
			t.Errorf("counts[%d] = %d, want %d", i, c, want)
		}
	}
}

// Tests that when value range exceeds chart width, adjacent values aggregate into the same
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

// Tests that an empty histogram returns a zero-filled column slice and zero peak.
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

// Tests the visibility floor: any non-zero count renders at least one row, so long-tail
// buckets next to a much larger peak don't disappear into 0-row rounding.
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

// Tests that a wide range across the full chart width yields five ticks (min, lower quartile,
// midpoint, upper quartile, max) in left-to-right order.
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

// Tests that tick labels are centred on their columns, with edge ticks clipped inward and
// colliding interior ticks dropped so min and max always render.
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

// Tests the bottom-axis rendering: leading "+" under the y-axis, dashes for the body, and a
// "+" anchor at each interior tick position.
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

// Tests that a tall peak emits four label rows: the peak itself plus 3/4, 1/2, 1/4 of it.
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

// Tests colForValue across stretch, compress, and degenerate regimes so labels stay aligned
// with their bars.
func TestColForValue_MatchesBarPositions(t *testing.T) {
	// Stretch regime: 15 distinct values at slot=3 → natural width 43, bars at (v-min)*3.
	// minV lands at col 0 and maxV at col 42 = (rng-1)*slot.
	if got := colForValue(7, 7, 21, 43); got != 0 {
		t.Errorf("colForValue(min) = %d, want 0", got)
	}
	if got := colForValue(21, 7, 21, 43); got != 42 {
		t.Errorf("colForValue(max) = %d, want 42 ((rng-1)*slot)", got)
	}
	// Compress regime: col = (v-min)*width/rng. At v=50, min=0, max=119, width=60 → 50*60/120 = 25.
	if got := colForValue(50, 0, 119, 60); got != 25 {
		t.Errorf("colForValue(50, 0, 119, 60) = %d, want 25 (compress regime)", got)
	}
	// Degenerate rng<=1: single bar centred regardless of value.
	if got := colForValue(42, 42, 42, 60); got != 29 {
		t.Errorf("colForValue with min==max = %d, want 29 (centred at (width-1)/2)", got)
	}
}

// Tests that a tiny peak collapses to a single label rather than printing duplicates.
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

// Tests histChartWidth across stretch, compress fallback, single-value, and degenerate
// inputs.
func TestHistChartWidth(t *testing.T) {
	// Short range (15 values): stretch to 14*3+1 = 43 cols.
	if got := histChartWidth(15); got != 43 {
		t.Errorf("histChartWidth(15) = %d, want 43 ((rng-1)*slot+1)", got)
	}
	// Medium range (25 values): stretch to 24*3+1 = 73 cols.
	if got := histChartWidth(25); got != 73 {
		t.Errorf("histChartWidth(25) = %d, want 73 ((rng-1)*slot+1)", got)
	}
	// Compress regime: rng large enough that stretch would exceed histMaxStretchWidth.
	if got := histChartWidth(histMaxStretchWidth); got != histWidth {
		t.Errorf("histChartWidth(%d) = %d, want %d (compress)", histMaxStretchWidth, got, histWidth)
	}
	// Single-value range keeps the full width (there's nothing to spread).
	if got := histChartWidth(1); got != histWidth {
		t.Errorf("histChartWidth(1) = %d, want %d", got, histWidth)
	}
	// Degenerate zero / negative range still returns a sensible value so the caller can proceed.
	if got := histChartWidth(0); got != histWidth {
		t.Errorf("histChartWidth(0) = %d, want %d", got, histWidth)
	}
}

// Tests that xAxisTicks short-circuits to nil for non-positive width or empty range.
func TestXAxisTicks_GuardsDegenerateInputs(t *testing.T) {
	if ticks := xAxisTicks(0, 10, 0); ticks != nil {
		t.Errorf("width=0 returned %v, want nil", ticks)
	}
	if ticks := xAxisTicks(5, 3, 60); ticks != nil {
		t.Errorf("min>max (rng<=0) returned %v, want nil", ticks)
	}
}

// Tests that an interior tick colliding with min or max is dropped rather than overwriting
// the anchor.
func TestXAxisTickRow_CollisionDropsInteriorWinsMinMax(t *testing.T) {
	// Interior tick value deliberately placed at col 1 with a label wide enough to overlap
	// the 4-character "1234" label anchored at col 0. width=10 leaves just enough room for
	// the min and max labels to dominate.
	ticks := []xAxisTick{
		{col: 0, value: 1234},
		{col: 9, value: 5678},
		{col: 1, value: 9999}, // overlaps min ("1234" occupies cols 0-3)
	}
	got := xAxisTickRow(ticks, 10)
	if got[:4] != "1234" {
		t.Errorf("min label clobbered: first four chars = %q, want \"1234\"", got[:4])
	}
	if got[6:] != "5678" {
		t.Errorf("max label clobbered: last four chars = %q, want \"5678\"", got[6:])
	}
	// Interior tick "9999" is dropped; the middle of the buffer stays blank.
	if got[4:6] != "  " {
		t.Errorf("interior label should be dropped on collision; got middle = %q", got[4:6])
	}
}

// Tests that unionHistogramScale picks the wider min/max range and the larger of the two
// peaks so side-by-side charts share identical scales.
func TestUnionHistogramScale_StretchesAxesToCoverBoth(t *testing.T) {
	// d1: values 9..21, peak 32135 — narrower range, taller peak.
	// d2: values 7..27, peak 19723 — wider range, shorter peak.
	s1 := deck.Stats{
		Hands:     100000,
		Histogram: map[int]int{9: 1000, 12: 32135, 15: 24000, 18: 16000, 21: 8000},
	}
	s2 := deck.Stats{
		Hands:     100000,
		Histogram: map[int]int{7: 1, 12: 14000, 17: 19723, 22: 9000, 27: 1},
	}

	got := unionHistogramScale(s1, s2)
	if got.minV != 7 {
		t.Errorf("minV = %d, want 7 (smaller of the two mins)", got.minV)
	}
	if got.maxV != 27 {
		t.Errorf("maxV = %d, want 27 (larger of the two maxes)", got.maxV)
	}
	if got.peak != 32135 {
		t.Errorf("peak = %d, want 32135 (larger of the two peaks)", got.peak)
	}
}

// Tests centerLabel: centred when shorter than width, left-biased on odd leftover, passed
// through verbatim when longer.
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
