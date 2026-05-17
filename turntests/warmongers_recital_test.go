package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero/heroes"
)

// Tests that Warmonger's Recital's OnHit recycle lands the buffed attack on the bottom
// of the deck. Hand: Bauble, Warmonger's Recital, Critical Strike. Empty starting deck.
// Optimal line: pitch Bauble {3} → Warmonger's Recital (cost 1, grants +3{p} and OnHit
// recycle to Critical Strike) → Critical Strike (now power 7, hits, OnHit pulls it from
// graveyard onto deck). End-of-turn deck should contain Critical Strike.
func TestWarmongersRecital_OnHitRecyclesToDeck(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, nil)
	hand := []deck.Card{
		cards.TitaniumBaubleBlue{},
		cards.WarmongersRecitalRed{},
		cards.CriticalStrikeYellow{},
	}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	csName := cards.CriticalStrikeYellow{}.DisplayName()
	if got.StartOfNextTurnDeck.NameCounts()[csName] == 0 {
		t.Fatalf("Critical Strike missing from end-of-turn deck; graveyard=%v", got.Graveyard)
	}
}
