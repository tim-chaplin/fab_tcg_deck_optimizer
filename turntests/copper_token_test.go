package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// stateWithItems returns a *GameState seeded with the supplied items. Tests use this to
// pass token carryover (Gold / Silver / Copper) into EvalOneTurnForTesting.
func stateWithItems(items ...*item.Item) *gameengine.GameState {
	b := gameengine.GameStateBuilder()
	for _, it := range items {
		b.AddItem(it)
	}
	return b.Build()
}

// Tests that the Copper token ability stays unspent when the chain can't fund its {4} cost.
func TestCopperToken_NotEnoughResourceSkipsSpend(t *testing.T) {
	cards := []deck.Card{
		testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack(),
	}
	d := deck.New(heroes.Viserai, nil, cards)
	hand := []card.Card{testutils.FakeBlueResource()}
	priorItems := []*item.Item{token.NewCopper(1)}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems(priorItems...), hand)
	if summary.State.CopperCount() != 1 {
		t.Fatalf("Copper after turn = %d, want 1 (single blue pitch can't fund {4})", summary.State.CopperCount())
	}
}

// Tests the Copper ability composes with a weapon swing when the pitch budget covers
// both. Two blue pitches (3+3=6 res) fund the Copper ability ({4}) plus a Reaping
// Blade swing ({1}), with 1 res to spare.
func TestCopperToken_SpendsAndSwings(t *testing.T) {
	cards := []deck.Card{
		testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack(),
	}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, cards)
	hand := []card.Card{testutils.FakeBlueResource(), testutils.FakeBlueResource()}
	priorItems := []*item.Item{token.NewCopper(1)}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems(priorItems...), hand)
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", summary.Value)
	}
	if summary.State.CopperCount() != 0 {
		t.Fatalf("Copper after turn = %d, want 0 (the only token spent)", summary.State.CopperCount())
	}
	if summary.State.Arsenal() == nil {
		t.Fatalf("Arsenal() = nil, want the Copper-drawn card promoted into the slot")
	}
}
