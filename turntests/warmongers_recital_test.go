package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Warmonger's Recital's OnHit recycle lands the buffed attack on the deck bottom.
// Line: pitch Bauble {3} → Warmonger's Recital (cost 1, +3{p}, OnHit recycle) → Critical
// Strike (power 7, hits, OnHit pulls from graveyard onto deck). The turn-2 refill draws
// the lone deck card, so we look across Hand/Deck/Arsenal.
func TestWarmongersRecital_OnHitRecyclesToDeck(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, nil)
	hand := []card.Card{
		cards.TitaniumBaubleBlue{},
		cards.WarmongersRecitalRed{},
		cards.CriticalStrikeYellow{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if got := countAcrossSurfaces(summary.State, ids.CriticalStrikeYellow); got != 1 {
		t.Fatalf("Critical Strike total across turn-2 surfaces = %d, want 1 "+
			"(OnHit recycled it; graveyard=%v)", got, summary.State.Graveyard())
	}
}
