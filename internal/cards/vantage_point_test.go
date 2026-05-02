package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that with no aura played or created the printed power is credited and Overpower stays false.
func TestVantagePoint_BaseDamageNoAura(t *testing.T) {
	cases := []struct {
		c    sim.Card
		base int
	}{
		{VantagePointRed{}, 7},
		{VantagePointYellow{}, 6},
		{VantagePointBlue{}, 5},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
		if s.Overpower {
			t.Errorf("%s: Overpower should stay false when no aura", tc.c.Name())
		}
	}
}

// Tests that an aura already played this turn flips s.Overpower.
func TestVantagePoint_AuraPlayedSetsOverpower(t *testing.T) {
	s := sim.TurnState{CardsPlayed: []sim.Card{testutils.Aura{}}}
	(VantagePointRed{}).Play(&s, &sim.CardState{Card: VantagePointRed{}})
	if got := s.Value; got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !s.Overpower {
		t.Errorf("Overpower should be set when an aura was played")
	}
}

// Tests that the AuraCreated flag also flips s.Overpower.
func TestVantagePoint_AuraCreatedSetsOverpower(t *testing.T) {
	s := sim.TurnState{AuraCreated: true}
	(VantagePointRed{}).Play(&s, &sim.CardState{Card: VantagePointRed{}})
	if got := s.Value; got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !s.Overpower {
		t.Errorf("Overpower should be set when AuraCreated is true")
	}
}
