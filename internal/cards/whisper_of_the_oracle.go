// Whisper of the Oracle — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 3.
//
// Text: "**Opt 4** **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func whisperOfTheOraclePlay(s *sim.TurnState, self *sim.CardState) {
	s.Log(self, 0)
	s.Opt(4)
}

func (WhisperOfTheOracleRed) Play(s *sim.TurnState, self *sim.CardState) {
	whisperOfTheOraclePlay(s, self)
}

func (WhisperOfTheOracleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	whisperOfTheOraclePlay(s, self)
}

func (WhisperOfTheOracleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	whisperOfTheOraclePlay(s, self)
}
