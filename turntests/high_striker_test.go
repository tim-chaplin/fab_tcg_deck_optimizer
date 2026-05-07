package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests the High Striker → Critical Strike chain: High Striker's on-attack-action-hit
// trigger lands when Critical Strike (power 4, hits LikelyToHit) resolves, creating
// the printed per-pitch Copper count (6/4/2 for R/Y/B). Pitch Blue to fund {1}.
func TestHighStriker_TriggersOnNextAttackActionHit(t *testing.T) {
	cases := []struct {
		name       string
		striker    sim.Card
		wantCopper int
	}{
		{"Red", cards.HighStrikerRed{}, 6},
		{"Yellow", cards.HighStrikerYellow{}, 4},
		{"Blue", cards.HighStrikerBlue{}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := sim.New(heroes.Viserai{}, nil, fillerDeck())
			hand := []sim.Card{
				tc.striker,
				cards.CriticalStrikeYellow{},
				testutils.BluePitch{},
			}
			state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
			if state.Copper() != tc.wantCopper {
				t.Fatalf("Copper = %d, want %d (next attack hit fires the rider)\nBestLine: %s",
					state.Copper(), tc.wantCopper, formatBestLine(state.BestLine))
			}
		})
	}
}
