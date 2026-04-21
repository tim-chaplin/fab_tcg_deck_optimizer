package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestBlessingOfOccult_PlayCreatesAuraNoThisTurnRunes: Play just flips AuraCreated so
// same-turn readers see an aura was created; no runes are made this turn.
func TestBlessingOfOccult_PlayCreatesAuraNoThisTurnRunes(t *testing.T) {
	cases := []card.Card{BlessingOfOccultRed{}, BlessingOfOccultYellow{}, BlessingOfOccultBlue{}}
	for _, c := range cases {
		var s card.TurnState
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune creation deferred to PlayNextTurn)", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", c.Name())
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (tokens are next-turn)", c.Name(), s.Runechants)
		}
	}
}

// TestBlessingOfOccult_PlayNextTurnCreatesRunesAndGraveyardsSelf: at next turn's upkeep,
// Blessing creates N live Runechants on the new state and adds itself to the graveyard.
// The returned Damage matches the token count so cumulative Value picks up the credit.
func TestBlessingOfOccult_PlayNextTurnCreatesRunesAndGraveyardsSelf(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{BlessingOfOccultRed{}, 3},
		{BlessingOfOccultYellow{}, 2},
		{BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		dp := tc.c.(card.DelayedPlay)
		r := dp.PlayNextTurn(&s)
		if r.Damage != tc.n {
			t.Errorf("%s: Damage = %d, want %d", tc.c.Name(), r.Damage, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d (live tokens on next turn)", tc.c.Name(), s.Runechants, tc.n)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != tc.c.ID() {
			t.Errorf("%s: Graveyard = %v, want [self]", tc.c.Name(), s.Graveyard)
		}
	}
}
