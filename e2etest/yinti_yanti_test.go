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
// aura and credits +1{p}. Red pitch funds Reduce's cost (no carryover Runechant — it'd
// pollute the test by satisfying Yinti Yanti's gate before Reduce ever played).
func TestYintiYanti_SeesRunechantFromReduceInDefense(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.YintiYantiRed{}, cards.ReduceToRunechantRed{}, testutils.RedPitch{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 4}, sim.TurnState{}, hand).Value
	if got != 9 {
		t.Fatalf("Value = %d, want 9 (Reduce defense 4 + Yinti Yanti 4 with +1 aura bonus + creation credit 1)", got)
	}
}

// Peace of Mind defends, creating a Ponder; Yinti Yanti's Play then sees the aura and
// credits +1{p}. Blue pitch funds Peace of Mind's cost.
func TestYintiYanti_SeesPonderFromPeaceOfMindInDefense(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.YintiYantiRed{}, cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 4}, sim.TurnState{}, hand).Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (Peace of Mind defense 4 + Yinti Yanti 4 with +1 aura bonus from Ponder)", got)
	}
}

// Tests that Yinti Yanti Blue plain-blocking alongside Reduce as DR sees Reduce's
// Runechant when its Block runs (DRs run first, populating s.Auras).
func TestYintiYanti_BlueBlockSeesRunechantFromReduceInDefense(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.YintiYantiBlue{}, cards.ReduceToRunechantRed{}, testutils.RedPitch{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 7}, sim.TurnState{}, hand).Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (Reduce defense 4 + Yinti Yanti block 3 with +1 aura bonus + creation credit 1)", got)
	}
}

// Defense-side mirror of TestYintiYanti_SeesPonderFromPeaceOfMindInDefense: Yinti
// Yanti Blue plain-blocks alongside Peace of Mind. PoM runs in the DR loop first, so
// the Ponder is in s.Auras when Yinti's Block runs, gating the +1{d}. Incoming bumped
// to 7 so blocking with Yinti is the optimal partition.
func TestYintiYanti_BlueBlockSeesPonderFromPeaceOfMindInDefense(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.YintiYantiBlue{}, cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 7}, sim.TurnState{}, hand).Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (Peace of Mind defense 4 + Yinti Yanti block 3 with +1 aura bonus from Ponder)", got)
	}
}
