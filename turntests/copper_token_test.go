package turntests

import (
	"testing"

	cardpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// Tests that the Copper token ability stays unspent when the chain can't fund its {4} cost.
func TestCopperAbility_NotEnoughResourceSkipsSpend(t *testing.T) {
	cards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, cards)
	hand := []deck.Card{testutils.BluePitch{}}
	priorItems := []*token.Item{cardpkg.NewCopper(1)}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.Prior{Items: priorItems}, hand)
	if got.CopperCount() != 1 {
		t.Fatalf("Copper after turn = %d, want 1 (single blue pitch can't fund {4})", got.CopperCount())
	}
}

// Tests the Copper ability composes with a weapon swing when the pitch budget covers
// both. Two blue pitches (3+3=6 res) fund the Copper ability ({4}) plus a Reaping
// Blade swing ({1}), with 1 res to spare.
func TestCopperAbility_SpendsAndSwings(t *testing.T) {
	cards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.ReapingBlade{}}, cards)
	hand := []deck.Card{testutils.BluePitch{}, testutils.BluePitch{}}
	priorItems := []*token.Item{cardpkg.NewCopper(1)}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.Prior{Items: priorItems}, hand)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.CopperCount() != 0 {
		t.Fatalf("Copper after turn = %d, want 0 (the only token spent)", got.CopperCount())
	}
	if got.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1 (Copper ability draws one card)", got.CardsDrawn)
	}
	if got.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted into the slot")
	}
}
