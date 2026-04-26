package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRunebloodIncantation_PlayRegistersStartOfTurnTriggerWithCountN: Play flips AuraCreated
// and registers a TriggerStartOfTurn AuraTrigger with Count=N (Red 3, Yellow 2, Blue 1). No
// same-turn damage credit — every Runechant lands on a real future-turn fire.
func TestRunebloodIncantation_PlayRegistersStartOfTurnTriggerWithCountN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{RunebloodIncantationRed{}, 3},
		{RunebloodIncantationYellow{}, 2},
		{RunebloodIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 0{
			t.Errorf("%s: Play() = %d, want 0 (every rune fires on a future turn)", tc.c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (no same-turn rune)", tc.c.Name(), s.Runechants)
		}
		if len(s.AuraTriggers) != 1 {
			t.Fatalf("%s: AuraTriggers len = %d, want 1", tc.c.Name(), len(s.AuraTriggers))
		}
		tr := s.AuraTriggers[0]
		if tr.Type != card.TriggerStartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", tc.c.Name(), tr.Type)
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
	for _, c := range []card.Card{RunebloodIncantationRed{}, RunebloodIncantationYellow{}, RunebloodIncantationBlue{}} {
		var play card.TurnState
		c.Play(&play, &card.CardState{Card: c})
		var fire card.TurnState
		got := play.AuraTriggers[0].Handler(&fire)
		if got != 1 {
			t.Errorf("%s: handler damage = %d, want 1", c.Name(), got)
		}
		if fire.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (one rune per fire)", c.Name(), fire.Runechants)
		}
	}
}
