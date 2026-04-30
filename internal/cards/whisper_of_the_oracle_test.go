package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Whisper of the Oracle resolves with no chain-step credit (the Opt 4 reorder is the
// only on-play effect and it isn't modelled).
func TestWhisperOfTheOracle_PlayCreditsNothing(t *testing.T) {
	for _, c := range []sim.Card{WhisperOfTheOracleRed{}, WhisperOfTheOracleYellow{}, WhisperOfTheOracleBlue{}} {
		var s sim.TurnState
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() Value = %d, want 0", c.Name(), got)
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
