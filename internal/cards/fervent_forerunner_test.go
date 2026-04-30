package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var ferventForerunnerVariants = []sim.Card{
	FerventForerunnerRed{},
	FerventForerunnerYellow{},
	FerventForerunnerBlue{},
}

// TestFerventForerunner_BaseGoAgainFalse pins the simplification: the only go-again trigger is
// "played from arsenal", which is gated on self.FromArsenal — printed GoAgain() must return
// false. Returning true would over-credit every sequence where it wasn't actually played from
// arsenal.
func TestFerventForerunner_BaseGoAgainFalse(t *testing.T) {
	for _, c := range ferventForerunnerVariants {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (arsenal-only go-again not modelled)", c.Name())
		}
	}
}

// Tests that the on-hit Opt 2 rider fires only when EffectiveAttack lands in the LikelyToHit
// window (1/4/7). Blue (printed power 1) hits, Red (3) and Yellow (2) miss.
func TestFerventForerunner_OnHitOptCreditsOnlyWhenInHitWindow(t *testing.T) {
	cases := []struct {
		c       sim.Card
		hitOpt  bool
		printed int
	}{
		{FerventForerunnerRed{}, false, 3},
		{FerventForerunnerYellow{}, false, 2},
		{FerventForerunnerBlue{}, true, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		want := tc.printed
		note := "no Opt — printed power outside 1/4/7"
		if tc.hitOpt {
			want += 2 * sim.OptValue
			note = "on-hit Opt 2"
		}
		if s.Value != want {
			t.Errorf("%s: Play() Value = %d, want %d (printed %d + %s)",
				tc.c.Name(), s.Value, want, tc.printed, note)
		}
	}
}

// Tests that a +1{p} grant (e.g. from a prior Force Sight) bumps Red's effective power into
// the 1/4/7 hit window, firing the on-hit Opt 2 rider.
func TestFerventForerunner_OnHitOptFiresWithBonusAttackInWindow(t *testing.T) {
	c := FerventForerunnerRed{}
	var s sim.TurnState
	c.Play(&s, &sim.CardState{Card: c, BonusAttack: 1})
	want := 3 + 1 + 2*sim.OptValue
	if s.Value != want {
		t.Errorf("Play() Value = %d, want %d (3 printed + 1 BonusAttack + Opt 2 on hit)", s.Value, want)
	}
}
