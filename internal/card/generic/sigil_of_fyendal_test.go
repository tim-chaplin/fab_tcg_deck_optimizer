package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfFyendal_PlaySetsAuraCreated verifies Play flips AuraCreated and credits 0 — the
// 1{h} gain defers to PlayNextTurn since the aura only leaves at the start of the next action
// phase.
func TestSigilOfFyendal_PlaySetsAuraCreated(t *testing.T) {
	s := card.TurnState{}
	if got := (SigilOfFyendalBlue{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (health gain deferred to PlayNextTurn)", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
}

// TestSigilOfFyendal_PlayNextTurnGainsHealth verifies the deferred 1{h} credit.
func TestSigilOfFyendal_PlayNextTurnGainsHealth(t *testing.T) {
	s := card.TurnState{}
	got := (SigilOfFyendalBlue{}).PlayNextTurn(&s)
	if got.Damage != 1 {
		t.Errorf("Damage = %d, want 1 (1{h} gain on leave)", got.Damage)
	}
	if got.ToHand != nil {
		t.Errorf("ToHand = %v, want nil (Fyendal doesn't reveal)", got.ToHand)
	}
}
