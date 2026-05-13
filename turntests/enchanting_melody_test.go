package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestEnchantingMelody_SetsAuraCreated verifies every variant flips TurnState.AuraCreated (so
// downstream cards like Yinti Yanti / Runerager Swarm see the aura entering play) and reports 0
// direct damage — the aura's value is in the flag it leaves behind.
func TestEnchantingMelody_SetsAuraCreated(t *testing.T) {
	cases := []card.Card{notimplemented.EnchantingMelodyRed{}, notimplemented.EnchantingMelodyYellow{}, notimplemented.EnchantingMelodyBlue{}}
	for _, c := range cases {
		s := gameengine.NewFromState(nil)
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
