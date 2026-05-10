package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func TestDestructiveDeliberation_PlayCreditsAttack(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{DestructiveDeliberationRed{}, 5},
		{DestructiveDeliberationYellow{}, 4},
		{DestructiveDeliberationBlue{}, 3},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &sim.CardState{Card: tc.c}
		tc.c.Play(&s, s.Logger(), self)
		if s.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), s.Value, tc.want)
		}
		if len(self.OnHit) != 1 {
			t.Errorf("%s: OnHit = %d, want 1 (Ponder rider)", tc.c.Name(), len(self.OnHit))
		}
	}
}

func TestDestructiveDeliberation_OnHitCreatesPonder(t *testing.T) {
	for _, c := range []sim.Card{
		DestructiveDeliberationRed{},
		DestructiveDeliberationYellow{},
		DestructiveDeliberationBlue{},
	} {
		s := sim.TurnState{}
		self := &sim.CardState{Card: c}
		c.Play(&s, s.Logger(), self)
		self.OnHit[0].Fire(&s, s.Logger(), self, &self.OnHit[0])
		if got := s.Ponders(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
