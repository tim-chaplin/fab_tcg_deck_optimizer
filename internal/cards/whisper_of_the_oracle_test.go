package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
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
	prev := sim.CurrentHero
	sim.CurrentHero = testutils.Hero{}
	defer func() { sim.CurrentHero = prev }()

	for _, card := range whisperOfTheOracleVariants {
		s := sim.NewTurnState([]sim.Card{a, b, c, d}, nil)
		card.Play(s, &sim.CardState{Card: card})
		if s.Value != 0 {
			t.Errorf("%s: Play() Value = %d, want 0", card.Name(), s.Value)
		}
		if len(s.Log) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (LogPlay + Opted ...)", card.Name(), len(s.Log))
			continue
		}
		want := "Opted [a, b, c, d], put [a, b, c, d] on top, put [] on bottom"
		if got := s.Log[1].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", card.Name(), got, want)
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
