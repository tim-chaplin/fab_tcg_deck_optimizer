package e2etest

// End-to-end tests for Yinti Yanti's "while you control an aura, +1{p}" rider seeing
// auras created by defenders earlier in the same turn.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Reduce to Runechant defends, creating a Runechant; Yinti Yanti's Play then sees the
// aura and credits +1{p}. Carryover Runechant funds Reduce's cost.
func TestYintiYanti_SeesRunechantFromReduceInDefense(t *testing.T) {
	hand := []sim.Card{cards.YintiYantiRed{}, cards.ReduceToRunechantRed{}}
	priorAuras := []sim.Aura{sim.NewRunechantAura(1)}
	got := sim.BestWithTriggers(heroes.Viserai{}, nil, hand, sim.Matchup{IncomingDamage: 4}, nil, nil, priorAuras)
	if got.Value != 9 {
		t.Fatalf("Value = %d, want 9 (Reduce defense 4 + Yinti Yanti 4 with +1 aura bonus + creation credit 1)", got.Value)
	}
}

// Peace of Mind defends, creating a Ponder; Yinti Yanti's Play then sees the aura and
// credits +1{p}. Blue pitch funds Peace of Mind's cost.
func TestYintiYanti_SeesPonderFromPeaceOfMindInDefense(t *testing.T) {
	hand := []sim.Card{cards.YintiYantiRed{}, cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	got := sim.Best(heroes.Viserai{}, nil, hand, sim.Matchup{IncomingDamage: 4}, nil, nil)
	if got.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Peace of Mind defense 4 + Yinti Yanti 4 with +1 aura bonus from Ponder)", got.Value)
	}
}
