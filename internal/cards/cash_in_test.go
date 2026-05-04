package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Cash In draws two cards from the deck top into hand.
func TestCashIn_DrawsTwo(t *testing.T) {
	deck := []sim.Card{testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	s := sim.NewTurnState(deck, nil)
	(CashInYellow{}).Play(s, &sim.CardState{Card: CashInYellow{}})
	if len(s.Hand) != 2 {
		t.Fatalf("Hand size = %d, want 2 (drew two cards)", len(s.Hand))
	}
	if s.CardsDrawn != 2 {
		t.Fatalf("CardsDrawn = %d, want 2", s.CardsDrawn)
	}
}

// Tests that Cash In gracefully handles a deck with fewer than two cards.
func TestCashIn_DrawsAvailableOnShortDeck(t *testing.T) {
	deck := []sim.Card{testutils.RedAttack{}}
	s := sim.NewTurnState(deck, nil)
	(CashInYellow{}).Play(s, &sim.CardState{Card: CashInYellow{}})
	if len(s.Hand) != 1 {
		t.Fatalf("Hand size = %d, want 1 (drew one of two)", len(s.Hand))
	}
	if s.CardsDrawn != 1 {
		t.Fatalf("CardsDrawn = %d, want 1", s.CardsDrawn)
	}
}
