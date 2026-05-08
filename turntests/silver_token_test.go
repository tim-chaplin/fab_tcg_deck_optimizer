package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests the Silver token activated ability end-to-end: spend Silver alone, verify the
// drawn card promotes into arsenal.
func TestSilverAbility_SpendsToFillArsenal(t *testing.T) {
	deck := []sim.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deck)
	hand := []sim.Card{testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewSilverItem(1)}
	got := d.EvalTwoTurnsForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	if got.Turn1.Value != 0 {
		t.Fatalf("Value = %d, want 0 (Silver ability has no damage)", got.Turn1.Value)
	}
	if got.Turn1.State.Silver() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", got.Turn1.State.Silver())
	}
	if got.Turn1.State.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Silver ability draws one card)", got.Turn1.State.CardsDrawn)
	}
	if got.Turn1.State.Arsenal == nil {
		t.Fatalf("Arsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.Turn2.DealtHand) != d.Hero.Intelligence() {
		t.Fatalf("turn 2 dealt hand size = %d, want %d (Silver-spend draw should leave enough deck for next turn's full deal)",
			len(got.Turn2.DealtHand), d.Hero.Intelligence())
	}
}

// Tests that the Silver ability composes with a weapon swing when the pitch budget covers
// both costs.
func TestSilverAbility_SpendsAndSwings(t *testing.T) {
	deck := []sim.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, deck)
	hand := []sim.Card{testutils.BluePitch{}, testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewSilverItem(1)}
	summary := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	state := summary.State
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", summary.Value)
	}
	if state.Silver() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", state.Silver())
	}
	if state.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Silver ability draws one card)", state.CardsDrawn)
	}
	if state.Arsenal == nil {
		t.Fatalf("Arsenal = nil, want the drawn card promoted into the slot")
	}
}
