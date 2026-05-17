package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that a plain-block card lands in the graveyard after defending.
func TestPlainBlock_LandsInGraveyard(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.OnTheHorizonRed{}}
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, hand)
	if extras.Value != 4 {
		t.Fatalf("Value = %d, want 4 (On the Horizon Red blocks 4)", extras.Value)
	}
	target := cards.OnTheHorizonRed{}.ID()
	for _, c := range gs.Graveyard() {
		if c.ID() == target {
			return
		}
	}
	t.Fatalf("On the Horizon [R] missing from graveyard after blocking; got %v", gs.Graveyard())
}
