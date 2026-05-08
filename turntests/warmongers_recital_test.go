package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Warmonger's Recital's OnHit recycle lands the buffed attack on the bottom
// of the deck. Hand: Bauble, Warmonger's Recital, Critical Strike. Empty starting deck.
// Optimal line: pitch Bauble {3} → Warmonger's Recital (cost 1, grants +3{p} and OnHit
// recycle to Critical Strike) → Critical Strike (now power 7, hits, OnHit pulls it from
// graveyard onto deck). End-of-turn deck should contain Critical Strike.
func TestWarmongersRecital_OnHitRecyclesToDeck(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, nil)
	hand := []sim.Card{
		cards.TitaniumBaubleBlue{},
		cards.WarmongersRecitalRed{},
		cards.CriticalStrikeYellow{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	found := false
	for _, c := range got.StartOfNextTurnDeck {
		if c.ID() == (cards.CriticalStrikeYellow{}).ID() {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Critical Strike missing from end-of-turn deck; deck=%v graveyard=%v", got.StartOfNextTurnDeck, got.Graveyard)
	}
}
