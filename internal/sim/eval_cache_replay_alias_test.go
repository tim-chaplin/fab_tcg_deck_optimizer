package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Tests that 20 shuffles of a seeded random Viserai deck produce identical TotalValue
// under cached vs uncached Evaluators.
func TestEvalCacheReplay_PromotesWinnerSlices(t *testing.T) {
	setupRNG := rand.New(rand.NewSource(6212487943906515912))
	baseline := deck.Random(heroes.Viserai, 40, 2, setupRNG, registry.Registry{})
	cachedStats := NewEvaluator().Evaluate(baseline.Copy(), 20, Matchup{IncomingPhysicalDamage: 7}, rand.New(rand.NewSource(99)))
	uncachedStats := NewEvaluatorWithoutCache().Evaluate(baseline.Copy(), 20, Matchup{IncomingPhysicalDamage: 7}, rand.New(rand.NewSource(99)))
	if cachedStats.TotalValue != uncachedStats.TotalValue {
		t.Fatalf("cached=%.0f uncached=%.0f", cachedStats.TotalValue, uncachedStats.TotalValue)
	}
}
