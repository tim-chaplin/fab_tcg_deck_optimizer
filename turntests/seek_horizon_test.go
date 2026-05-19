package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that an empty hand suppresses the go-again grant.
func TestSeekHorizon_EmptyHandNoGoAgain(t *testing.T) {
	for _, c := range []card.Card{cards.SeekHorizonRed{}, cards.SeekHorizonYellow{}, cards.SeekHorizonBlue{}} {
		ge := gameengine.New()
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = true with empty hand, want false", c.Name(), c.Pitch())
		}
	}
}

// Tests that paying the alt cost grants go-again and prepends the spent hand card to the deck.
func TestSeekHorizon_PaysAltCostAndGrantsGoAgain(t *testing.T) {
	for _, c := range []card.Card{cards.SeekHorizonRed{}, cards.SeekHorizonYellow{}, cards.SeekHorizonBlue{}} {
		ge := gameengine.New()
		spare := testutils.GenericAttack(0, 1)
		ge.GameState.AppendHandRaw(spare)
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)

		if !self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after paying alt cost, want true", c.Name(), c.Pitch())
		}
		if len(ge.Hand()) != 0 {
			t.Errorf("%s [%d{p}]: hand size = %d after alt cost, want 0 (card moved to deck top)", c.Name(), c.Pitch(), len(ge.Hand()))
		}
		top, ok := ge.PeekDeck()
		if !ok || top != spare {
			t.Errorf("%s [%d{p}]: deck top = %v, want %v (the popped hand card)", c.Name(), c.Pitch(), top, spare)
		}
	}
}
