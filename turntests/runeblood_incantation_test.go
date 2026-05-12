package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestRunebloodIncantation_PlayRegistersStartOfTurnTriggerWithCountN: Play flips AuraCreated
// and registers a TriggerStartOfTurn Aura with Count=N (Red 3, Yellow 2, Blue 1). No
// same-turn damage credit — every Runechant lands on a real future-turn fire.
func TestRunebloodIncantation_PlayRegistersStartOfTurnTriggerWithCountN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.RunebloodIncantationRed{}, 3},
		{cards.RunebloodIncantationYellow{}, 2},
		{cards.RunebloodIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (every rune fires on a future turn)", tc.c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.RunechantCount() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (no same-turn rune)", tc.c.Name(), s.RunechantCount())
		}
		if len(s.Auras()) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(s.Auras()))
		}
		tr := s.Auras()[0]
		if tr.TriggerType != sim.TriggerStartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", tc.c.Name(), tr.TriggerType)
		}
		if tr.Count != tc.n {
			t.Errorf("%s: Count = %d, want %d (one per verse counter)", tc.c.Name(), tr.Count, tc.n)
		}
	}
}

// TestRunebloodIncantation_HandlerCreatesOneRunechantPerFire: each invocation of the handler
// creates exactly one live Runechant — the multi-fire behaviour comes from the sim ticking
// Count, not from the handler doing more work each call.
func TestRunebloodIncantation_HandlerCreatesOneRunechantPerFire(t *testing.T) {
	for _, c := range []card.Card{cards.RunebloodIncantationRed{}, cards.RunebloodIncantationYellow{}, cards.RunebloodIncantationBlue{}} {
		var play sim.TurnState
		sim.ResolveChainStep(&play, play.Logger(), &card.CardState{Card: c})
		fire := sim.NewTurnStateFromCards(nil, nil)
		fire.SetAuras(append(fire.Auras(), play.Auras()[0]))
		fire.SetCurrentAuraIdxForTesting(0)
		fire.FireAuraForTesting(0)
		if fire.Value() != 1 {
			t.Errorf("%s: handler Value = %d, want 1", c.Name(), fire.Value())
		}
		if fire.RunechantCount() != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (one rune per fire)", c.Name(), fire.RunechantCount())
		}
	}
}
