package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that PlayPrecondition passes when a cost-2-or-greater card sits in hand.
func TestDemolitionCrew_PreconditionPassesWithEligibleReveal(t *testing.T) {
	for _, c := range []card.Card{cards.DemolitionCrewRed{}, cards.DemolitionCrewYellow{}, cards.DemolitionCrewBlue{}} {
		s := gameengine.NewFromState(nil)
		s.SetHand([]card.Card{testutils.GenericAttack(2, 0)})
		if ok := c.(card.PlayPrecondition).PlayPrecondition(s, &card.CardState{Card: c}); !ok {
			t.Errorf("%s: PlayPrecondition with cost-2 card in hand returned false, want true", c.Name())
		}
	}
}

// Tests that PlayPrecondition fails when only sub-cost-2 cards sit in hand.
func TestDemolitionCrew_PreconditionFailsWithoutEligibleReveal(t *testing.T) {
	for _, c := range []card.Card{cards.DemolitionCrewRed{}, cards.DemolitionCrewYellow{}, cards.DemolitionCrewBlue{}} {
		s := gameengine.NewFromState(nil)
		s.SetHand([]card.Card{testutils.GenericAttack(1, 0)})
		if ok := c.(card.PlayPrecondition).PlayPrecondition(s, &card.CardState{Card: c}); ok {
			t.Errorf("%s: PlayPrecondition with no cost-2 card returned true, want false", c.Name())
		}
	}
}

// Tests that an empty hand fails the additional-cost check.
func TestDemolitionCrew_PreconditionFailsOnEmptyHand(t *testing.T) {
	s := gameengine.NewFromState(nil)
	if ok := (cards.DemolitionCrewRed{}).PlayPrecondition(s, &card.CardState{Card: cards.DemolitionCrewRed{}}); ok {
		t.Errorf("PlayPrecondition with empty hand returned true, want false")
	}
}

// Tests that Play attacks for printed power once the precondition has been satisfied.
func TestDemolitionCrew_PlayAttacksForPrintedPower(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.DemolitionCrewRed{}, 6},
		{cards.DemolitionCrewYellow{}, 5},
		{cards.DemolitionCrewBlue{}, 4},
	}
	for _, tc := range cases {
		s := gameengine.NewFromState(nil)
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
