package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestSigilOfProtection_SetsAuraCreated verifies every variant flips AuraCreated and returns 0.
func TestSigilOfProtection_SetsAuraCreated(t *testing.T) {
	cases := []card.Card{notimplemented.SigilOfProtectionRed{}, notimplemented.SigilOfProtectionYellow{}, notimplemented.SigilOfProtectionBlue{}}
	for _, c := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
