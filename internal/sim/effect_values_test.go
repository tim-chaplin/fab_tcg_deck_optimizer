package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that without Dominate only 1/4/7 damage are treated as likely-to-hit.
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

// Tests that Dominate makes 5+ power likely-to-hit; 0/2/3 stay blockable.
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

// Tests that LikelyToHit folds EffectiveAttack (printed + BonusAttack, clamp 0) and
// EffectiveDominate (printed marker OR granted) into LikelyDamageHits.
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
		base := NewFakeCard(tc.name).WithAttack(tc.printed)
		var c card.Card = base
		if tc.printedDom {
			c = DominatingFakeCard{FakeCard: base}
		}
		p := &card.CardState{Card: c, BonusAttack: tc.bonusAttack, GrantedDominate: tc.grantedDom}
		if got := LikelyToHit(p); got != tc.want {
			t.Errorf("%s: LikelyToHit() = %v, want %v", tc.name, got, tc.want)
		}
	}
}
