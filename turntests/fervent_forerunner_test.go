package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

var ferventForerunnerVariants = []card.Card{
	cards.FerventForerunnerRed{},
	cards.FerventForerunnerYellow{},
	cards.FerventForerunnerBlue{},
}

// Tests that printed GoAgain() is false; the only grant is the arsenal-gated rider.
func TestFerventForerunner_BaseGoAgainFalse(t *testing.T) {
	for _, c := range ferventForerunnerVariants {
		if c.GoAgain(nil) {
			t.Errorf("%s: GoAgain() = true, want false (arsenal-only go-again not modelled)", c.Name())
		}
	}
}

// Tests that the on-hit Opt 2 fires only when EffectiveAttack lands in the 1/4/7 window.
func TestFerventForerunner_OnHitOptFiresOnlyWhenInHitWindow(t *testing.T) {
	defer testutils.SwapCurrentHero(testutils.Hero{})()

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	cases := []struct {
		c       card.Card
		hitOpt  bool
		printed int
	}{
		{cards.FerventForerunnerRed{}, false, 3},
		{cards.FerventForerunnerYellow{}, false, 2},
		{cards.FerventForerunnerBlue{}, true, 1},
	}
	for _, tc := range cases {
		s := sim.NewTurnStateFromCards([]card.Card{a, b}, nil)
		cs := &card.CardState{Card: tc.c}
		sim.ResolveChainStep(s, s.Logger(), cs)
		testutils.FireOnHitIfLikely(s, s.Logger(), cs)
		if s.Value() != tc.printed {
			t.Errorf("%s: Play() Value = %d, want %d (printed power)",
				tc.c.Name(), s.Value(), tc.printed)
		}
		wantLogLen := 1
		if tc.hitOpt {
			wantLogLen = 2
		}
		if len(s.LogEntries()) != wantLogLen {
			t.Errorf("%s: Log len = %d, want %d", tc.c.Name(), len(s.LogEntries()), wantLogLen)
			continue
		}
		if tc.hitOpt {
			want := "Opted [a, b], put [a, b] on top, put [] on bottom"
			if got := s.LogEntries()[1].Text; got != want {
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
	c := cards.FerventForerunnerRed{}
	s := sim.NewTurnStateFromCards([]card.Card{a, b}, nil)
	cs := &card.CardState{Card: c, BonusAttack: 1}
	sim.ResolveChainStep(s, s.Logger(), cs)
	testutils.FireOnHitIfLikely(s, s.Logger(), cs)
	want := 3 + 1
	if s.Value() != want {
		t.Errorf("Play() Value = %d, want %d (3 printed + 1 BonusAttack)", s.Value(), want)
	}
	if len(s.LogEntries()) != 2 {
		t.Fatalf("Log len = %d, want 2 (chain step + Opted ...)", len(s.LogEntries()))
	}
	wantOpt := "Opted [a, b], put [a, b] on top, put [] on bottom"
	if got := s.LogEntries()[1].Text; got != wantOpt {
		t.Errorf("Opt log entry = %q, want %q", got, wantOpt)
	}
}
