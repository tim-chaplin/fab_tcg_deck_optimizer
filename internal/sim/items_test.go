package sim

import (
	"testing"
)

// Tests that GoldTokenAbility.Play decrements Count and removes the entry at zero. Drives
// Play directly because the optimizer credits no Value for spending Gold.
func TestGoldAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	s := NewTurnStateFromCards([]Card{FakeRedAttack{}}, nil)
	s.Items = []Item{NewGoldItem(1)}
	ResolveChainStep(s, s.Logger(), &CardState{Card: GoldTokenAbility{}})
	if s.Gold() != 0 {
		t.Fatalf("Gold = %d after spending the only token, want 0", s.Gold())
	}
	if len(s.Items) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(s.Items))
	}
	if h := s.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests that spending one of multiple Gold tokens leaves the entry at decremented Count.
func TestGoldAbility_PlayDecrementsCountWhenMultiple(t *testing.T) {
	s := NewTurnStateFromCards([]Card{FakeRedAttack{}}, nil)
	s.Items = []Item{NewGoldItem(3)}
	ResolveChainStep(s, s.Logger(), &CardState{Card: GoldTokenAbility{}})
	if s.Gold() != 2 {
		t.Fatalf("Gold = %d after spending 1 of 3, want 2", s.Gold())
	}
}

// Tests SilverTokenAbility.Play decrement + draw behaviour, mirroring the Gold case.
func TestSilverAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	s := NewTurnStateFromCards([]Card{FakeRedAttack{}}, nil)
	s.Items = []Item{NewSilverItem(1)}
	ResolveChainStep(s, s.Logger(), &CardState{Card: SilverTokenAbility{}})
	if s.Silver() != 0 {
		t.Fatalf("Silver = %d after spending the only token, want 0", s.Silver())
	}
	if len(s.Items) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(s.Items))
	}
	if h := s.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests CopperTokenAbility.Play decrement + draw behaviour, mirroring the Gold case.
func TestCopperAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	s := NewTurnStateFromCards([]Card{FakeRedAttack{}}, nil)
	s.Items = []Item{NewCopperItem(1)}
	ResolveChainStep(s, s.Logger(), &CardState{Card: CopperTokenAbility{}})
	if s.Copper() != 0 {
		t.Fatalf("Copper = %d after spending the only token, want 0", s.Copper())
	}
	if len(s.Items) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(s.Items))
	}
	if h := s.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests CreateSilver and CreateCopper consolidate by token type.
func TestCreateSilverCopper_BumpsExistingEntry(t *testing.T) {
	s := NewTurnState(nil, nil)
	s.CreateSilver(2)
	s.CreateSilver(1)
	s.CreateCopper(1)
	if s.Silver() != 3 {
		t.Errorf("Silver = %d, want 3 (2 + 1 consolidated)", s.Silver())
	}
	if s.Copper() != 1 {
		t.Errorf("Copper = %d, want 1", s.Copper())
	}
	if got := len(s.Items); got != 2 {
		t.Errorf("Items entries = %d, want 2 (one Silver + one Copper)", got)
	}
}

// Tests that the eval cache fingerprints priorItems so calls with different gold
// counts don't collide.
func TestEvalCache_PriorItemsKeyedDistinctly(t *testing.T) {
	hand := []Card{FakeRedAttack{}}
	ev := NewEvaluator()
	mp := Matchup{IncomingDamage: 0}
	_ = ev.Best(FakeHero{Intel: 4}, nil, hand, mp, nil, NewTurnStateFromSpec(TurnStateSpec{Items: []Item{NewGoldItem(1)}}))
	_ = ev.Best(FakeHero{Intel: 4}, nil, hand, mp, nil, NewTurnStateFromSpec(TurnStateSpec{Items: []Item{NewGoldItem(2)}}))
	stats := ev.CacheStats()
	// Two distinct item-count keys → both miss; neither hits.
	if stats.Hits != 0 {
		t.Errorf("hits = %d, want 0 (different gold counts must not collide)", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("misses = %d, want 2 (one per distinct item key)", stats.Misses)
	}
	// Repeating the second call should hit.
	_ = ev.Best(FakeHero{Intel: 4}, nil, hand, mp, nil, NewTurnStateFromSpec(TurnStateSpec{Items: []Item{NewGoldItem(2)}}))
	stats = ev.CacheStats()
	if stats.Hits != 1 {
		t.Errorf("hits after repeat = %d, want 1 (matching item key should hit)", stats.Hits)
	}
}
