package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestMaleficIncantation_NoFollowingAttackIsFlatDamage(t *testing.T) {
	// No attack in CardsRemaining → fall back to flat N damage; no tokens on state.
	cases := []struct {
		c card.Card
		n int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != 0 || s.DelayedRunechants != 0 {
			t.Errorf("%s: Runechants=%d DelayedRunechants=%d, want both 0 (flat fallback)",
				tc.c.Name(), s.Runechants, s.DelayedRunechants)
		}
	}
}

func TestMaleficIncantation_ExactlyOneFollowingAttackDelaysOne(t *testing.T) {
	// Exactly one future attack in CardsRemaining → 1 token routed through DelayRunechants,
	// remaining N-1 are flat. Play still returns N total damage.
	cases := []struct {
		c card.Card
		n int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}}}
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.DelayedRunechants != 1 {
			t.Errorf("%s: DelayedRunechants = %d, want 1", tc.c.Name(), s.DelayedRunechants)
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (delayed, not live)", tc.c.Name(), s.Runechants)
		}
	}
}

func TestMaleficIncantation_MultipleFollowingAttacksIsFlatDamage(t *testing.T) {
	// Two future attacks → flat N (the "exactly one" branch doesn't trigger).
	s := card.TurnState{
		CardsRemaining: []*card.PlayedCard{
			{Card: stubRunebladeAttack{}},
			{Card: stubRunebladeAttack{}},
		},
	}
	if got := (MaleficIncantationRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
	if s.Runechants != 0 || s.DelayedRunechants != 0 {
		t.Errorf("Runechants=%d DelayedRunechants=%d, want both 0", s.Runechants, s.DelayedRunechants)
	}
}
