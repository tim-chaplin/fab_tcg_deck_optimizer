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

// Tests that Lead the Charge's action-point grant lets a second action card play after it.
func TestLeadTheCharge_GrantsActionPointForNextAction(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{
		cards.LeadTheChargeRed{},
		testutils.FakeBlueResource(),
		testutils.FakeRedAttack().WithPower(3),
		testutils.FakeRedAttack().WithPower(3),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 6 {
		t.Fatalf("Value = %d, want 6 (two CostlyAttacks follow-up via the rider's action point)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}

// Tests that Lead the Charge's grant fizzles when no action card follows.
func TestLeadTheCharge_NoExtraPointWithoutFollowingAction(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{
		cards.LeadTheChargeRed{},
		testutils.FakeRedAttack().WithPower(3),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (single CostlyAttack off Lead the Charge's own go again)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}

// Tests that the granted action point stacks with the next action card's own go again — it
// is a real action point, not a redundant copy of that card's Go again keyword.
func TestLeadTheCharge_ActionPointStacksWithGoAgain(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{
		cards.LeadTheChargeRed{},
		testutils.FakeRedAttack().
			WithPower(3).
			WithGoAgain().
			WithTypes(card.TypeRuneblade),
		testutils.FakeBlueResource(),
		testutils.FakeRedAttack().WithPower(3),
		testutils.FakeRedAttack().WithPower(3),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 9 {
		t.Fatalf("Value = %d, want 9 (AttackWithPower + two CostlyAttacks; the rider's point stacks atop AttackWithPower's own go again)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}
