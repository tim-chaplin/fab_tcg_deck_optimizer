package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Plunder Run and Smashing Good Time grant +N{p} to the next scheduled attack
// action's BonusAttack iff self.FromArsenal is true.
func TestFromArsenalNextAttackBonus_GrantsOnArsenalCopyOnly(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.PlunderRunRed{}, 3},
		{cards.PlunderRunYellow{}, 2},
		{cards.PlunderRunBlue{}, 1},
		{notimpl.SmashingGoodTimeRed{}, 3},
		{notimpl.SmashingGoodTimeYellow{}, 2},
		{notimpl.SmashingGoodTimeBlue{}, 1},
	}
	for _, tc := range cases {
		// Hand-played copy: the bonus must NOT land on the queued attack action.
		handTarget := &card.CardState{Card: testutils.FakeRedAttack()}
		handState := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{handTarget}).Build()}
		handState.ResolveAttackStep(handState.Logger(), &card.CardState{Card: tc.c})
		if handTarget.BonusAttack != 0 {
			t.Errorf("%s: hand-play target BonusAttack = %d, want 0", tc.c.Name(), handTarget.BonusAttack)
		}
		// Arsenal-played copy: the bonus must land on the next attack action.
		arsenalTarget := &card.CardState{Card: testutils.FakeRedAttack()}
		arsenalState := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{arsenalTarget}).Build()}
		arsenalState.ResolveAttackStep(arsenalState.Logger(), &card.CardState{Card: tc.c, FromArsenal: true})
		if arsenalTarget.BonusAttack != tc.n {
			t.Errorf("%s: arsenal-play target BonusAttack = %d, want %d", tc.c.Name(), arsenalTarget.BonusAttack, tc.n)
		}
	}
}
