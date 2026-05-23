package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// Tests that the eval cache fingerprints priorItems so calls with different gold counts
// don't collide. Uses real cards / hero (not testutils fakes) so the cache key path actually
// runs — the cache bails out on InvalidCard / InvalidHero inputs.
func TestEvalCache_PriorItemsKeyedDistinctly(t *testing.T) {
	hand := []card.Card{cards.AetherSlashRed{}}
	ev := NewEvaluator()
	h := heroes.Viserai
	stateWithItems := func(items []*item.Item) *gameengine.GameState {
		b := gameengine.GameStateBuilder().SetHero(h)
		for _, it := range items {
			b.AddItem(it)
		}
		return b.Build()
	}
	_ = ev.Best(nil, hand, nil, stateWithItems([]*item.Item{token.NewGold(1)}))
	_ = ev.Best(nil, hand, nil, stateWithItems([]*item.Item{token.NewGold(2)}))
	stats := ev.CacheStats()
	if stats.Hits != 0 {
		t.Errorf("hits = %d, want 0 (different gold counts must not collide)", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("misses = %d, want 2 (one per distinct item key)", stats.Misses)
	}
	_ = ev.Best(nil, hand, nil, stateWithItems([]*item.Item{token.NewGold(2)}))
	stats = ev.CacheStats()
	if stats.Hits != 1 {
		t.Errorf("hits after repeat = %d, want 1 (matching item key should hit)", stats.Hits)
	}
}
