package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that every variant credits Opt 4 at 4 * sim.OptValue.
func TestWhisperOfTheOracle_PlayCreditsOpt4(t *testing.T) {
	want := 4 * sim.OptValue
	for _, c := range []sim.Card{WhisperOfTheOracleRed{}, WhisperOfTheOracleYellow{}, WhisperOfTheOracleBlue{}} {
		var s sim.TurnState
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != want {
			t.Errorf("%s: Play() Value = %d, want %d", c.Name(), got, want)
		}
	}
}

// Tests that every variant carries Go again so the chain runner can keep playing.
func TestWhisperOfTheOracle_GoAgain(t *testing.T) {
	for _, c := range []sim.Card{WhisperOfTheOracleRed{}, WhisperOfTheOracleYellow{}, WhisperOfTheOracleBlue{}} {
		if !c.GoAgain() {
			t.Errorf("%s: GoAgain() = false, want true", c.Name())
		}
	}
}
