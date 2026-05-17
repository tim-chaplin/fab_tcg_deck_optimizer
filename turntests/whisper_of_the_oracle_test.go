package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

var whisperOfTheOracleVariants = []card.Card{
	cards.WhisperOfTheOracleRed{},
	cards.WhisperOfTheOracleYellow{},
	cards.WhisperOfTheOracleBlue{},
}

// Tests that every variant emits a LogPlay step and an Opt 4 log entry.
func TestWhisperOfTheOracle_PlayCallsOpt4(t *testing.T) {
	a, b, c, d := testutils.NewStubCard("a"), testutils.NewStubCard("b"),
		testutils.NewStubCard("c"), testutils.NewStubCard("d")

	for _, variant := range whisperOfTheOracleVariants {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b, c, d}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: variant})
		if ge.Value() != 0 {
			t.Errorf("%s: Play() Value = %d, want 0", variant.Name(), ge.Value())
		}
		if len(ge.LogEntries()) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (Opted... + chain step)", variant.Name(), len(ge.LogEntries()))
			continue
		}
		// Play emits the Opted line; the chain step is auto-appended after Play
		// returns, so the Opted entry lands first.
		want := "Opted [a, b, c, d], put [a, b, c, d] on top, put [] on bottom"
		if got := ge.LogEntries()[0].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", variant.Name(), got, want)
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
