package main

import "testing"

// TestBuildHistogramColumns_Stretch verifies the "range ≤ width" regime: when the histogram
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

// TestXAxisTickRow_FitsWithinWidth pins the tick-row layout: minV is left-justified at col 0,
// maxV is right-justified at the last column, spaces fill the middle.
func TestXAxisTickRow_FitsWithinWidth(t *testing.T) {
	got := xAxisTickRow(7, 21, 20)
	if len(got) != 20 {
		t.Errorf("len = %d, want 20 (chart width)", len(got))
	}
	if got[0] != '7' {
		t.Errorf("first char = %q, want '7' (left tick aligned at col 0)", got[0])
	}
	if got[len(got)-2:] != "21" {
		t.Errorf("last two chars = %q, want \"21\" (right tick aligned at last col)", got[len(got)-2:])
	}
}

// TestXAxisTickRow_TooNarrowFallsBack covers the degenerate case: when the two labels plus
// a separating space don't fit in width, the row collapses to a compact "min..max" form
// rather than overlapping labels.
func TestXAxisTickRow_TooNarrowFallsBack(t *testing.T) {
	got := xAxisTickRow(1000, 2000, 5)
	if got != "1000..2000" {
		t.Errorf("got %q, want %q", got, "1000..2000")
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
