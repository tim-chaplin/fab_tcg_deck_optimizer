package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMaleficIncantation_NoFollowUpAttackIsFlatN: no attack action follows Malefic in the
// chain, so no verse counter ticks this turn — credit flat n for the future-turn ticks and
// leave the aura in play (it heads to next turn's graveyard via PlayNextTurn).
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
		if len(s.Graveyard) != 0 {
			t.Errorf("%s: Graveyard = %v, want empty (aura persists until next turn)",
				tc.c.Name(), s.Graveyard)
		}
	}
}

// TestMaleficIncantation_FollowUpAttackActionTicksOnce: an attack action card in the chain
// after Malefic triggers the "once per turn" clause — a live Runechant appears this turn
// and the total damage credited is n (1 rune + n-1 flat for the remaining ticks).
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

// TestMaleficIncantationBlue_FollowUpAttackGraveyardsImmediately: Blue starts with a single
// verse counter, so the same-turn tick takes it to zero and the aura lands in this turn's
// graveyard right away. Red/Yellow still have counters left, so they stay in play.
func TestMaleficIncantationBlue_FollowUpAttackGraveyardsImmediately(t *testing.T) {
	cases := []struct {
		c             card.Card
		wantGraveyard bool
	}{
		{MaleficIncantationRed{}, false},
		{MaleficIncantationYellow{}, false},
		{MaleficIncantationBlue{}, true},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubRunebladeAttack{}}}}
		tc.c.Play(&s, &card.CardState{})
		inGrav := len(s.Graveyard) == 1 && s.Graveyard[0].ID() == tc.c.ID()
		if inGrav != tc.wantGraveyard {
			t.Errorf("%s: in graveyard = %v, want %v (Graveyard=%v)",
				tc.c.Name(), inGrav, tc.wantGraveyard, s.Graveyard)
		}
	}
}

// TestMaleficIncantation_PlayNextTurnGraveyardsSelf: the aura heads to next turn's graveyard
// when it still had counters left at end of this turn (any variant without a same-turn tick,
// or Red/Yellow even with one). The callback only graveyards; damage credit is already on
// the previous turn's Play return.
func TestMaleficIncantation_PlayNextTurnGraveyardsSelf(t *testing.T) {
	cases := []card.Card{
		MaleficIncantationRed{},
		MaleficIncantationYellow{},
		MaleficIncantationBlue{},
	}
	for _, c := range cases {
		var s card.TurnState
		dp := c.(card.DelayedPlay)
		r := dp.PlayNextTurn(&s)
		if r.Damage != 0 {
			t.Errorf("%s: Damage = %d, want 0 (damage was credited on Play)", c.Name(), r.Damage)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != c.ID() {
			t.Errorf("%s: Graveyard = %v, want [self]", c.Name(), s.Graveyard)
		}
	}
}
