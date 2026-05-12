package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Play credits the flat 3-damage prevention.
func TestCalmingBreeze_PreventsFlat3(t *testing.T) {
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{IncomingDamage: 5})
	self := &card.CardState{Card: cards.CalmingBreezeRed{}}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if s.Value() != 3 {
		t.Errorf("Value = %d, want 3", s.Value())
	}
	if s.IncomingDamage() != 2 {
		t.Errorf("IncomingDamage = %d, want 2", s.IncomingDamage())
	}
}
