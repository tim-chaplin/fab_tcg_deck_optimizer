package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestBlessingOfOccult_PlayCreatesAuraNoThisTurnRunes: Play flips AuraCreated so same-turn
// readers see an aura was created; no runes are made this turn (deferred to the trigger).
// The registered trigger is TriggerStartOfTurn with Count=1.
func TestBlessingOfOccult_PlayCreatesAuraNoThisTurnRunes(t *testing.T) {
	cases := []card.Card{BlessingOfOccultRed{}, BlessingOfOccultYellow{}, BlessingOfOccultBlue{}}
	for _, c := range cases {
		var s card.TurnState
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0{
			t.Errorf("%s: Play() = %d, want 0 (rune creation deferred to trigger)", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", c.Name())
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (tokens are next-turn)", c.Name(), s.Runechants)
		}
		if len(s.AuraTriggers) != 1 {
			t.Fatalf("%s: AuraTriggers len = %d, want 1", c.Name(), len(s.AuraTriggers))
		}
		if s.AuraTriggers[0].Type != card.TriggerStartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", c.Name(), s.AuraTriggers[0].Type)
		}
		if s.AuraTriggers[0].Count != 1 {
			t.Errorf("%s: Count = %d, want 1", c.Name(), s.AuraTriggers[0].Count)
		}
	}
}

// TestBlessingOfOccult_TriggerHandlerCreatesNRunes: invoking the trigger's handler on a
// fresh TurnState creates N live Runechants and credits matching damage.
func TestBlessingOfOccult_TriggerHandlerCreatesNRunes(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{BlessingOfOccultRed{}, 3},
		{BlessingOfOccultYellow{}, 2},
		{BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		var play card.TurnState
		tc.c.Play(&play, &card.CardState{Card: tc.c})
		var next card.TurnState
		got := play.AuraTriggers[0].Handler(&next)
		if got != tc.n {
			t.Errorf("%s: handler damage = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if next.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d (live tokens on next turn)",
				tc.c.Name(), next.Runechants, tc.n)
		}
	}
}
