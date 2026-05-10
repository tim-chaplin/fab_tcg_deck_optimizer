package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Play credits 0 immediately and registers an OncePerTurn TriggerAttackAction
// aura with the variant's Count.
func TestMaleficIncantation_PlayRegistersAttackActionTrigger(t *testing.T) {
	cases := []struct {
		c sim.Card
		n int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune comes from trigger, not Play)", tc.c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.Runechants() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (trigger not yet fired)", tc.c.Name(), s.Runechants())
		}
		if len(s.Auras) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(s.Auras))
		}
		tr := s.Auras[0]
		if tr.TriggerType != sim.TriggerAttackAction {
			t.Errorf("%s: trigger Type = %d, want TriggerAttackAction", tc.c.Name(), tr.TriggerType)
		}
		if !tr.OncePerTurn {
			t.Errorf("%s: OncePerTurn = false, want true", tc.c.Name())
		}
		if tr.Count != tc.n {
			t.Errorf("%s: Count = %d, want %d (one per verse counter)", tc.c.Name(), tr.Count, tc.n)
		}
	}
}

// Tests that one handler invocation creates one Runechant and credits 1 damage.
func TestMaleficIncantation_HandlerCreatesOneRunechantPerFire(t *testing.T) {
	for _, c := range []sim.Card{MaleficIncantationRed{}, MaleficIncantationYellow{}, MaleficIncantationBlue{}} {
		var s sim.TurnState
		c.Play(&s, &sim.CardState{Card: c})
		chain := sim.NewTurnStateFromCards(nil, nil)
		chain.TriggeringCard = c
		chain.Auras = append(chain.Auras, s.Auras[0])
		chain.SetCurrentAuraIdxForTesting(0)
		chain.Auras[0].Handler(chain, &chain.Auras[0].Trigger, &chain.Auras[0])
		if chain.Value != 1 {
			t.Errorf("%s: handler Value = %d, want 1", c.Name(), chain.Value)
		}
		if chain.Runechants() != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (handler creates one live rune)", c.Name(), chain.Runechants())
		}
	}
}
