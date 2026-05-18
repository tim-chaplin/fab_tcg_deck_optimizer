package turntests

// End-to-end tests for Yinti Yanti's "while you control an aura, +1{p}" rider seeing
// auras created by defenders earlier in the same turn.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Reduce to Runechant defends, creating a Runechant; Yinti Yanti's Play then sees the
// aura and credits +1{p}. Red pitch funds Reduce's cost (no carryover Runechant — it'd
// pollute the test by satisfying Yinti Yanti's gate before Reduce ever played).
func TestYintiYanti_SeesRunechantFromReduceInDefense(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{cards.YintiYantiRed{}, cards.ReduceToRunechantRed{}, testutils.RedPitch{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), hand)
	got := summary.Value
	if got != 9 {
		t.Fatalf("Value = %d, want 9 (Reduce defense 4 + Yinti Yanti 4 with +1 aura bonus + creation credit 1)", got)
	}
}

// Peace of Mind defends, creating a Ponder; Yinti Yanti's Play then sees the aura and
// credits +1{p}. Blue pitch funds Peace of Mind's cost.
func TestYintiYanti_SeesPonderFromPeaceOfMindInDefense(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{cards.YintiYantiRed{}, cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), hand)
	got := summary.Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (Peace of Mind defense 4 + Yinti Yanti 4 with +1 aura bonus from Ponder)", got)
	}
}

// Tests that Yinti Yanti Blue plain-blocking alongside Reduce as DR sees Reduce's
// Runechant when its Block runs (DRs run first, populating auras before plain-block).
func TestYintiYanti_BlueBlockSeesRunechantFromReduceInDefense(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{cards.YintiYantiBlue{}, cards.ReduceToRunechantRed{}, testutils.RedPitch{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(7).Build(), hand)
	got := summary.Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (Reduce defense 4 + Yinti Yanti block 3 with +1 aura bonus + creation credit 1)", got)
	}
}

// Tests that Yinti Yanti Blue plain-blocking alongside Peace of Mind sees the Ponder PoM's
// DR puts in auras (DR runs first, populating auras before the plain-block hook).
func TestYintiYanti_BlueBlockSeesPonderFromPeaceOfMindInDefense(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{cards.YintiYantiBlue{}, cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(7).Build(), hand)
	got := summary.Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (Peace of Mind defense 4 + Yinti Yanti block 3 with +1 aura bonus from Ponder)", got)
	}
}
