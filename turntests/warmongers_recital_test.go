package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Warmonger's Recital's OnHit recycle lands the buffed attack on the bottom
// of the deck. Hand: Bauble, Warmonger's Recital, Snatch. Single-card filler deck (so
// Snatch's own "draw a card on hit" rider pulls the filler off the top, leaving the
// recycled Snatch sitting on the bottom). Optimal line: pitch Bauble {3} → Warmonger's
// Recital (cost 1, grants +3{p} and OnHit recycle to Snatch) → Snatch (now power 7,
// hits, both OnHit handlers fire: Warmonger's recycle puts Snatch on deck bottom,
// Snatch's draw pulls the filler off the top). End-of-turn deck should contain Snatch.
func TestWarmongersRecital_OnHitRecyclesToDeck(t *testing.T) {
	filler := testutils.NewStubCard("filler")
	d := sim.New(heroes.Viserai{}, nil, []sim.Card{filler})
	hand := []sim.Card{
		cards.TitaniumBaubleBlue{},
		cards.WarmongersRecitalRed{},
		cards.SnatchRed{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	found := false
	for _, c := range got.StartOfNextTurnDeck {
		if c.ID() == (cards.SnatchRed{}).ID() {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Snatch missing from end-of-turn deck; deck=%v graveyard=%v", got.StartOfNextTurnDeck, got.Graveyard)
	}
}
