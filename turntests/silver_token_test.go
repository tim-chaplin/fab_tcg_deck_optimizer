package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests the Silver token activated ability end-to-end: spend Silver alone, verify the
// drawn card promotes into arsenal.
func TestSilverToken_SpendsToFillArsenal(t *testing.T) {
	cards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai, nil, cards)
	hand := []card.Card{testutils.BluePitch{}}
	priorItems := []*item.Item{token.NewSilver(1)}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems(priorItems...), hand)
	if summary.Value != 0 {
		t.Fatalf("Value = %d, want 0 (Silver ability has no damage)", summary.Value)
	}
	if summary.State.SilverCount() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", summary.State.SilverCount())
	}
	if summary.State.Arsenal() == nil {
		t.Fatalf("Arsenal() = nil, want the Silver-drawn card promoted into the slot")
	}
	if len(summary.State.Hand()) != d.Hero.(hero.Hero).Intelligence() {
		t.Fatalf("Hand() size = %d, want %d (Silver-spend draw should leave enough deck for next turn's full deal)",
			len(summary.State.Hand()), d.Hero.(hero.Hero).Intelligence())
	}
}

// Tests that the Silver ability composes with a weapon swing when the pitch budget covers
// both costs.
func TestSilverToken_SpendsAndSwings(t *testing.T) {
	cards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, cards)
	hand := []card.Card{testutils.BluePitch{}, testutils.BluePitch{}}
	priorItems := []*item.Item{token.NewSilver(1)}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems(priorItems...), hand)
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", summary.Value)
	}
	if summary.State.SilverCount() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", summary.State.SilverCount())
	}
	if summary.State.Arsenal() == nil {
		t.Fatalf("Arsenal() = nil, want the Silver-drawn card promoted into the slot")
	}
}
