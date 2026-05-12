package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that a plain-block card lands in the graveyard after defending.
func TestPlainBlock_LandsInGraveyard(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.OnTheHorizonRed{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, gameengine.Spec{}, hand)
	if state.Value != 4 {
		t.Fatalf("Value = %d, want 4 (On the Horizon Red blocks 4)", state.Value)
	}
	target := cards.OnTheHorizonRed{}.ID()
	for _, c := range state.Graveyard {
		if c.ID() == target {
			return
		}
	}
	t.Fatalf("On the Horizon [R] missing from graveyard after blocking; got %v", state.Graveyard)
}
