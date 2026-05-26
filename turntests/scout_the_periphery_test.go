package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestScoutThePeriphery_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestScoutThePeriphery_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{cards.ScoutThePeripheryRed{}, cards.ScoutThePeripheryYellow{}, cards.ScoutThePeripheryBlue{}} {
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestScoutThePeriphery_NonAttackInRemainingFizzles: non-attack action (even from arsenal)
// fails the predicate — only attack actions count as the rider's target.
func TestScoutThePeriphery_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAction(), FromArsenal: true}}).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.ScoutThePeripheryRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestScoutThePeriphery_HandPlayedAttackFizzles: queued attack action that wasn't played from
// arsenal fails the rider's "next attack action card you play from arsenal" target gate.
func TestScoutThePeriphery_HandPlayedAttackFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAttack()}}).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.ScoutThePeripheryRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (target attack not from arsenal)", got)
	}
}

// Tests that the per-variant bonus lands on the next arsenal-played attack's BonusAttack
// (granter credits 0; bonus rides on the target).
func TestScoutThePeriphery_NextArsenalAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ScoutThePeripheryRed{}, 3},
		{cards.ScoutThePeripheryYellow{}, 2},
		{cards.ScoutThePeripheryBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack(), FromArsenal: true}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: granter credits %d, want 0 (bonus rides on target)", tc.c.Name(), got)
		}
		if got := target.BonusAttack; got != tc.want {
			t.Errorf("%s: target.BonusAttack = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
