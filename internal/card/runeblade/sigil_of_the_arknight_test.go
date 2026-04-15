package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestSigilOfTheArknight_EmptyDeckReturnsZero(t *testing.T) {
	// With no deck to reveal from, the expected value collapses to 0. AuraCreated still flips.
	var s card.TurnState
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() with empty deck = %d, want 0", got)
	}
	if !s.AuraCreated {
		t.Errorf("Play() did not set AuraCreated")
	}
}

func TestSigilOfTheArknight_ExpectedValueFromDeck(t *testing.T) {
	// Deck of 4: 2 attack actions, 2 non-attack. EV = (2*3)/4 = 1 (integer truncation).
	deck := []card.Card{
		stubRunebladeAttack{},
		stubRunebladeAttack{},
		stubNonAttack{},
		stubAura{},
	}
	s := card.TurnState{Deck: deck}
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 1 {
		t.Errorf("Play() = %d, want 1 (2 attack actions of 4 cards × 3 / 4)", got)
	}
}

func TestSigilOfTheArknight_AllAttackActionsDeck(t *testing.T) {
	// Every card in the deck is an attack action → EV = 1.0 × 3 = 3.
	deck := []card.Card{stubRunebladeAttack{}, stubRunebladeAttack{}, stubRunebladeAttack{}}
	s := card.TurnState{Deck: deck}
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
}

func TestSigilOfTheArknight_NoMemoMarker(t *testing.T) {
	// Sigil's Play depends on deck composition, so it must opt out of the hand-evaluation memo.
	var c card.Card = SigilOfTheArknightBlue{}
	if _, ok := c.(card.NoMemo); !ok {
		t.Errorf("SigilOfTheArknightBlue should implement card.NoMemo")
	}
}
