// Whisper of the Oracle — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 3.
//
// Text: "**Opt 4** **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func whisperOfTheOraclePlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.Opt(l, 4)
}

func (WhisperOfTheOracleRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(s, l, self)
}

func (WhisperOfTheOracleYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(s, l, self)
}

func (WhisperOfTheOracleBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(s, l, self)
}
