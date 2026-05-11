package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that PlayPrecondition passes when a cost-2-or-greater card sits in hand.
func TestDemolitionCrew_PreconditionPassesWithEligibleReveal(t *testing.T) {
	for _, c := range []card.Card{DemolitionCrewRed{}, DemolitionCrewYellow{}, DemolitionCrewBlue{}} {
		var s sim.TurnState
		s.SetHandForTesting([]card.Card{testutils.GenericAttack(2, 0)})
		if ok := c.(sim.PlayPrecondition).PlayPrecondition(&s, &card.CardState{Card: c}); !ok {
			t.Errorf("%s: PlayPrecondition with cost-2 card in hand returned false, want true", c.Name())
		}
	}
}

// Tests that PlayPrecondition fails when only sub-cost-2 cards sit in hand.
func TestDemolitionCrew_PreconditionFailsWithoutEligibleReveal(t *testing.T) {
	for _, c := range []card.Card{DemolitionCrewRed{}, DemolitionCrewYellow{}, DemolitionCrewBlue{}} {
		var s sim.TurnState
		s.SetHandForTesting([]card.Card{testutils.GenericAttack(1, 0)})
		if ok := c.(sim.PlayPrecondition).PlayPrecondition(&s, &card.CardState{Card: c}); ok {
			t.Errorf("%s: PlayPrecondition with no cost-2 card returned true, want false", c.Name())
		}
	}
}

// Tests that an empty hand fails the additional-cost check.
func TestDemolitionCrew_PreconditionFailsOnEmptyHand(t *testing.T) {
	var s sim.TurnState
	if ok := (DemolitionCrewRed{}).PlayPrecondition(&s, &card.CardState{Card: DemolitionCrewRed{}}); ok {
		t.Errorf("PlayPrecondition with empty hand returned true, want false")
	}
}

// Tests that Play attacks for printed power once the precondition has been satisfied.
func TestDemolitionCrew_PlayAttacksForPrintedPower(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{DemolitionCrewRed{}, 6},
		{DemolitionCrewYellow{}, 5},
		{DemolitionCrewBlue{}, 4},
	}
	for _, tc := range cases {
		var s sim.TurnState
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
