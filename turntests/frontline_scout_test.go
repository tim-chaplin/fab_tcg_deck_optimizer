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

// Tests that an arsenal-played Frontline Scout gains go again and extends with a hand attack.
func TestFrontlineScout_ArsenalPlayGoesAgainIntoHandAttack(t *testing.T) {
	handAttack := testutils.FakeRedAttack().WithPower(3)
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	state := gameengine.GameStateBuilder().SetArsenal(cards.FrontlineScoutRed{}).Build()
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{handAttack})
	if summary.Value != 6 {
		t.Errorf("Value = %d, want 6 (Frontline Scout 3 + hand attack 3 via arsenal go-again)", summary.Value)
	}
}

// Tests that a hand-played Frontline Scout gains no go again, so it cannot extend with a second attack.
func TestFrontlineScout_HandPlayDoesNotGoAgain(t *testing.T) {
	handAttack := testutils.FakeRedAttack().WithPower(3)
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{cards.FrontlineScoutRed{}, handAttack})
	if summary.Value != 3 {
		t.Errorf("Value = %d, want 3 (no arsenal go-again; only one attack resolves)", summary.Value)
	}
}
