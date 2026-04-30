package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Demolition Crew with a cost-2-or-greater card in hand attacks for printed power.
func TestDemolitionCrew_EligibleRevealAttacks(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{DemolitionCrewRed{}, 6},
		{DemolitionCrewYellow{}, 5},
		{DemolitionCrewBlue{}, 4},
	}
	for _, tc := range cases {
		s := sim.TurnState{Hand: []sim.Card{testutils.GenericAttack(2, 0)}}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that Demolition Crew with no cost-2-or-greater card in hand resolves as a no-op.
func TestDemolitionCrew_NoEligibleRevealNoOps(t *testing.T) {
	for _, c := range []sim.Card{DemolitionCrewRed{}, DemolitionCrewYellow{}, DemolitionCrewBlue{}} {
		s := sim.TurnState{Hand: []sim.Card{testutils.GenericAttack(1, 0)}}
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (additional cost unmet → no-op)", c.Name(), got)
		}
	}
}

// Tests that an empty Hand fails the additional-cost check.
func TestDemolitionCrew_EmptyHandNoOps(t *testing.T) {
	var s sim.TurnState
	(DemolitionCrewRed{}).Play(&s, &sim.CardState{Card: DemolitionCrewRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (empty hand can't reveal anything)", got)
	}
}
