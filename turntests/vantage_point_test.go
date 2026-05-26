package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
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
		pc := &card.CardState{Card: tc.c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if got := ge.Value(); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
		if pc.GrantedOverpower {
			t.Errorf("%s: GrantedOverpower should stay false when no aura", tc.c.Name())
		}
	}
}

// Tests that the AuraCreated flag flips pc.GrantedOverpower.
func TestVantagePoint_AuraCreatedSetsOverpower(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
	pc := &card.CardState{Card: cards.VantagePointRed{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if got := ge.Value(); got != 7 {
		t.Errorf("Play() = %d, want 7", got)
	}
	if !pc.GrantedOverpower {
		t.Errorf("GrantedOverpower should be set when AuraCreated is true")
	}
}
