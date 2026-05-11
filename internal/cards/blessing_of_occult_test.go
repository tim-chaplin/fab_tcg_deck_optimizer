package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Play flips AuraCreated, makes no runes this turn, and registers an aura with
// the per-variant Count for the deferred trigger.
func TestBlessingOfOccult_PlayCreatesAuraNoThisTurnRunes(t *testing.T) {
	cases := []struct {
		c         sim.Card
		wantCount int
	}{
		{BlessingOfOccultRed{}, 3},
		{BlessingOfOccultYellow{}, 2},
		{BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		sim.ResolveChainStep(&s, s.Logger(), &sim.CardState{Card: tc.c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune creation deferred to trigger)", tc.c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.Runechants() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (tokens are next-turn)", tc.c.Name(), s.Runechants())
		}
		if len(s.Auras()) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(s.Auras()))
		}
		if s.Auras()[0].TriggerType != sim.TriggerStartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", tc.c.Name(), s.Auras()[0].TriggerType)
		}
		if s.Auras()[0].Count != tc.wantCount {
			t.Errorf("%s: Count = %d, want %d", tc.c.Name(), s.Auras()[0].Count, tc.wantCount)
		}
	}
}

// TestBlessingOfOccult_TriggerHandlerCreatesNRunes: invoking the trigger's handler on a
// fresh TurnState creates N live Runechants and credits matching damage.
func TestBlessingOfOccult_TriggerHandlerCreatesNRunes(t *testing.T) {
	cases := []struct {
		c sim.Card
		n int
	}{
		{BlessingOfOccultRed{}, 3},
		{BlessingOfOccultYellow{}, 2},
		{BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		var play sim.TurnState
		sim.ResolveChainStep(&play, play.Logger(), &sim.CardState{Card: tc.c})
		next := sim.NewTurnStateFromCards(nil, nil)
		next.SetAuras(append(next.Auras(), play.Auras()[0]))
		next.SetCurrentAuraIdxForTesting(0)
		next.Auras()[0].Handler(next, next.Logger(), &next.Auras()[0].Trigger, &next.Auras()[0])
		if next.Value() != tc.n {
			t.Errorf("%s: handler Value = %d, want %d", tc.c.Name(), next.Value(), tc.n)
		}
		if next.Runechants() != tc.n {
			t.Errorf("%s: Runechants = %d, want %d (live tokens on next turn)",
				tc.c.Name(), next.Runechants(), tc.n)
		}
	}
}
