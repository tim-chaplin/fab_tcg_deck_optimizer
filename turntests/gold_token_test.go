package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

func TestGoldToken_SpendsToFillArsenalAndSwings(t *testing.T) {
	cards := []deck.Card{
		// Five fillers covers the gold-spend draw plus next-turn's 4 dealt cards.
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(hero.Viserai{}, []deck.Weapon{weapons.ReapingBlade{}}, cards)
	hand := []deck.Card{testutils.BluePitch{}}
	priorItems := []*token.Item{token.NewGold(1)}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, stateWithItems(priorItems...), hand)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.GoldCount() != 0 {
		t.Fatalf("Gold after turn = %d, want 0 (the only token spent)", got.GoldCount())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Gold ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.StartOfNextTurnHand) != d.Hero.(hero.Hero).Intelligence() {
		t.Fatalf("StartOfNextTurnHand size = %d, want %d (Gold-spend draw should leave enough deck for next turn's full deal)",
			len(got.StartOfNextTurnHand), d.Hero.(hero.Hero).Intelligence())
	}
}
