package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Emissary of Tides cycles a hand card to the deck bottom and credits +2{p}.
func TestEmissaryOfTides_CyclesHandForBonus(t *testing.T) {
	cycled := testutils.FakeRedAttack()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{cards.EmissaryOfTidesRed{}, cycled})
	if summary.Value != 6 {
		t.Errorf("Value = %d, want 6 (printed 4 + 2{p} bonus)", summary.Value)
	}
}

// Tests that an empty hand suppresses the cycle and the bonus.
func TestEmissaryOfTides_EmptyHandSuppressesBonus(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{cards.EmissaryOfTidesRed{}})
	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (printed power; no cycle without a spare hand card)", summary.Value)
	}
}
