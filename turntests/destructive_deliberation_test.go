package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func TestDestructiveDeliberation_PlayCreditsAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.DestructiveDeliberationRed{}, 5},
		{cards.DestructiveDeliberationYellow{}, 4},
		{cards.DestructiveDeliberationBlue{}, 3},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &card.CardState{Card: tc.c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if s.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), s.Value(), tc.want)
		}
		if len(self.OnHit) != 1 {
			t.Errorf("%s: OnHit = %d, want 1 (Ponder rider)", tc.c.Name(), len(self.OnHit))
		}
	}
}

func TestDestructiveDeliberation_OnHitCreatesPonder(t *testing.T) {
	for _, c := range []card.Card{
		cards.DestructiveDeliberationRed{},
		cards.DestructiveDeliberationYellow{},
		cards.DestructiveDeliberationBlue{},
	} {
		s := sim.TurnState{}
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		self.OnHit[0].Fire(&s, s.Logger(), self, &self.OnHit[0])
		if got := s.Ponders(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
