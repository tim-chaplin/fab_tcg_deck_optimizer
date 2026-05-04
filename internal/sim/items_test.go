package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that GoldTokenAbility.Play decrements Count and removes the entry at zero.
// Drives Play directly because the optimizer credits no Value for spending Gold (the
// drawn card's value lands in the carry state for the next turn) and so wouldn't
// reliably select the spend over the no-op.
func TestGoldAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	s := NewTurnState([]Card{testutils.RedAttack{}}, nil)
	s.Items = []Item{NewGoldItem(1)}
	(GoldTokenAbility{}).Play(s, &CardState{Card: GoldTokenAbility{}})
	if s.Gold() != 0 {
		t.Fatalf("Gold = %d after spending the only token, want 0", s.Gold())
	}
	if len(s.Items) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(s.Items))
	}
	if len(s.Hand) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(s.Hand))
	}
}

// Tests that spending one of multiple Gold tokens leaves the entry at decremented Count.
func TestGoldAbility_PlayDecrementsCountWhenMultiple(t *testing.T) {
	s := NewTurnState([]Card{testutils.RedAttack{}}, nil)
	s.Items = []Item{NewGoldItem(3)}
	(GoldTokenAbility{}).Play(s, &CardState{Card: GoldTokenAbility{}})
	if s.Gold() != 2 {
		t.Fatalf("Gold = %d after spending 1 of 3, want 2", s.Gold())
	}
}

// Tests that the eval cache fingerprints priorItems so calls with different gold
// counts don't collide.
func TestEvalCache_PriorItemsKeyedDistinctly(t *testing.T) {
	hand := []Card{testutils.RedAttack{}}
	ev := NewEvaluator()
	mp := Matchup{IncomingDamage: 0}
	_ = ev.BestWithTriggers(testutils.Hero{Intel: 4}, nil, hand, mp, nil, nil, nil, []Item{NewGoldItem(1)})
	_ = ev.BestWithTriggers(testutils.Hero{Intel: 4}, nil, hand, mp, nil, nil, nil, []Item{NewGoldItem(2)})
	stats := ev.CacheStats()
	// Two distinct item-count keys → both miss; neither hits.
	if stats.Hits != 0 {
		t.Errorf("hits = %d, want 0 (different gold counts must not collide)", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("misses = %d, want 2 (one per distinct item key)", stats.Misses)
	}
	// Repeating the second call should hit.
	_ = ev.BestWithTriggers(testutils.Hero{Intel: 4}, nil, hand, mp, nil, nil, nil, []Item{NewGoldItem(2)})
	stats = ev.CacheStats()
	if stats.Hits != 1 {
		t.Errorf("hits after repeat = %d, want 1 (matching item key should hit)", stats.Hits)
	}
}
