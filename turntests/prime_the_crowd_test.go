package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Play with no qualifying next attack returns zero.
func TestPrimeTheCrowd_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{cards.PrimeTheCrowdRed{}, cards.PrimeTheCrowdYellow{}, cards.PrimeTheCrowdBlue{}} {
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// Tests that a non-attack action in CardsRemaining fails the rider predicate.
func TestPrimeTheCrowd_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.PrimeTheCrowdRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// Tests that the first attack-action card in CardsRemaining receives the per-variant bonus.
func TestPrimeTheCrowd_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.PrimeTheCrowdRed{}, 4},
		{cards.PrimeTheCrowdYellow{}, 3},
		{cards.PrimeTheCrowdBlue{}, 2},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack()}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
