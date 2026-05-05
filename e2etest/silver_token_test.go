package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests the Silver token activated ability end-to-end via the chain runner. Same
// shape as TestGoldAbility_SpendsToFillArsenalAndSwings: blue pitch (3) funds the
// Silver ability ({3}) — leaves 0 for a Reaping Blade swing, so we'd need pitch=4
// to cover both. Here we test the simpler case: spend Silver alone, verify the
// drawn card promotes into arsenal.
func TestSilverAbility_SpendsToFillArsenal(t *testing.T) {
	deck := []sim.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deck)
	hand := []sim.Card{testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewSilverItem(1)}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	if got.Value != 0 {
		t.Fatalf("Value = %d, want 0 (Silver ability has no damage)", got.Value)
	}
	if got.Silver() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", got.Silver())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Silver ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.StartOfNextTurnHand) != d.Hero.Intelligence() {
		t.Fatalf("StartOfNextTurnHand size = %d, want %d (Silver-spend draw should leave enough deck for next turn's full deal)",
			len(got.StartOfNextTurnHand), d.Hero.Intelligence())
	}
}

// Tests the Silver ability composes with a weapon swing when the pitch budget covers
// both. Two blue pitches (3+3=6 res) fund the Silver ability ({3}) plus a Reaping
// Blade swing ({1}); chain order doesn't matter since Silver and the swing both
// have go again.
func TestSilverAbility_SpendsAndSwings(t *testing.T) {
	deck := []sim.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, deck)
	hand := []sim.Card{testutils.BluePitch{}, testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewSilverItem(1)}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.Silver() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", got.Silver())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Silver ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
}
