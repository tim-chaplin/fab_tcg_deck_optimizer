package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// Tests that cards.GoldToken.Play decrements Count and removes the entry at zero.
func TestGoldToken_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(token.NewGold(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.GoldToken{}})
	if ge.GoldCount() != 0 {
		t.Fatalf("Gold = %d after spending the only token, want 0", ge.GoldCount())
	}
	if len(ge.Items()) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(ge.Items()))
	}
	if h := ge.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests that spending one of multiple Gold tokens leaves the entry at decremented Count.
func TestGoldToken_PlayDecrementsCountWhenMultiple(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(token.NewGold(3))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.GoldToken{}})
	if ge.GoldCount() != 2 {
		t.Fatalf("Gold = %d after spending 1 of 3, want 2", ge.GoldCount())
	}
}

// Tests cards.SilverToken.Play decrement + draw behaviour.
func TestSilverToken_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(token.NewSilver(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SilverToken{}})
	if ge.SilverCount() != 0 {
		t.Fatalf("Silver = %d after spending the only token, want 0", ge.SilverCount())
	}
	if len(ge.Items()) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(ge.Items()))
	}
	if h := ge.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests cards.CopperToken.Play decrement + draw behaviour.
func TestCopperToken_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(token.NewCopper(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.CopperToken{}})
	if ge.CopperCount() != 0 {
		t.Fatalf("Copper = %d after spending the only token, want 0", ge.CopperCount())
	}
	if len(ge.Items()) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(ge.Items()))
	}
	if h := ge.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests CreateSilver and CreateCopper consolidate by token type.
func TestCreateSilverCopper_BumpsExistingEntry(t *testing.T) {
	ge := gameengine.New()
	ge.CreateSilver(2)
	ge.CreateSilver(1)
	ge.CreateCopper(1)
	if ge.SilverCount() != 3 {
		t.Errorf("Silver = %d, want 3 (2 + 1 consolidated)", ge.SilverCount())
	}
	if ge.CopperCount() != 1 {
		t.Errorf("Copper = %d, want 1", ge.CopperCount())
	}
	if got := len(ge.Items()); got != 2 {
		t.Errorf("Items entries = %d, want 2 (one Silver + one Copper)", got)
	}
}

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
