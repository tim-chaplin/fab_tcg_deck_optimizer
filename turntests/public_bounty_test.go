package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Public Bounty marks the opposing hero on Play.
func TestPublicBounty_MarksOpponent(t *testing.T) {
	for _, c := range []card.Card{cards.PublicBountyRed{}, cards.PublicBountyYellow{}, cards.PublicBountyBlue{}} {
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if !ge.OpponentMarked() {
			t.Errorf("%s: OpponentMarked = false after Play, want true", c.Name())
		}
	}
}

// Tests that Public Bounty grants the per-variant +N{p} bonus to the next IsAttack target.
func TestPublicBounty_GrantsBonusToNextAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.PublicBountyRed{}, 3},
		{cards.PublicBountyYellow{}, 2},
		{cards.PublicBountyBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack()}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
