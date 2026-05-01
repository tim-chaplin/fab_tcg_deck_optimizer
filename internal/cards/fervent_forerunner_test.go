package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

var ferventForerunnerVariants = []sim.Card{
	FerventForerunnerRed{},
	FerventForerunnerYellow{},
	FerventForerunnerBlue{},
}

// TestFerventForerunner_BaseGoAgainFalse pins printed GoAgain() = false; the only grant is the
// arsenal-gated rider on self.FromArsenal.
func TestFerventForerunner_BaseGoAgainFalse(t *testing.T) {
	for _, c := range ferventForerunnerVariants {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (arsenal-only go-again not modelled)", c.Name())
		}
	}
}

// Tests that the on-hit Opt 2 fires only when EffectiveAttack lands in the 1/4/7 window.
func TestFerventForerunner_OnHitOptFiresOnlyWhenInHitWindow(t *testing.T) {
	defer testutils.SwapCurrentHero(testutils.Hero{})()

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
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
		s := sim.NewTurnState([]sim.Card{a, b}, nil)
		tc.c.Play(s, &sim.CardState{Card: tc.c})
		if s.Value != tc.printed {
			t.Errorf("%s: Play() Value = %d, want %d (printed power)",
				tc.c.Name(), s.Value, tc.printed)
		}
		wantLogLen := 1
		if tc.hitOpt {
			wantLogLen = 2
		}
		if len(s.Log) != wantLogLen {
			t.Errorf("%s: Log len = %d, want %d", tc.c.Name(), len(s.Log), wantLogLen)
			continue
		}
		if tc.hitOpt {
			want := "Opted [a, b], put [a, b] on top, put [] on bottom"
			if got := s.Log[1].Text; got != want {
				t.Errorf("%s: Opt log entry = %q, want %q", tc.c.Name(), got, want)
			}
		}
	}
}

// Tests that a +1{p} grant bumps Red's effective power into the 1/4/7 hit window, firing
// the on-hit Opt 2.
func TestFerventForerunner_OnHitOptFiresWithBonusAttackInWindow(t *testing.T) {
	defer testutils.SwapCurrentHero(testutils.Hero{})()

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	c := FerventForerunnerRed{}
	s := sim.NewTurnState([]sim.Card{a, b}, nil)
	c.Play(s, &sim.CardState{Card: c, BonusAttack: 1})
	want := 3 + 1
	if s.Value != want {
		t.Errorf("Play() Value = %d, want %d (3 printed + 1 BonusAttack)", s.Value, want)
	}
	if len(s.Log) != 2 {
		t.Fatalf("Log len = %d, want 2 (chain step + Opted ...)", len(s.Log))
	}
	wantOpt := "Opted [a, b], put [a, b] on top, put [] on bottom"
	if got := s.Log[1].Text; got != wantOpt {
		t.Errorf("Opt log entry = %q, want %q", got, wantOpt)
	}
}
