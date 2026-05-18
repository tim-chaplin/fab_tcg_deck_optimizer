package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Titanium Bauble's printed 3{d} blocks at block time — Resources have no
// Action subtype so they can't be played, but they can still defend like any non-Action
// hand card.
func TestResource_TitaniumBaubleBlocks(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.TitaniumBaubleBlue{}}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 5}, nil, hand)
	got := summary.Value
	if got != 3 {
		t.Fatalf("Value = %d, want 3 (Titanium blocks 3 of 5 incoming)", got)
	}
}

// Tests that a sole Resource card in hand at end of turn isn't promoted to arsenal.
func TestResource_DoesNotPromoteToArsenal(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.TitaniumBaubleBlue{}}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if summary.State.Arsenal() != nil {
		t.Fatalf("Arsenal = %v, want nil (Resource cards skip post-hoc promotion)", summary.State.Arsenal())
	}
}

// Tests that a sole Block-typed card (no Action / DR) isn't promoted to arsenal — the
// arsenal-in slot's only legal Block move is via Defense Reaction, so a pure Block card
// would lock there forever.
func TestArsenalPromotion_SkipsPureBlock(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.OnTheHorizonRed{}}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if summary.State.Arsenal() != nil {
		t.Fatalf("Arsenal = %v, want nil (pure Block cards skip post-hoc promotion)", summary.State.Arsenal())
	}
}
