package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMaleficIncantation_NoFollowUpAttackIsFlatN: no attack action follows Malefic in the
// chain, so no verse counter ticks this turn — credit flat n for the future-turn ticks.
func TestMaleficIncantation_NoFollowUpAttackIsFlatN(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (no same-turn tick)", tc.c.Name(), s.Runechants)
		}
	}
}

// TestMaleficIncantation_FollowUpAttackActionTicksOnce: an attack action card in the chain
// after Malefic triggers the "once per turn" clause — a live Runechant appears this turn
// and n-1 flat damage is credited for the remaining future-turn ticks. Total n either way.
func TestMaleficIncantation_FollowUpAttackActionTicksOnce(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubRunebladeAttack{}}}}
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if s.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (same-turn tick fired)", tc.c.Name(), s.Runechants)
		}
	}
}
