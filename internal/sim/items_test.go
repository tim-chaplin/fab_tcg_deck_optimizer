package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that GoldTokenAbility.Play decrements Count and removes the entry at zero.
func TestGoldAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(NewGoldItem(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: GoldTokenAbility{}})
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
func TestGoldAbility_PlayDecrementsCountWhenMultiple(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(NewGoldItem(3))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: GoldTokenAbility{}})
	if ge.GoldCount() != 2 {
		t.Fatalf("Gold = %d after spending 1 of 3, want 2", ge.GoldCount())
	}
}

// Tests SilverTokenAbility.Play decrement + draw behaviour.
func TestSilverAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(NewSilverItem(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: SilverTokenAbility{}})
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

// Tests CopperTokenAbility.Play decrement + draw behaviour.
func TestCopperAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := gameengine.New()
	ge.SetDeck(DeckOf(FakeRedAttack{}))
	ge.CreateItem(NewCopperItem(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: CopperTokenAbility{}})
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
	hero := FakeHero{Intel: 4}
	_ = ev.Best(nil, hand, mp, nil, Prior{Hero: hero, Items: []*Item{NewGoldItem(1)}})
	_ = ev.Best(nil, hand, mp, nil, Prior{Hero: hero, Items: []*Item{NewGoldItem(2)}})
	stats := ev.CacheStats()
	if stats.Hits != 0 {
		t.Errorf("hits = %d, want 0 (different gold counts must not collide)", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("misses = %d, want 2 (one per distinct item key)", stats.Misses)
	}
	_ = ev.Best(nil, hand, mp, nil, Prior{Hero: hero, Items: []*Item{NewGoldItem(2)}})
	stats = ev.CacheStats()
	if stats.Hits != 1 {
		t.Errorf("hits after repeat = %d, want 1 (matching item key should hit)", stats.Hits)
	}
}
