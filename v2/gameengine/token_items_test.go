package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// deckOf builds a *deck.Deck whose draw order is the supplied cards. Used to seed a
// non-empty deck before token-ability DrawOne calls.
func deckOf(cs ...card.Card) *deck.Deck {
	dc := make([]deck.Card, len(cs))
	for i, c := range cs {
		dc[i] = c
	}
	return deck.New(nil, nil, dc)
}

// Tests that cards.GoldToken.Play decrements Count and removes the entry at zero.
func TestGoldToken_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := New()
	ge.SetDeck(deckOf(stubCard{name: "filler"}))
	ge.CreateItem(token.NewGold(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.GoldToken{}})
	if ge.GoldCount() != 0 {
		t.Fatalf("Gold = %d after spending the only token, want 0", ge.GoldCount())
	}
	if len(ge.Items()) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(ge.Items()))
	}
	if h := ge.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests that spending one of multiple Gold tokens leaves the entry at decremented Count.
func TestGoldToken_PlayDecrementsCountWhenMultiple(t *testing.T) {
	ge := New()
	ge.SetDeck(deckOf(stubCard{name: "filler"}))
	ge.CreateItem(token.NewGold(3))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.GoldToken{}})
	if ge.GoldCount() != 2 {
		t.Fatalf("Gold = %d after spending 1 of 3, want 2", ge.GoldCount())
	}
}

// Tests cards.SilverToken.Play decrement + draw behaviour.
func TestSilverToken_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := New()
	ge.SetDeck(deckOf(stubCard{name: "filler"}))
	ge.CreateItem(token.NewSilver(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.SilverToken{}})
	if ge.SilverCount() != 0 {
		t.Fatalf("Silver = %d after spending the only token, want 0", ge.SilverCount())
	}
	if len(ge.Items()) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(ge.Items()))
	}
	if h := ge.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests cards.CopperToken.Play decrement + draw behaviour.
func TestCopperToken_PlaysDecrementsAndDestroys(t *testing.T) {
	ge := New()
	ge.SetDeck(deckOf(stubCard{name: "filler"}))
	ge.CreateItem(token.NewCopper(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.CopperToken{}})
	if ge.CopperCount() != 0 {
		t.Fatalf("Copper = %d after spending the only token, want 0", ge.CopperCount())
	}
	if len(ge.Items()) != 0 {
		t.Fatalf("Items still has %d entries after destroy, want 0", len(ge.Items()))
	}
	if h := ge.Hand(); len(h) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one card)", len(h))
	}
}

// Tests CreateSilver and CreateCopper consolidate by token type.
func TestCreateSilverCopper_BumpsExistingEntry(t *testing.T) {
	ge := New()
	ge.CreateSilver(2)
	ge.CreateSilver(1)
	ge.CreateCopper(1)
	if ge.SilverCount() != 3 {
		t.Errorf("Silver = %d, want 3 (2 + 1 consolidated)", ge.SilverCount())
	}
	if ge.CopperCount() != 1 {
		t.Errorf("Copper = %d, want 1", ge.CopperCount())
	}
	if got := len(ge.Items()); got != 2 {
		t.Errorf("Items entries = %d, want 2 (one Silver + one Copper)", got)
	}
}
