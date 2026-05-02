package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// LikelyDamageHits is true at 1/4/7 or 5+ with dominate. Of printed powers (Red 3, Yellow 2,
// Blue 1) only Blue lands in the window unassisted, so only Blue flips GrantedGoAgain.
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

// BonusAttack pushes Red 3 into the hit window: +1 → 4 (1/4/7 path), +2 → 5 (dominate path).
// Both flip GrantedGoAgain.
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
