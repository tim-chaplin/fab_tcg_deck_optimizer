package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that playing the Gold ability decrements Count and removes the entry at zero.
func TestGoldAbility_PlaysDecrementsAndDestroys(t *testing.T) {
	hand := []Card{testutils.RedAttack{}, testutils.RedAttack{}}
	priorItems := []Item{NewGoldItem(1)}
	got := BestWithTriggers(testutils.Hero{Intel: 4}, nil, hand, Matchup{IncomingDamage: 0}, nil, nil, nil, priorItems)
	if got.State.Gold() != 0 {
		t.Fatalf("Gold = %d after spending the only token, want 0", got.State.Gold())
	}
	if got.Value != DrawValue {
		t.Fatalf("Value = %d, want %d (Gold ability draws a card)", got.Value, DrawValue)
	}
}

// Tests that an unaffordable Gold ability leaves the token intact in the carry state.
func TestGoldAbility_UnspentTokenCarriesAcross(t *testing.T) {
	hand := []Card{testutils.RedAttack{}}
	priorItems := []Item{NewGoldItem(2)}
	got := BestWithTriggers(testutils.Hero{Intel: 4}, nil, hand, Matchup{IncomingDamage: 0}, nil, nil, nil, priorItems)
	if got.State.Gold() != 2 {
		t.Fatalf("Gold = %d (single red can't fund {2} ability), want 2 unchanged",
			got.State.Gold())
	}
}
