package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

var whisperOfTheOracleVariants = []card.Card{
	cards.WhisperOfTheOracleRed{},
	cards.WhisperOfTheOracleYellow{},
	cards.WhisperOfTheOracleBlue{},
}

// Tests that every variant returns Value 0 (the Opt 4 rider reshapes the deck, doesn't
// credit damage).
func TestWhisperOfTheOracle_PlayCallsOpt4(t *testing.T) {
	a, b, c, d := testutils.NewStubCard("a"), testutils.NewStubCard("b"),
		testutils.NewStubCard("c"), testutils.NewStubCard("d")

	for _, variant := range whisperOfTheOracleVariants {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b, c, d}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: variant})
		if ge.Value() != 0 {
			t.Errorf("%s: Play() Value = %d, want 0", variant.Name(), ge.Value())
		}
	}
}

// Tests that every variant carries Go again so the chain runner can keep playing.
func TestWhisperOfTheOracle_GoAgain(t *testing.T) {
	for _, c := range whisperOfTheOracleVariants {
		if !c.GoAgain(nil) {
			t.Errorf("%s: GoAgain() = false, want true", c.Name())
		}
	}
}
