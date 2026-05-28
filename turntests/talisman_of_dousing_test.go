package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// dousingInItems reports whether any item in items carries Dousing's display name.
func dousingInItems(items []gameengine.Item) bool {
	for _, it := range items {
		if it.CardName() == (cards.TalismanOfDousingYellow{}).DisplayName() {
			return true
		}
	}
	return false
}

// Tests that Talisman of Dousing, in play at start of turn, prevents 1 arcane damage
// when the matchup deals arcane — crediting Value and consuming the talisman.
func TestTalismanOfDousing_PreventsOneArcaneAndDestroys(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateItemFromCard(cards.TalismanOfDousingYellow{}).
		SetIncomingArcaneDamage(2).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 1 {
		t.Errorf("Value = %d, want 1 (Dousing prevents 1 arcane)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if dousingInItems(summary.State.Items()) {
		t.Errorf("Talisman of Dousing still in items; should self-destruct after preventing")
	}
}

// Tests that on a physical-only damage turn Dousing's arcane prevention finds nothing to
// prevent — handler no-ops, talisman survives, no extra Value is credited.
func TestTalismanOfDousing_PhysicalOnlyTurnLeavesItemIntact(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateItemFromCard(cards.TalismanOfDousingYellow{}).
		SetIncomingPhysicalDamage(5).
		SetIncomingArcaneDamage(0).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 0 {
		t.Errorf("Value = %d, want 0 (no arcane damage, Dousing no-ops)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if !dousingInItems(summary.State.Items()) {
		t.Errorf("Talisman of Dousing missing from items; expected it to survive a physical-only turn\nFinal items: %v",
			summary.State.Items())
	}
}
