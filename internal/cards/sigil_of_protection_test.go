package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestSigilOfProtection_SetsAuraCreated verifies every variant flips AuraCreated and returns 0.
func TestSigilOfProtection_SetsAuraCreated(t *testing.T) {
	cases := []sim.Card{SigilOfProtectionRed{}, SigilOfProtectionYellow{}, SigilOfProtectionBlue{}}
	for _, c := range cases {
		s := sim.TurnState{}
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
