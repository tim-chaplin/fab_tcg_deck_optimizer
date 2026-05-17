package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

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
		ge := gameengine.New()
		self := &card.CardState{Card: tc.c}
		ge.ResolveChainStep(ge.Logger(), self)
		if got := ge.Value(); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
		if self.GrantedOverpower {
			t.Errorf("%s: GrantedOverpower should stay false when no aura", tc.c.Name())
		}
	}
}

// Tests that the AuraCreated flag flips self.GrantedOverpower.
func TestVantagePoint_AuraCreatedSetsOverpower(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
	self := &card.CardState{Card: cards.VantagePointRed{}}
	ge.ResolveChainStep(ge.Logger(), self)
	if got := ge.Value(); got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !self.GrantedOverpower {
		t.Errorf("GrantedOverpower should be set when AuraCreated is true")
	}
}
