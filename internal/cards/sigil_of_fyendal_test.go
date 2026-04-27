package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfFyendal_PlayRegistersStartOfTurnTrigger: Play flips AuraCreated and appends a
// start-of-turn AuraTrigger with Count=1 — no same-turn damage, the 1{h} gain is credited
// when the sim fires the trigger next turn.
func TestSigilOfFyendal_PlayRegistersStartOfTurnTrigger(t *testing.T) {
	var s card.TurnState
	(SigilOfFyendalBlue{}).Play(&s, &card.CardState{Card: SigilOfFyendalBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (1{h} gain deferred to trigger)", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
	if len(s.AuraTriggers) != 1 || s.AuraTriggers[0].Type != card.TriggerStartOfTurn {
		t.Errorf("AuraTriggers = %+v, want one TriggerStartOfTurn entry", s.AuraTriggers)
	}
	if s.AuraTriggers[0].Count != 1 {
		t.Errorf("Count = %d, want 1", s.AuraTriggers[0].Count)
	}
}

// TestSigilOfFyendal_TriggerHandlerCredits1Damage: the registered handler credits +1 damage
// (the 1{h} gain, valued 1-to-1 with damage).
func TestSigilOfFyendal_TriggerHandlerCredits1Damage(t *testing.T) {
	var s card.TurnState
	(SigilOfFyendalBlue{}).Play(&s, &card.CardState{Card: SigilOfFyendalBlue{}})
	if got := s.AuraTriggers[0].Handler(&card.TurnState{}); got != 1 {
		t.Errorf("Handler damage = %d, want 1", got)
	}
}

// TestSigilOfFyendal_AddsFutureValue pins the marker so the solver's beatsBest tiebreaker
// favours playing the sigil over Held → arsenal at equal Value.
func TestSigilOfFyendal_AddsFutureValue(t *testing.T) {
	var c card.Card = SigilOfFyendalBlue{}
	if _, ok := c.(card.AddsFutureValue); !ok {
		t.Error("SigilOfFyendalBlue should implement card.AddsFutureValue")
	}
}
