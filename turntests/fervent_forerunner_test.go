package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b}).Build()}
		cs := &card.CardState{Card: tc.c}
		ge.ResolveChainStep(ge.Logger(), cs)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
		if ge.Value() != tc.printed {
			t.Errorf("%s: Play() Value = %d, want %d (printed power)",
				tc.c.Name(), ge.Value(), tc.printed)
		}
		wantLogLen := 1
		if tc.hitOpt {
			wantLogLen = 2
		}
		if len(ge.LogEntries()) != wantLogLen {
			t.Errorf("%s: Log len = %d, want %d", tc.c.Name(), len(ge.LogEntries()), wantLogLen)
			continue
		}
		if tc.hitOpt {
			want := "Opted [a, b], put [a, b] on top, put [] on bottom"
			if got := ge.LogEntries()[1].Text; got != want {
				t.Errorf("%s: Opt log entry = %q, want %q", tc.c.Name(), got, want)
			}
		}
	}
}

// Tests that a +1{p} grant bumps Red's effective power into the 1/4/7 hit window, firing
// the on-hit Opt 2.
func TestFerventForerunner_OnHitOptFiresWithBonusAttackInWindow(t *testing.T) {

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	c := cards.FerventForerunnerRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b}).Build()}
	cs := &card.CardState{Card: c, BonusAttack: 1}
	ge.ResolveChainStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	want := 3 + 1
	if ge.Value() != want {
		t.Errorf("Play() Value = %d, want %d (3 printed + 1 BonusAttack)", ge.Value(), want)
	}
	if len(ge.LogEntries()) != 2 {
		t.Fatalf("Log len = %d, want 2 (chain step + Opted ...)", len(ge.LogEntries()))
	}
	wantOpt := "Opted [a, b], put [a, b] on top, put [] on bottom"
	if got := ge.LogEntries()[1].Text; got != wantOpt {
		t.Errorf("Opt log entry = %q, want %q", got, wantOpt)
	}
}
