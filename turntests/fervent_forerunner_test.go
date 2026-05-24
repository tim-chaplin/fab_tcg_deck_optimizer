package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
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

// Tests that Play credits the printed power; the on-hit Opt 2 only fires when
// EffectiveAttack lands in the 1/4/7 window, and its observable effect is the deck
// reshape, which TestFerventForerunner_OnHitOptReshapeIntoDeck pins separately via the
// post-Opt deck order.
func TestFerventForerunner_PlayCreditsPrintedPower(t *testing.T) {
	a, b := testutils.FakeRedAction().WithName("a"), testutils.FakeRedAction().WithName("b")
	cases := []struct {
		c       card.Card
		printed int
	}{
		{cards.FerventForerunnerRed{}, 3},
		{cards.FerventForerunnerYellow{}, 2},
		{cards.FerventForerunnerBlue{}, 1},
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
	}
}

// Tests that a +1{p} grant bumps Red's effective power so Play credits the buffed total.
func TestFerventForerunner_OnHitOptFiresWithBonusAttackInWindow(t *testing.T) {
	a, b := testutils.FakeRedAction().WithName("a"), testutils.FakeRedAction().WithName("b")
	c := cards.FerventForerunnerRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b}).Build()}
	cs := &card.CardState{Card: c, PerPerm: card.PerPerm{BonusAttack: 1}}
	ge.ResolveChainStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	want := 3 + 1
	if ge.Value() != want {
		t.Errorf("Play() Value = %d, want %d (3 printed + 1 BonusAttack)", ge.Value(), want)
	}
}
