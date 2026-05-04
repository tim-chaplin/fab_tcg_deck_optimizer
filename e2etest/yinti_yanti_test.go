package e2etest

// End-to-end tests for Yinti Yanti's "while you control an aura, +1{p}" rider in the
// FaB-correct turn order: defense resolves first, so a Runechant or Ponder created by a
// defender on the same turn is in play when Yinti Yanti attacks during the action phase.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests Yinti Yanti's attack-side aura bonus when Reduce to Runechant defends earlier
// in the same turn. Reduce creates a Runechant when it resolves; in the new turn order
// that Runechant is visible to chain cards, so Yinti Yanti's Play sees len(s.Auras) > 0
// and credits the +1{p} buff. Reduce gets one runechant carryover to fund its cost.
func TestYintiYanti_SeesRunechantFromReduceInDefense(t *testing.T) {
	hand := []sim.Card{cards.YintiYantiRed{}, cards.ReduceToRunechantRed{}}
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	priorAuras := []sim.Aura{sim.NewRunechantAura(1)}
	got := sim.BestWithTriggers(heroes.Viserai{}, nil, hand, sim.Matchup{IncomingDamage: 4}, nil, nil, priorAuras)
	_ = d
	// Expected Value:
	//   defense: Reduce prevents 4 (cost paid by the carryover Runechant),
	//   action:  Yinti Yanti Red attacks for printed 3 + aura bonus 1 = 4,
	//   creation credit: Reduce's Play creates one Runechant (+1 damage-equivalent at
	//                    creation time per the optimistic-credit rule).
	// Total = 4 + 4 + 1 = 9.
	if got.Value != 9 {
		t.Fatalf("Value = %d, want 9 (Reduce defense 4 + Yinti Yanti 4 with +1 aura bonus + creation credit 1)", got.Value)
	}
}

// Tests Yinti Yanti's attack-side aura bonus when Peace of Mind defends earlier in the
// same turn. Peace of Mind creates a Ponder; the Ponder is in play during action so
// Yinti Yanti's Play sees an aura and credits +1{p}.
func TestYintiYanti_SeesPonderFromPeaceOfMindInDefense(t *testing.T) {
	hand := []sim.Card{cards.YintiYantiRed{}, cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	got := sim.Best(heroes.Viserai{}, nil, hand, sim.Matchup{IncomingDamage: 4}, nil, nil)
	// Expected Value:
	//   defense: Blue pitch funds Peace of Mind cost 2; Peace of Mind prevents 4,
	//   action:  Yinti Yanti Red attacks for printed 3 + aura bonus 1 = 4.
	// Total = 4 + 4 = 8.
	if got.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Peace of Mind defense 4 + Yinti Yanti 4 with +1 aura bonus from Ponder)", got.Value)
	}
}
