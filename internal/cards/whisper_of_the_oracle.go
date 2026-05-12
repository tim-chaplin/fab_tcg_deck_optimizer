// Whisper of the Oracle — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 3.
//
// Text: "**Opt 4** **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func whisperOfTheOraclePlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.Opt(l, 4)
}

func (WhisperOfTheOracleRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(g, l, self)
}

func (WhisperOfTheOracleYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(g, l, self)
}

func (WhisperOfTheOracleBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	whisperOfTheOraclePlay(g, l, self)
}
