package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that PlayPrecondition passes when a cost-2-or-greater card sits in hand.
func TestDemolitionCrew_PreconditionPassesWithEligibleReveal(t *testing.T) {
	for _, c := range []card.Card{cards.DemolitionCrewRed{}, cards.DemolitionCrewYellow{}, cards.DemolitionCrewBlue{}} {
		ge := gameengine.New()
		ge.SetHand([]card.Card{testutils.FakeRedAttack().WithCost(2)})
		if ok := c.(card.PlayPrecondition).PlayPrecondition(ge, &card.CardState{Card: c}); !ok {
			t.Errorf("%s: PlayPrecondition with cost-2 card in hand returned false, want true", c.Name())
		}
	}
}

// Tests that PlayPrecondition fails when only sub-cost-2 cards sit in hand.
func TestDemolitionCrew_PreconditionFailsWithoutEligibleReveal(t *testing.T) {
	for _, c := range []card.Card{cards.DemolitionCrewRed{}, cards.DemolitionCrewYellow{}, cards.DemolitionCrewBlue{}} {
		ge := gameengine.New()
		ge.SetHand([]card.Card{testutils.FakeRedAttack().WithCost(1)})
		if ok := c.(card.PlayPrecondition).PlayPrecondition(ge, &card.CardState{Card: c}); ok {
			t.Errorf("%s: PlayPrecondition with no cost-2 card returned true, want false", c.Name())
		}
	}
}

// Tests that an empty hand fails the additional-cost check.
func TestDemolitionCrew_PreconditionFailsOnEmptyHand(t *testing.T) {
	ge := gameengine.New()
	if ok := (cards.DemolitionCrewRed{}).PlayPrecondition(ge, &card.CardState{Card: cards.DemolitionCrewRed{}}); ok {
		t.Errorf("PlayPrecondition with empty hand returned true, want false")
	}
}

// Tests that Play attacks for printed power once the precondition has been satisfied (a
// second cost-2 card in hand serves as the reveal) and the attack turn has enough pitch to fund
// the cost.
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
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{tc.c, testutils.FakeRedAttack().WithCost(2), testutils.FakeBlueResource()}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
		if summary.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), summary.Value, tc.want)
		}
	}
}
