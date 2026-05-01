package cards_test

import (
	"testing"

	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Plunder Run and Smashing Good Time grant +N{p} to the next scheduled attack
// action's BonusAttack iff self.FromArsenal is true.
func TestFromArsenalNextAttackBonus_GrantsOnArsenalCopyOnly(t *testing.T) {
	cases := []struct {
		c sim.Card
		n int
	}{
		{notimpl.PlunderRunRed{}, 3},
		{notimpl.PlunderRunYellow{}, 2},
		{notimpl.PlunderRunBlue{}, 1},
		{notimpl.SmashingGoodTimeRed{}, 3},
		{notimpl.SmashingGoodTimeYellow{}, 2},
		{notimpl.SmashingGoodTimeBlue{}, 1},
	}
	for _, tc := range cases {
		// Hand-played copy: the bonus must NOT land on the queued attack action.
		handTarget := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
		handState := sim.TurnState{CardsRemaining: []*sim.CardState{handTarget}}
		tc.c.Play(&handState, &sim.CardState{Card: tc.c})
		if handTarget.BonusAttack != 0 {
			t.Errorf("%s: hand-play target BonusAttack = %d, want 0", tc.c.Name(), handTarget.BonusAttack)
		}
		// Arsenal-played copy: the bonus must land on the next attack action.
		arsenalTarget := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
		arsenalState := sim.TurnState{CardsRemaining: []*sim.CardState{arsenalTarget}}
		tc.c.Play(&arsenalState, &sim.CardState{Card: tc.c, FromArsenal: true})
		if arsenalTarget.BonusAttack != tc.n {
			t.Errorf("%s: arsenal-play target BonusAttack = %d, want %d", tc.c.Name(), arsenalTarget.BonusAttack, tc.n)
		}
	}
}
