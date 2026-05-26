package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Battlefront Bastion's +1 alone-bonus fires when it's the only plain blocker.
func TestBlock_BattlefrontBastionAloneFiresBonus(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.BattlefrontBastionRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(5).Build(), hand)
	got := summary.Value
	if got != 3 {
		t.Fatalf("Value = %d, want 3 (BB block 2 + alone bonus 1)", got)
	}
}

// Tests that a DR blocking alongside Battlefront Bastion doesn't cancel its alone-bonus —
// only a second simultaneous plain block does.
func TestBlock_BattlefrontBastionAloneFiresBesideDR(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.BattlefrontBastionRed{},
		cards.ToughenUpBlue{},
		testutils.FakeBlueResource(),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(10).Build(), hand)
	got := summary.Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (Toughen Up DR 4 + BB 2 + alone bonus 1)", got)
	}
}

// Tests that another Block-typed card (On the Horizon) alongside Battlefront Bastion
// counts as a second plain blocker — so BB's alone-bonus does NOT fire.
func TestBlock_BattlefrontBastionAloneCancelledByBlockCard(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.BattlefrontBastionRed{}, cards.OnTheHorizonRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(10).Build(), hand)
	got := summary.Value
	if got != 6 {
		t.Fatalf("Value = %d, want 6 (OTH 4 + BB 2; no alone bonus, OTH is a plain blocker)", got)
	}
}

// Tests that a second plain blocker cancels Battlefront Bastion's alone-bonus while
// firing Right Behind You's together-bonus on the same attack turn.
func TestBlock_BattlefrontBastionAloneCancelledByPlainBlocker(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.BattlefrontBastionRed{}, cards.RightBehindYouRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(5).Build(), hand)
	got := summary.Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (BB 2 + RBY 2 + RBY together bonus 1)", got)
	}
}

// Tests that Right Behind You's together-bonus stays silent when no other plain blocker
// shares the defenders slot.
func TestBlock_RightBehindYouAloneNoBonus(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.RightBehindYouRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(5).Build(), hand)
	got := summary.Value
	if got != 2 {
		t.Fatalf("Value = %d, want 2 (RBY block 2; no together bonus)", got)
	}
}
