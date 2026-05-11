package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Public Bounty marks the opposing hero on Play.
func TestPublicBounty_MarksOpponent(t *testing.T) {
	for _, c := range []card.Card{PublicBountyRed{}, PublicBountyYellow{}, PublicBountyBlue{}} {
		s := sim.TurnState{}
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: c})
		if !s.OpponentMarked() {
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
		{PublicBountyRed{}, 3},
		{PublicBountyYellow{}, 2},
		{PublicBountyBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{target}})
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
