package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// TestSigilOfProtection_SetsAuraCreated verifies every variant flips AuraCreated and returns 0.
func TestSigilOfProtection_SetsAuraCreated(t *testing.T) {
	cases := []card.Card{notimplemented.SigilOfProtectionRed{}, notimplemented.SigilOfProtectionYellow{}, notimplemented.SigilOfProtectionBlue{}}
	for _, c := range cases {
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
