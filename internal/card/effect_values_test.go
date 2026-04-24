package card

import "testing"

// TestLikelyToHit_OnlyAwkwardAmounts: without Dominate, 1 / 4 / 7 damage slip past typical
// blocks (cards are ~3 points of value, so opponents won't over-pay with a 3-block to soak 1
// damage, etc.); everything else the sim treats as reliably blockable.
func TestLikelyToHit_OnlyAwkwardAmounts(t *testing.T) {
	for _, n := range []int{1, 4, 7} {
		if !LikelyToHit(n, false) {
			t.Errorf("LikelyToHit(%d, false) = false, want true (awkward amount)", n)
		}
	}
	for _, n := range []int{0, 2, 3, 5, 6, 8, 10} {
		if LikelyToHit(n, false) {
			t.Errorf("LikelyToHit(%d, false) = true, want false", n)
		}
	}
}

// TestLikelyToHit_DominateClearsFive: a Dominate attack caps the defender at one blocking
// card, so 5+ power slips past that one block. The awkward-amount rule still applies below 5.
func TestLikelyToHit_DominateClearsFive(t *testing.T) {
	for _, n := range []int{5, 6, 8, 10} {
		if !LikelyToHit(n, true) {
			t.Errorf("LikelyToHit(%d, true) = false, want true (dominate 5+)", n)
		}
	}
	// Still-blockable amounts under Dominate: 2 and 3 don't clear a single 3-block.
	for _, n := range []int{0, 2, 3} {
		if LikelyToHit(n, true) {
			t.Errorf("LikelyToHit(%d, true) = true, want false", n)
		}
	}
}
