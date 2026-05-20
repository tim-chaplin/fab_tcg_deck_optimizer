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

// Tests that Lead the Charge's action-point grant lets a second action card chain off it.
func TestLeadTheCharge_GrantsActionPointForNextAction(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{
		cards.LeadTheChargeRed{},
		testutils.BluePitch{},
		testutils.CostlyAttack{},
		testutils.CostlyAttack{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 6 {
		t.Fatalf("Value = %d, want 6 (two CostlyAttacks chained via go again + the rider's action point)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}

// Tests that Lead the Charge's grant fizzles when no action card follows.
func TestLeadTheCharge_NoExtraPointWithoutFollowingAction(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{
		cards.LeadTheChargeRed{},
		testutils.CostlyAttack{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (single CostlyAttack off Lead the Charge's own go again)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}
