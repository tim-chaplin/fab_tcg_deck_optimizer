package turntests

// Token basics — one focused lifecycle test per built-in token kind. Each test plays a card
// that creates a single token, exercises the token's defining behaviour (trigger fire or
// activated ability), and asserts the resulting count + destroy semantics. These guard the
// shared token-slot machinery in internal/gameengine/tokens.go.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that a Runechant created on turn 1 fires on turn 2's attack and self-destroys.
func TestTokenBasics_Runechant_CreateFireDestroy(t *testing.T) {
	// Top-of-deck cards become turn 2's dealt hand. Sized to the hero's intelligence (4)
	// so the hand is fully populated with attacks.
	deckCards := []deck.Card{
		testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack(),
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	turn1Hand := []card.Card{cards.HocusPocusRed{}, testutils.FakeBlueResource(), testutils.FakeBlueResource()}
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), turn1Hand)
	if got := turn1.State.RunechantCount(); got != 1 {
		t.Errorf("turn 1 RunechantCount = %d, want 1 (Hocus Pocus rider created one)", got)
	}
	if got := turn2.State.RunechantCount(); got != 0 {
		t.Errorf("turn 2 RunechantCount = %d, want 0 (carryover Runechant fired on the attack action)", got)
	}
}

// Tests that a Ponder created by Peace of Mind fires at end-of-turn and self-destroys.
func TestTokenBasics_Ponder_CreateFireDestroy(t *testing.T) {
	beacon := testutils.FakeRedAttack()
	deckCards := []deck.Card{beacon, testutils.FakeBlueAttack(), testutils.FakeBlueAttack(),
		testutils.FakeBlueAttack(), testutils.FakeBlueAttack(), testutils.FakeBlueAttack(), testutils.FakeBlueAttack()}
	d := deck.New(heroes.Viserai, nil, deckCards)
	hand := []card.Card{cards.PeaceOfMindRed{}, testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), hand)
	if got := summary.State.PonderCount(); got != 0 {
		t.Errorf("PonderCount after turn = %d, want 0 (end-of-turn fire destroys the Ponder)", got)
	}
	if summary.State.Arsenal() == nil || summary.State.Arsenal().Name() != beacon.Name() {
		t.Errorf("Arsenal = %v, want %v (Ponder's deck pop should fill the empty arsenal)", summary.State.Arsenal(), beacon)
	}
}

// Tests that a carryover Gold token spends via its activated ability and zeroes count.
func TestTokenBasics_Gold_CreateSpendDestroy(t *testing.T) {
	deckCards := []deck.Card{testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack()}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
	hand := []card.Card{testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems([]*item.Item{token.NewGold(1)}...), hand)
	if got := summary.State.GoldCount(); got != 0 {
		t.Errorf("GoldCount after turn = %d, want 0 (Spend ability consumed the token)", got)
	}
	if summary.State.Arsenal() == nil {
		t.Errorf("Arsenal = nil, want the Gold-drawn card promoted into the empty slot")
	}
}

// Tests that a carryover Silver token spends via its activated ability and zeroes count.
func TestTokenBasics_Silver_CreateSpendDestroy(t *testing.T) {
	deckCards := []deck.Card{testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack()}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
	hand := []card.Card{testutils.FakeBlueResource(), testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems([]*item.Item{token.NewSilver(1)}...), hand)
	if got := summary.State.SilverCount(); got != 0 {
		t.Errorf("SilverCount after turn = %d, want 0 (Spend ability consumed the token)", got)
	}
	if summary.State.Arsenal() == nil {
		t.Errorf("Arsenal = nil, want the Silver-drawn card promoted into the empty slot")
	}
}

// Tests that a carryover Copper token spends via its activated ability and zeroes count.
func TestTokenBasics_Copper_CreateSpendDestroy(t *testing.T) {
	deckCards := []deck.Card{testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack()}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, deckCards)
	hand := []card.Card{testutils.FakeBlueResource(), testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, stateWithItems([]*item.Item{token.NewCopper(1)}...), hand)
	if got := summary.State.CopperCount(); got != 0 {
		t.Errorf("CopperCount after turn = %d, want 0 (Spend ability consumed the token)", got)
	}
}
