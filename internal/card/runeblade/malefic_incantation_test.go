package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMaleficIncantation_NoFollowingAttackLingers: with no future attack in CardsRemaining, the
// aura sits in the arena — Play credits only the N-1 "future counters" as flat damage and the
// aura does not move to the graveyard this turn.
func TestMaleficIncantation_NoFollowingAttackLingers(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MaleficIncantationRed{}, 2},    // 3 - 1 (no same-turn pop)
		{MaleficIncantationYellow{}, 1}, // 2 - 1
		{MaleficIncantationBlue{}, 0},   // 1 - 1
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (no counter popped)", tc.c.Name(), s.Runechants)
		}
		if len(s.Graveyard) != 0 {
			t.Errorf("%s: Graveyard = %v, want empty (aura lingers)", tc.c.Name(), s.Graveyard)
		}
	}
}

// TestMaleficIncantation_FollowingAttackPopsCounter: a future attack in CardsRemaining ticks a
// counter this turn — CreateRunechant fires and the aura moves to the graveyard same turn so
// subsequent effects see it there.
func TestMaleficIncantation_FollowingAttackPopsCounter(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}},
		}
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if s.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (one counter popped)", tc.c.Name(), s.Runechants)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != tc.c.ID() {
			t.Errorf("%s: Graveyard = %v, want [%s] (destroyed same turn)", tc.c.Name(), s.Graveyard, tc.c.Name())
		}
	}
}

// TestMaleficIncantation_PlayNextTurnDestroys: the lingering aura is swept into the graveyard
// at the start of the subsequent turn via PlayNextTurn.
func TestMaleficIncantation_PlayNextTurnDestroys(t *testing.T) {
	cases := []card.DelayedPlay{
		MaleficIncantationRed{},
		MaleficIncantationYellow{},
		MaleficIncantationBlue{},
	}
	for _, c := range cases {
		var s card.TurnState
		got := c.PlayNextTurn(&s)
		if got.Damage != 0 {
			t.Errorf("%s: Damage = %d, want 0 (leave rider not modelled)", c.(card.Card).Name(), got.Damage)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != c.(card.Card).ID() {
			t.Errorf("%s: Graveyard = %v, want [%s]", c.(card.Card).Name(), s.Graveyard, c.(card.Card).Name())
		}
	}
}
