package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestSmashingGoodTime_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSmashingGoodTime_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{notimplemented.SmashingGoodTimeRed{}, notimplemented.SmashingGoodTimeYellow{}, notimplemented.SmashingGoodTimeBlue{}} {
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSmashingGoodTime_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSmashingGoodTime_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.SmashingGoodTimeRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSmashingGoodTime_NextAttackGrantsBonusAttack: arsenal-played copy with a queued
// attack action grants the per-variant bonus (Red +3, Yellow +2, Blue +1) onto the target's
// BonusAttack. Granter returns 0; the +N attributes to the buffed attack.
func TestSmashingGoodTime_NextAttackGrantsBonusAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{notimplemented.SmashingGoodTimeRed{}, 3},
		{notimplemented.SmashingGoodTimeYellow{}, 2},
		{notimplemented.SmashingGoodTimeBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		self := &card.CardState{Card: tc.c, FromArsenal: true}
		ge.ResolveChainStep(ge.Logger(), self)
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}

// TestSmashingGoodTime_HandPlayedFizzles: hand-played copy fails the from-arsenal gate.
func TestSmashingGoodTime_HandPlayedFizzles(t *testing.T) {
	for _, c := range []card.Card{notimplemented.SmashingGoodTimeRed{}, notimplemented.SmashingGoodTimeYellow{}, notimplemented.SmashingGoodTimeBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAttack(0, 0)}}).Build()}
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (hand-played)", c.Name(), got)
		}
	}
}
