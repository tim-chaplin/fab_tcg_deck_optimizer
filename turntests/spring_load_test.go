package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Spring Load with a non-empty hand attacks for printed power only.
func TestSpringLoad_BasePower(t *testing.T) {
	for _, c := range []card.Card{cards.SpringLoadRed{}, cards.SpringLoadYellow{}, cards.SpringLoadBlue{}} {
		ge := gameengine.New()
		ge.SetHand([]card.Card{testutils.GenericAttack(0, 0)})
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 2 {
			t.Errorf("%s: Play() with non-empty hand = %d, want 2", c.Name(), got)
		}
	}
}

// Tests that Spring Load with an empty hand fires the +3{p} rider on every variant.
func TestSpringLoad_EmptyHandFiresRider(t *testing.T) {
	for _, c := range []card.Card{cards.SpringLoadRed{}, cards.SpringLoadYellow{}, cards.SpringLoadBlue{}} {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 5 {
			t.Errorf("%s: Play() with empty hand = %d, want 5 (2 printed + 3 rider)", c.Name(), got)
		}
	}
}
