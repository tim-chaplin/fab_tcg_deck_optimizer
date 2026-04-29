package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestSpringLoad_BasePower: hand non-empty → printed power only, no rider.
func TestSpringLoad_BasePower(t *testing.T) {
	for _, c := range []sim.Card{SpringLoadRed{}, SpringLoadYellow{}, SpringLoadBlue{}} {
		s := sim.TurnState{Hand: []sim.Card{testutils.GenericAttack(0, 0)}}
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 2 {
			t.Errorf("%s: Play() with non-empty hand = %d, want 2 (printed power, no rider)", c.Name(), got)
		}
	}
}

// TestSpringLoad_EmptyHandFiresRider: every variant gains +3{p} when len(Hand) == 0.
func TestSpringLoad_EmptyHandFiresRider(t *testing.T) {
	for _, c := range []sim.Card{SpringLoadRed{}, SpringLoadYellow{}, SpringLoadBlue{}} {
		s := sim.TurnState{}
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 5 {
			t.Errorf("%s: Play() with empty hand = %d, want 5 (2 printed + 3 rider)", c.Name(), got)
		}
	}
}
