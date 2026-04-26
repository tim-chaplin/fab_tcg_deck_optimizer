package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestEnchantingMelody_SetsAuraCreated verifies every variant flips TurnState.AuraCreated (so
// downstream cards like Yinti Yanti / Runerager Swarm see the aura entering play) and reports 0
// direct damage — the aura's value is in the flag it leaves behind.
func TestEnchantingMelody_SetsAuraCreated(t *testing.T) {
	cases := []card.Card{EnchantingMelodyRed{}, EnchantingMelodyYellow{}, EnchantingMelodyBlue{}}
	for _, c := range cases {
		s := card.TurnState{}
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0{
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
