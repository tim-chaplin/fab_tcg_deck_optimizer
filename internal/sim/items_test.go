package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// Tests that the eval cache fingerprints priorItems so calls with different gold counts
// don't collide.
func TestEvalCache_PriorItemsKeyedDistinctly(t *testing.T) {
	hand := []card.Card{FakeRedAttack{}}
	ev := NewEvaluator()
	mp := Matchup{IncomingDamage: 0}
	h := FakeHero{Intel: 4}
	stateWithItems := func(items []*item.Item) *gameengine.GameState {
		b := gameengine.GameStateBuilder().SetHero(h)
		for _, it := range items {
			b.AddItem(it)
		}
		return b.Build()
	}
	_ = ev.Best(nil, hand, mp, nil, stateWithItems([]*item.Item{token.NewGold(1)}))
	_ = ev.Best(nil, hand, mp, nil, stateWithItems([]*item.Item{token.NewGold(2)}))
	stats := ev.CacheStats()
	if stats.Hits != 0 {
		t.Errorf("hits = %d, want 0 (different gold counts must not collide)", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("misses = %d, want 2 (one per distinct item key)", stats.Misses)
	}
	_ = ev.Best(nil, hand, mp, nil, stateWithItems([]*item.Item{token.NewGold(2)}))
	stats = ev.CacheStats()
	if stats.Hits != 1 {
		t.Errorf("hits after repeat = %d, want 1 (matching item key should hit)", stats.Hits)
	}
}
