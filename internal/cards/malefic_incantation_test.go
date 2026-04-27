package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMaleficIncantation_PlayRegistersAttackActionTrigger: Play credits 0 same-turn damage —
// every rune comes from the trigger firing on each turn's first attack action. AuraCreated
// fires so same-turn aura-readers see Malefic. The registered trigger is TriggerAttackAction
// + OncePerTurn with Count=N (Red 3, Yellow 2, Blue 1).
func TestMaleficIncantation_PlayRegistersAttackActionTrigger(t *testing.T) {
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
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune comes from trigger, not Play)", tc.c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (trigger not yet fired)", tc.c.Name(), s.Runechants)
		}
		if len(s.AuraTriggers) != 1 {
			t.Fatalf("%s: AuraTriggers len = %d, want 1", tc.c.Name(), len(s.AuraTriggers))
		}
		tr := s.AuraTriggers[0]
		if tr.Type != card.TriggerAttackAction {
			t.Errorf("%s: trigger Type = %d, want TriggerAttackAction", tc.c.Name(), tr.Type)
		}
		if !tr.OncePerTurn {
			t.Errorf("%s: OncePerTurn = false, want true", tc.c.Name())
		}
		if tr.Count != tc.n {
			t.Errorf("%s: Count = %d, want %d (one per verse counter)", tc.c.Name(), tr.Count, tc.n)
		}
	}
}

// TestMaleficIncantation_HandlerCreatesOneRunechantPerFire: invoking the handler creates one
// live Runechant and credits 1 damage. Count tick + OncePerTurn gate are sim-managed, not
// the handler's job. chain.TriggeringCard is seeded to mimic the sim — the handler reads it
// to source-attribute the log entry it writes.
func TestMaleficIncantation_HandlerCreatesOneRunechantPerFire(t *testing.T) {
	for _, c := range []card.Card{MaleficIncantationRed{}, MaleficIncantationYellow{}, MaleficIncantationBlue{}} {
		var s card.TurnState
		c.Play(&s, &card.CardState{Card: c})
		chain := card.TurnState{TriggeringCard: c}
		got := s.AuraTriggers[0].Handler(&chain)
		if got != 1 {
			t.Errorf("%s: handler damage = %d, want 1", c.Name(), got)
		}
		if chain.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (handler creates one live rune)", c.Name(), chain.Runechants)
		}
	}
}

// TestMaleficIncantation_ImplementsAddsFutureValue pins the marker so the solver's
// beatsBest tiebreaker counts this card as future-value-adding — without it a lone Malefic
// loses to Held → arsenal promotion at equal current-turn Value.
func TestMaleficIncantation_ImplementsAddsFutureValue(t *testing.T) {
	for _, c := range []card.Card{MaleficIncantationRed{}, MaleficIncantationYellow{}, MaleficIncantationBlue{}} {
		if _, ok := c.(card.AddsFutureValue); !ok {
			t.Errorf("%s should implement card.AddsFutureValue", c.Name())
		}
	}
}
