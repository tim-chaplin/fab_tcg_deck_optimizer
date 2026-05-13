// Whisper of the Oracle — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 3.
//
// Text: "**Opt 4** **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func whisperOfTheOraclePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.Opt(l, 4)
}

func (WhisperOfTheOracleRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(ge, l, self)
}

func (WhisperOfTheOracleYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(ge, l, self)
}

func (WhisperOfTheOracleBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(ge, l, self)
}
