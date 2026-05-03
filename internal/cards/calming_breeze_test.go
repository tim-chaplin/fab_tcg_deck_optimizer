package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Play credits the flat 3-damage prevention.
func TestCalmingBreeze_PreventsFlat3(t *testing.T) {
	s := sim.TurnState{IncomingDamage: 5}
	self := &sim.CardState{Card: CalmingBreezeRed{}}
	CalmingBreezeRed{}.Play(&s, self)
	if s.Value != 3 {
		t.Errorf("Value = %d, want 3", s.Value)
	}
	if s.IncomingDamage != 2 {
		t.Errorf("IncomingDamage = %d, want 2", s.IncomingDamage)
	}
}
