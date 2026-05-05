package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestSigilOfFyendal_PlayRegistersStartOfTurnTrigger: Play flips AuraCreated and appends a
// start-of-turn Aura with Count=1 — no same-turn damage, the 1{h} gain is credited
// when the sim fires the trigger next turn.
func TestSigilOfFyendal_PlayRegistersStartOfTurnTrigger(t *testing.T) {
	var s sim.TurnState
	(SigilOfFyendalBlue{}).Play(&s, &sim.CardState{Card: SigilOfFyendalBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (1{h} gain deferred to trigger)", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
	if len(s.Auras) != 1 || s.Auras[0].TriggerType != sim.TriggerStartOfTurn {
		t.Errorf("Auras = %+v, want one TriggerStartOfTurn entry", s.Auras)
	}
	if s.Auras[0].Count != 1 {
		t.Errorf("Count = %d, want 1", s.Auras[0].Count)
	}
}

// TestSigilOfFyendal_TriggerHandlerCredits1Damage: the registered handler credits +1 damage
// (the 1{h} gain, valued 1-to-1 with damage).
func TestSigilOfFyendal_TriggerHandlerCredits1Damage(t *testing.T) {
	var s sim.TurnState
	(SigilOfFyendalBlue{}).Play(&s, &sim.CardState{Card: SigilOfFyendalBlue{}})
	var fire sim.TurnState
	s.Auras[0].Handler(&fire, &s.Auras[0])
	if fire.Value != 1 {
		t.Errorf("Handler Value = %d, want 1", fire.Value)
	}
}
