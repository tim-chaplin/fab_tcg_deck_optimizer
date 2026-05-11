package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

var whisperOfTheOracleVariants = []sim.Card{
	WhisperOfTheOracleRed{},
	WhisperOfTheOracleYellow{},
	WhisperOfTheOracleBlue{},
}

// Tests that every variant emits a LogPlay step and an Opt 4 log entry.
func TestWhisperOfTheOracle_PlayCallsOpt4(t *testing.T) {
	a, b, c, d := testutils.NewStubCard("a"), testutils.NewStubCard("b"),
		testutils.NewStubCard("c"), testutils.NewStubCard("d")
	defer testutils.SwapCurrentHero(testutils.Hero{})()

	for _, variant := range whisperOfTheOracleVariants {
		s := sim.NewTurnStateFromCards([]sim.Card{a, b, c, d}, nil)
		sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: variant})
		if s.Value() != 0 {
			t.Errorf("%s: Play() Value = %d, want 0", variant.Name(), s.Value())
		}
		if len(s.LogEntries()) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (Opted... + chain step)", variant.Name(), len(s.LogEntries()))
			continue
		}
		// Play emits the Opted line; the chain step is auto-appended after Play
		// returns, so the Opted entry lands first.
		want := "Opted [a, b, c, d], put [a, b, c, d] on top, put [] on bottom"
		if got := s.LogEntries()[0].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", variant.Name(), got, want)
		}
	}
}

// Tests that every variant carries Go again so the chain runner can keep playing.
func TestWhisperOfTheOracle_GoAgain(t *testing.T) {
	for _, c := range whisperOfTheOracleVariants {
		if !c.GoAgain() {
			t.Errorf("%s: GoAgain() = false, want true", c.Name())
		}
	}
}
