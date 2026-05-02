package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Overload's on-hit go-again rider is modelled eagerly via sim.LikelyToHit.
// LikelyDamageHits is true at 1/4/7 power or at 5+ with dominate. Printed
// powers are Red 3, Yellow 2, Blue 1 — only Blue (n==1) lands in the window
// without help, so only Blue flips GrantedGoAgain on a clean Play.
func TestOverload_OnHitGoAgainEagerByLikelyToHit(t *testing.T) {
	cases := []struct {
		c       sim.Card
		wantDmg int
		wantGA  bool
	}{
		{OverloadRed{}, 3, false},
		{OverloadYellow{}, 2, false},
		{OverloadBlue{}, 1, true},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &sim.CardState{Card: tc.c}
		tc.c.Play(&s, self)
		if got := s.Value; got != tc.wantDmg {
			t.Errorf("%s: Play() Value = %d, want %d", tc.c.Name(), got, tc.wantDmg)
		}
		if self.GrantedGoAgain != tc.wantGA {
			t.Errorf("%s: GrantedGoAgain = %v, want %v", tc.c.Name(), self.GrantedGoAgain, tc.wantGA)
		}
		if tc.c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (rider is conditional, not printed)", tc.c.Name())
		}
	}
}

// A +2{p} BonusAttack on Red Overload bumps it from 3 → 5 (dominate window),
// flipping GrantedGoAgain. A +1 bonus stops at 4 — also a hit window — so
// likewise. Cover both the dominate-5+ and the 1/4/7 paths.
func TestOverload_BonusAttackPushesIntoHitWindow(t *testing.T) {
	cases := []struct {
		bonus int
	}{
		{1}, // Red 3 + 1 = 4 → hit window
		{2}, // Red 3 + 2 = 5 with dominate → hit window
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &sim.CardState{Card: OverloadRed{}, BonusAttack: tc.bonus}
		OverloadRed{}.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("Red + BonusAttack %d: GrantedGoAgain = false, want true", tc.bonus)
		}
	}
}
