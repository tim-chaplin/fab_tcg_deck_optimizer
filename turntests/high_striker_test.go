package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests the High Striker → Critical Strike chain: High Striker's on-attack-action-hit
// trigger lands when Critical Strike (power 4, hits LikelyToHit) resolves, creating
// the printed per-pitch Copper count (6/4/2 for R/Y/B). Pitch Blue to fund {1}.
func TestHighStriker_TriggersOnNextAttackActionHit(t *testing.T) {
	cases := []struct {
		name       string
		striker    card.Card
		wantCopper int
	}{
		{"Red", cards.HighStrikerRed{}, 6},
		{"Yellow", cards.HighStrikerYellow{}, 4},
		{"Blue", cards.HighStrikerBlue{}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := deck.New(heroes.Viserai, nil, fillerDeck())
			hand := []card.Card{
				tc.striker,
				cards.CriticalStrikeYellow{},
				testutils.BluePitch{},
			}
			summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
			if summary.State.CopperCount() != tc.wantCopper {
				t.Fatalf("Copper = %d, want %d (next attack hit fires the rider)\nBestLine: %s",
					summary.State.CopperCount(), tc.wantCopper, formatBestLine(summary.BestLine))
			}
		})
	}
}
