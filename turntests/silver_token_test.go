package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests the Silver token activated ability end-to-end: spend Silver alone, verify the
// drawn card promotes into arsenal.
func TestSilverAbility_SpendsToFillArsenal(t *testing.T) {
	cards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, cards)
	hand := []deck.Card{testutils.BluePitch{}}
	priorItems := []*sim.Item{sim.NewSilverItem(1)}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.Prior{Items: priorItems}, hand)
	if got.Value != 0 {
		t.Fatalf("Value = %d, want 0 (Silver ability has no damage)", got.Value)
	}
	if got.SilverCount() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", got.SilverCount())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Silver ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.StartOfNextTurnHand) != d.Hero.(gameengine.Hero).Intelligence() {
		t.Fatalf("StartOfNextTurnHand size = %d, want %d (Silver-spend draw should leave enough deck for next turn's full deal)",
			len(got.StartOfNextTurnHand), d.Hero.(gameengine.Hero).Intelligence())
	}
}

// Tests that the Silver ability composes with a weapon swing when the pitch budget covers
// both costs.
func TestSilverAbility_SpendsAndSwings(t *testing.T) {
	cards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.ReapingBlade{}}, cards)
	hand := []deck.Card{testutils.BluePitch{}, testutils.BluePitch{}}
	priorItems := []*sim.Item{sim.NewSilverItem(1)}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.Prior{Items: priorItems}, hand)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.SilverCount() != 0 {
		t.Fatalf("Silver after turn = %d, want 0 (the only token spent)", got.SilverCount())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Silver ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
}
