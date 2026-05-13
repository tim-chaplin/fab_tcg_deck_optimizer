package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that the on-hit 1{h} gain credits +1 on a likely-hit attack.
func TestLifeForALife_LikelyHitCreditsHeal(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	c := cards.LifeForALifeRed{}
	cs := &card.CardState{Card: c}
	s.ResolveChainStep(s.Logger(), cs)
	testutils.FireOnHitIfLikely(s, s.Logger(), cs)
	if got := s.Value(); got != 4+1 {
		t.Errorf("Red: Play() = %d, want 5 (4 likely to hit + 1 heal)", got)
	}
}

// Tests that the heal rider doesn't fire on blockable variants.
func TestLifeForALife_BlockableSuppressesHeal(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.LifeForALifeYellow{}, 3},
		{cards.LifeForALifeBlue{}, 2},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no heal)", tc.c.Name(), got, tc.want)
		}
	}
}
