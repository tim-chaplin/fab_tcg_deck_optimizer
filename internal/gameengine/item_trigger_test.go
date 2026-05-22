package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that a card-sourced item trigger fires on a matching event and, when it does not
// self-destroy, stays in the arena.
func TestItemTrigger_FiresAndStays(t *testing.T) {
	ge := New()
	fired := 0
	ge.CreateItem(item.NewFromCard(stubCard{name: "Test Talisman"}, triggertype.Hit,
		func(_ card.GameEngine, _ card.Logger, _ card.Item) { fired++ }, false))

	ge.FireTriggers(triggertype.Hit, stubCard{name: "attacker"})

	if fired != 1 {
		t.Fatalf("handler fired %d times, want 1", fired)
	}
	if len(ge.Items()) != 1 {
		t.Fatalf("Items has %d entries, want 1 (item should persist)", len(ge.Items()))
	}
}

// Tests that an item trigger does not run for an event type it doesn't subscribe to.
func TestItemTrigger_SkipsNonMatchingEvent(t *testing.T) {
	ge := New()
	fired := 0
	ge.CreateItem(item.NewFromCard(stubCard{name: "Test Talisman"}, triggertype.Hit,
		func(_ card.GameEngine, _ card.Logger, _ card.Item) { fired++ }, false))

	ge.FireTriggers(triggertype.EndOfTurn, nil)

	if fired != 0 {
		t.Fatalf("handler fired %d times on a non-matching event, want 0", fired)
	}
}

// Tests that an item trigger whose handler calls Destroy is spliced out of the arena and
// its source card lands in the graveyard.
func TestItemTrigger_SelfDestructRemovesItemAndGraveyards(t *testing.T) {
	ge := New()
	src := stubCard{name: "Test Talisman"}
	ge.CreateItem(item.NewFromCard(src, triggertype.Hit,
		func(_ card.GameEngine, _ card.Logger, self card.Item) { self.Destroy(true) }, false))

	ge.FireTriggers(triggertype.Hit, stubCard{name: "attacker"})

	if len(ge.Items()) != 0 {
		t.Fatalf("Items has %d entries after self-destruct, want 0", len(ge.Items()))
	}
	if g := ge.Graveyard(); len(g) != 1 || g[0] != card.Card(src) {
		t.Fatalf("Graveyard = %v, want one entry (the source card)", g)
	}
}

// Tests that a token item carries no trigger, so FireTriggers never fires or destroys it.
func TestItemTrigger_TokenItemNeverFires(t *testing.T) {
	ge := New()
	ge.CreateItem(token.NewGold(1))

	ge.FireTriggers(triggertype.Hit, stubCard{name: "attacker"})
	ge.FireTriggers(triggertype.EndOfTurn, nil)

	if ge.GoldCount() != 1 {
		t.Fatalf("Gold count = %d after FireTriggers, want 1 (token must not fire)", ge.GoldCount())
	}
}
