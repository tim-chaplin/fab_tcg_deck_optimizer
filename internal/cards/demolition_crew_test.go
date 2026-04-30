package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that PlayPrecondition passes when a cost-2-or-greater card sits in hand.
func TestDemolitionCrew_PreconditionPassesWithEligibleReveal(t *testing.T) {
	for _, c := range []sim.Card{DemolitionCrewRed{}, DemolitionCrewYellow{}, DemolitionCrewBlue{}} {
		s := sim.TurnState{Hand: []sim.Card{testutils.GenericAttack(2, 0)}}
		if ok := c.(sim.PlayPrecondition).PlayPrecondition(&s, &sim.CardState{Card: c}); !ok {
			t.Errorf("%s: PlayPrecondition with cost-2 card in hand returned false, want true", c.Name())
		}
	}
}

// Tests that PlayPrecondition fails when only sub-cost-2 cards sit in hand.
func TestDemolitionCrew_PreconditionFailsWithoutEligibleReveal(t *testing.T) {
	for _, c := range []sim.Card{DemolitionCrewRed{}, DemolitionCrewYellow{}, DemolitionCrewBlue{}} {
		s := sim.TurnState{Hand: []sim.Card{testutils.GenericAttack(1, 0)}}
		if ok := c.(sim.PlayPrecondition).PlayPrecondition(&s, &sim.CardState{Card: c}); ok {
			t.Errorf("%s: PlayPrecondition with no cost-2 card returned true, want false", c.Name())
		}
	}
}

// Tests that an empty hand fails the additional-cost check.
func TestDemolitionCrew_PreconditionFailsOnEmptyHand(t *testing.T) {
	var s sim.TurnState
	if ok := (DemolitionCrewRed{}).PlayPrecondition(&s, &sim.CardState{Card: DemolitionCrewRed{}}); ok {
		t.Errorf("PlayPrecondition with empty hand returned true, want false")
	}
}

// Tests that Play attacks for printed power once the precondition has been satisfied.
func TestDemolitionCrew_PlayAttacksForPrintedPower(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{DemolitionCrewRed{}, 6},
		{DemolitionCrewYellow{}, 5},
		{DemolitionCrewBlue{}, 4},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
