package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that the on-hit Runechant rider fires only on likely-hit variants.
func TestMeatAndGreet_OnHitRunechantGatedByLikelyToHit(t *testing.T) {
	cases := []struct {
		c       card.Card
		wantDmg int
	}{
		{cards.MeatAndGreetRed{}, 4 + 1},
		{cards.MeatAndGreetYellow{}, 3},
		{cards.MeatAndGreetBlue{}, 2},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		self := &card.CardState{Card: tc.c}
		s.ResolveChainStep(s.Logger(), self)
		testutils.FireOnHitIfLikely(s, s.Logger(), self)
		if got := s.Value(); got != tc.wantDmg {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.wantDmg)
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true, want false (no prior arcane damage → no go again)", tc.c.Name())
		}
		// Card's printed GoAgain must also be false — the rider is the only source.
		if tc.c.GoAgain(nil) {
			t.Errorf("%s: GoAgain() = true, want false (rider is conditional, not printed)", tc.c.Name())
		}
	}
}

// Tests that ArcaneDamageDealt at Play time grants conditional go again.
func TestMeatAndGreet_ArcaneDamageDealtGrantsGoAgain(t *testing.T) {
	cases := []card.Card{
		cards.MeatAndGreetRed{},
		cards.MeatAndGreetYellow{},
		cards.MeatAndGreetBlue{},
	}
	for _, c := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetArcaneDamageDealt(true).Build()}
		self := &card.CardState{Card: c}
		s.ResolveChainStep(s.Logger(), self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (ArcaneDamageDealt → go again)", c.Name())
		}
	}
}
