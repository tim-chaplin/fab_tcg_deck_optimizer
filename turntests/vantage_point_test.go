package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that with no aura played or created the printed power is credited and
// GrantedOverpower stays false.
func TestVantagePoint_BaseDamageNoAura(t *testing.T) {
	cases := []struct {
		c    card.Card
		base int
	}{
		{cards.VantagePointRed{}, 7},
		{cards.VantagePointYellow{}, 6},
		{cards.VantagePointBlue{}, 5},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.NewState()}
		self := &card.CardState{Card: tc.c}
		s.ResolveChainStep(s.Logger(), self)
		if got := s.Value(); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
		if self.GrantedOverpower {
			t.Errorf("%s: GrantedOverpower should stay false when no aura", tc.c.Name())
		}
	}
}

// Tests that the AuraCreated flag flips self.GrantedOverpower.
func TestVantagePoint_AuraCreatedSetsOverpower(t *testing.T) {
	s := gameengine.NewFromSpec(gameengine.Spec{AuraCreated: true})
	self := &card.CardState{Card: cards.VantagePointRed{}}
	s.ResolveChainStep(s.Logger(), self)
	if got := s.Value(); got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !self.GrantedOverpower {
		t.Errorf("GrantedOverpower should be set when AuraCreated is true")
	}
}
