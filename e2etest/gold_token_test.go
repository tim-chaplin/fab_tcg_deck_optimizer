package e2etest

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
	got := sim.BestWithTriggers(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, hand, sim.Matchup{IncomingDamage: 0}, deck, nil, nil, priorItems)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.State.Gold() != 0 {
		t.Fatalf("Gold after turn = %d, want 0 (the only token spent)", got.State.Gold())
	}
	if got.State.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Gold ability draws one card)", got.State.CardsDrawn)
	}
	if got.State.Arsenal == nil {
		t.Fatalf("Arsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.State.Deck) < d.Hero.Intelligence() {
		t.Fatalf("State.Deck has %d cards, want >= %d so next turn deals a full hand",
			len(got.State.Deck), d.Hero.Intelligence())
	}
	_ = d
}
