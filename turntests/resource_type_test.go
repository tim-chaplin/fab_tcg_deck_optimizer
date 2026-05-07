package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Titanium Bauble's printed 3{d} is rejected at block time — Resource-typed
// cards can only Pitch, never Defend, even when their printed Defense would otherwise
// absorb incoming damage.
func TestResource_TitaniumBaubleCannotBlock(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.TitaniumBaubleBlue{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, sim.TurnState{}, hand).Value
	if got != 0 {
		t.Fatalf("Value = %d, want 0 (Resource Titanium Bauble can't block, hand has nothing else)", got)
	}
}

// Tests that a sole Resource card in hand at end of turn isn't promoted to arsenal.
func TestResource_DoesNotPromoteToArsenal(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.TitaniumBaubleBlue{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if got.StartOfNextTurnArsenal != nil {
		t.Fatalf("Arsenal = %v, want nil (Resource cards skip post-hoc promotion)", got.StartOfNextTurnArsenal)
	}
}
