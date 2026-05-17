package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestPrimeTheCrowd_NoAttackReturnsZero: no qualifying next attack card → +4 rider fizzles.
func TestPrimeTheCrowd_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{notimplemented.PrimeTheCrowdRed{}, notimplemented.PrimeTheCrowdYellow{}, notimplemented.PrimeTheCrowdBlue{}} {
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestPrimeTheCrowd_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestPrimeTheCrowd_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.PrimeTheCrowdRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestPrimeTheCrowd_NextAttackReturnsBonus: first attack-action triggers the per-variant bonus
// (Red +4, Yellow +3, Blue +2).
func TestPrimeTheCrowd_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{notimplemented.PrimeTheCrowdRed{}, 4},
		{notimplemented.PrimeTheCrowdYellow{}, 3},
		{notimplemented.PrimeTheCrowdBlue{}, 2},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
