package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that a card-sourced item trigger fires on a matching event and, when it does not
// self-destroy, stays in the arena.
func TestItemTrigger_FiresAndStays(t *testing.T) {
	ge := New()
	fired := 0
	ge.CreateItem(fakeCard{name: "Test Talisman"}, triggertype.Hit,
		func(_ card.GameEngine, _ card.Logger, _ card.Item, _ card.FireContext) { fired++ }, false, nil)

	ge.FireTriggers(card.FireContext{FiringType: triggertype.Hit, TriggeringCard: &card.CardState{Card: fakeCard{name: "attacker"}}})

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
	ge.CreateItem(fakeCard{name: "Test Talisman"}, triggertype.Hit,
		func(_ card.GameEngine, _ card.Logger, _ card.Item, _ card.FireContext) { fired++ }, false, nil)

	ge.FireTriggers(card.FireContext{FiringType: triggertype.EndOfTurn})

	if fired != 0 {
		t.Fatalf("handler fired %d times on a non-matching event, want 0", fired)
	}
}

// Tests that an item trigger whose handler calls Destroy is spliced out of the arena and
// its source card lands in the graveyard.
func TestItemTrigger_SelfDestructRemovesItemAndGraveyards(t *testing.T) {
	ge := New()
	src := fakeCard{name: "Test Talisman"}
	ge.CreateItem(src, triggertype.Hit,
		func(_ card.GameEngine, _ card.Logger, self card.Item, _ card.FireContext) {
			self.Destroy(true)
		}, false, nil)

	ge.FireTriggers(card.FireContext{FiringType: triggertype.Hit, TriggeringCard: &card.CardState{Card: fakeCard{name: "attacker"}}})

	if len(ge.Items()) != 0 {
		t.Fatalf("Items has %d entries after self-destruct, want 0", len(ge.Items()))
	}
	if g := ge.Graveyard(); len(g) != 1 || g[0] != card.Card(src) {
		t.Fatalf("Graveyard = %v, want one entry (the source card)", g)
	}
}

// Tests that DestroyItemObject removes exactly the item it's handed (not a same-source
// duplicate) and graveyards that item's source card.
func TestDestroyItemObject_RemovesTheGivenItem(t *testing.T) {
	ge := New()
	src := fakeCard{name: "Amulet"}
	ge.CreateItemWithAbility(src, fakeCard{name: "Amulet"})
	ge.CreateItemWithAbility(src, fakeCard{name: "Amulet"})
	items := ge.Items()
	if len(items) != 2 {
		t.Fatalf("Items has %d entries, want 2", len(items))
	}
	first, second := items[0], items[1]

	// Destroy the second object; a name-matched destroy would wrongly take the first.
	ge.DestroyItemObject(second.(card.Item), true)

	remaining := ge.Items()
	if len(remaining) != 1 {
		t.Fatalf("Items has %d entries after destroy, want 1", len(remaining))
	}
	if remaining[0] != first {
		t.Error("DestroyItemObject removed the wrong item — the first object should remain")
	}
	if g := ge.Graveyard(); len(g) != 1 || g[0] != card.Card(src) {
		t.Fatalf("Graveyard = %v, want one entry (the destroyed item's source)", g)
	}
}

// Tests that a token item carries no trigger, so FireTriggers never fires or destroys it.
func TestItemTrigger_TokenItemNeverFires(t *testing.T) {
	ge := New()
	ge.CreateGold(1)

	ge.FireTriggers(card.FireContext{FiringType: triggertype.Hit, TriggeringCard: &card.CardState{Card: fakeCard{name: "attacker"}}})
	ge.FireTriggers(card.FireContext{FiringType: triggertype.EndOfTurn})

	if ge.GoldCount() != 1 {
		t.Fatalf("Gold count = %d after FireTriggers, want 1 (token must not fire)", ge.GoldCount())
	}
}
