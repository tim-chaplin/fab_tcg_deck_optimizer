package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestTremorOfIArathael_NoBanishReturnsBaseAttack: empty Banish → printed power, no rider.
func TestTremorOfIArathael_NoBanishReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		base int
	}{
		{cards.TremorOfIArathaelRed{}, 4},
		{cards.TremorOfIArathaelYellow{}, 3},
		{cards.TremorOfIArathaelBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if ge.Value() != tc.base {
			t.Errorf("%s: Value = %d, want %d (no banish, base attack only)", tc.c.Name(), ge.Value(), tc.base)
		}
	}
}

// TestTremorOfIArathael_BanishGrantsPlusTwo: CardBanished set flips +2{p}.
func TestTremorOfIArathael_BanishGrantsPlusTwo(t *testing.T) {
	cases := []struct {
		c    card.Card
		base int
	}{
		{cards.TremorOfIArathaelRed{}, 4},
		{cards.TremorOfIArathaelYellow{}, 3},
		{cards.TremorOfIArathaelBlue{}, 2},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardBanished(true).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if want := tc.base + 2; ge.Value() != want {
			t.Errorf("%s: Value = %d, want %d (CardBanished set, +2{p} rider on)", tc.c.Name(), ge.Value(), want)
		}
	}
}
