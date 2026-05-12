// Seek Horizon — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Seek Horizon, you may put a card from your hand on top of
// your deck. If you do, Seek Horizon gains **go again**."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// not implemented: hand-on-top alt cost and conditional go-again rider

func (c SeekHorizonRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: hand-on-top alt cost and conditional go-again rider

func (c SeekHorizonYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: hand-on-top alt cost and conditional go-again rider

func (c SeekHorizonBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
