package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the hand cycle pops the first hand card to deck bottom and draws one.
func TestScourTheBattlescape_CyclesHandCardToDeckBottomAndDraws(t *testing.T) {
	for _, c := range []card.Card{cards.ScourTheBattlescapeRed{}, cards.ScourTheBattlescapeYellow{}, cards.ScourTheBattlescapeBlue{}} {
		cycled := testutils.GenericAttack(0, 1).WithName("cycled")
		drawn := testutils.GenericAttack(0, 2).WithName("drawn")
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{drawn}).Build()}
		ge.SetHand([]card.Card{cycled})
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)

		hand := ge.Hand()
		if len(hand) != 1 || hand[0].DisplayName() != "drawn" {
			t.Errorf("%s: Hand = %v, want [drawn] (deck top drawn into hand)", c.Name(), hand)
		}
		if d := ge.Deck(); d.Size() != 1 || d.PeekTop().(card.Card).DisplayName() != "cycled" {
			t.Errorf("%s: Deck = %v, want [cycled] (the cycled card sits alone at the bottom)", c.Name(), d)
		}
	}
}

// Tests that an empty hand suppresses the cycle: no card is drawn from the seeded deck
// and the deck stays untouched.
func TestScourTheBattlescape_EmptyHandSuppressesCycle(t *testing.T) {
	for _, c := range []card.Card{cards.ScourTheBattlescapeRed{}, cards.ScourTheBattlescapeYellow{}, cards.ScourTheBattlescapeBlue{}} {
		top := testutils.GenericAttack(0, 1).WithName("top")
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{top}).Build()}
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if len(ge.Hand()) != 0 {
			t.Errorf("%s: Hand = %v, want [] (no DrawOne should fire)", c.Name(), ge.Hand())
		}
		if d := ge.Deck(); d.Size() != 1 || d.PeekTop().(card.Card).DisplayName() != "top" {
			t.Errorf("%s: Deck = %v, want [top] (untouched)", c.Name(), d)
		}
	}
}
