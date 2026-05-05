package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests the Copper token activated ability end-to-end. {4} cost is the highest of
// the three token tiers — needs two blue pitches (3+3=6 res) to fund. With a single
// blue pitch (3 res), the chain is infeasible; the optimizer should leave the
// Copper alone and just hold the blue.
func TestCopperAbility_NotEnoughResourceSkipsSpend(t *testing.T) {
	deck := []sim.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deck)
	hand := []sim.Card{testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewCopperItem(1)}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	if got.Copper() != 1 {
		t.Fatalf("Copper after turn = %d, want 1 (single blue pitch can't fund {4})", got.Copper())
	}
}

// Tests the Copper ability composes with a weapon swing when the pitch budget covers
// both. Two blue pitches (3+3=6 res) fund the Copper ability ({4}) plus a Reaping
// Blade swing ({1}), with 1 res to spare.
func TestCopperAbility_SpendsAndSwings(t *testing.T) {
	deck := []sim.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, deck)
	hand := []sim.Card{testutils.BluePitch{}, testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewCopperItem(1)}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.Copper() != 0 {
		t.Fatalf("Copper after turn = %d, want 0 (the only token spent)", got.Copper())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Copper ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
}
