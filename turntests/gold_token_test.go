package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

func TestGoldAbility_SpendsToFillArsenalAndSwings(t *testing.T) {
	deck := []sim.Card{
		// Five fillers covers the gold-spend draw plus next-turn's 4 dealt cards.
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, deck)
	hand := []sim.Card{testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewGoldItem(1)}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.TurnState{Items: priorItems}, hand)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.Gold() != 0 {
		t.Fatalf("Gold after turn = %d, want 0 (the only token spent)", got.Gold())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Gold ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.StartOfNextTurnHand) != d.Hero.(sim.Hero).Intelligence() {
		t.Fatalf("StartOfNextTurnHand size = %d, want %d (Gold-spend draw should leave enough deck for next turn's full deal)",
			len(got.StartOfNextTurnHand), d.Hero.(sim.Hero).Intelligence())
	}
}
