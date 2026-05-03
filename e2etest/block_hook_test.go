package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Battlefront Bastion's +1 alone-bonus fires when it's the only plain blocker.
func TestBlock_BattlefrontBastionAloneFiresBonus(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.BattlefrontBastionRed{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value
	if got != 3 {
		t.Fatalf("Value = %d, want 3 (BB block 2 + alone bonus 1)", got)
	}
}

// Tests that a Defense Reaction (On the Horizon) blocking alongside doesn't cancel
// Battlefront Bastion's alone-bonus — only a simultaneous additional plain block does.
func TestBlock_BattlefrontBastionAloneFiresBesideDR(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.BattlefrontBastionRed{}, cards.OnTheHorizonRed{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 10}, nil, hand).Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (OTH 4 + BB 2 + alone bonus 1)", got)
	}
}

// Tests that a second plain blocker cancels Battlefront Bastion's alone-bonus while
// firing Right Behind You's together-bonus on the same chain.
func TestBlock_BattlefrontBastionAloneCancelledByPlainBlocker(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.BattlefrontBastionRed{}, cards.RightBehindYouRed{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (BB 2 + RBY 2 + RBY together bonus 1)", got)
	}
}

// Tests that Right Behind You's together-bonus stays silent when no other plain blocker
// shares the defenders slot.
func TestBlock_RightBehindYouAloneNoBonus(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.RightBehindYouRed{}}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value
	if got != 2 {
		t.Fatalf("Value = %d, want 2 (RBY block 2; no together bonus)", got)
	}
}
