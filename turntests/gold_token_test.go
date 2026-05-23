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

func TestGoldToken_SpendsToFillArsenalAndSwings(t *testing.T) {
	cards := []deck.Card{
		// Five fillers covers the gold-spend draw plus next-turn's 4 dealt cards.
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, cards)
	hand := []card.Card{testutils.BluePitch{}}
	priorItems := []*item.Item{token.NewGold(1)}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems(priorItems...), hand)
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", summary.Value)
	}
	if summary.State.GoldCount() != 0 {
		t.Fatalf("Gold after turn = %d, want 0 (the only token spent)", summary.State.GoldCount())
	}
	if summary.State.Arsenal() == nil {
		t.Fatalf("Arsenal() = nil, want the Gold-drawn card promoted into the slot")
	}
	if len(summary.State.Hand()) != d.Hero.(hero.Hero).Intelligence() {
		t.Fatalf("Hand() size = %d, want %d (Gold-spend draw should leave enough deck for next turn's full deal)",
			len(summary.State.Hand()), d.Hero.(hero.Hero).Intelligence())
	}
}
