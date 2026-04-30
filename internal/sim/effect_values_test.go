package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestLikelyDamageHits_OnlyAwkwardAmounts: without Dominate, 1 / 4 / 7 damage slip past
// typical blocks (cards are ~3 points of value, so opponents won't over-pay with a 3-block
// to soak 1 damage, etc.); everything else the sim treats as reliably blockable.
func TestLikelyDamageHits_OnlyAwkwardAmounts(t *testing.T) {
	for _, n := range []int{1, 4, 7} {
		if !LikelyDamageHits(n, false) {
			t.Errorf("LikelyDamageHits(%d, false) = false, want true (awkward amount)", n)
		}
	}
	for _, n := range []int{0, 2, 3, 5, 6, 8, 10} {
		if LikelyDamageHits(n, false) {
			t.Errorf("LikelyDamageHits(%d, false) = true, want false", n)
		}
	}
}

// TestLikelyDamageHits_DominateClearsFive: a Dominate attack caps the defender at one
// blocking card, so 5+ power slips past that one block. The awkward-amount rule still
// applies below 5.
func TestLikelyDamageHits_DominateClearsFive(t *testing.T) {
	for _, n := range []int{5, 6, 8, 10} {
		if !LikelyDamageHits(n, true) {
			t.Errorf("LikelyDamageHits(%d, true) = false, want true (dominate 5+)", n)
		}
	}
	// Still-blockable amounts under Dominate: 2 and 3 don't clear a single 3-block.
	for _, n := range []int{0, 2, 3} {
		if LikelyDamageHits(n, true) {
			t.Errorf("LikelyDamageHits(%d, true) = true, want false", n)
		}
	}
}

// TestLikelyToHit_FoldsEffectiveAttackAndDominate: the CardState-typed wrapper threads
// EffectiveAttack (printed + BonusAttack, clamp at 0) and EffectiveDominate (printed marker
// OR granted) into LikelyDamageHits. A +1 BonusAttack bumping a 3-power attack to 4 makes
// the rider fire even though the printed value alone would have been blocked; a -3
// BonusAttack flooring a 4-power attack at 0 turns the hit off; a Dominate grant on a
// 5-power attack clears the 5+ threshold.
func TestLikelyToHit_FoldsEffectiveAttackAndDominate(t *testing.T) {
	cases := []struct {
		name        string
		printed     int
		bonusAttack int
		grantedDom  bool
		printedDom  bool
		want        bool
	}{
		{"printed 3, no bonus → not in 1/4/7", 3, 0, false, false, false},
		{"printed 3, +1 bonus → 4 hits", 3, 1, false, false, true},
		{"printed 4, -10 bonus → clamped to 0, no hit", 4, -10, false, false, false},
		{"printed 5 with granted Dominate → 5+ clears one block", 5, 0, true, false, true},
		{"printed 5 with printed Dominate marker → 5+ clears one block", 5, 0, false, true, true},
		{"printed 5, no Dominate → blockable", 5, 0, false, false, false},
		{"printed 4, +1 bonus, granted Dominate → still in 1/4/7 OR 5+", 4, 1, true, false, true},
	}
	for _, tc := range cases {
		base := testutils.NewStubCard(tc.name).WithAttack(tc.printed)
		var c Card = base
		if tc.printedDom {
			c = testutils.DominatingStubCard{StubCard: base}
		}
		p := &CardState{Card: c, BonusAttack: tc.bonusAttack, GrantedDominate: tc.grantedDom}
		if got := LikelyToHit(p); got != tc.want {
			t.Errorf("%s: LikelyToHit() = %v, want %v", tc.name, got, tc.want)
		}
	}
}
