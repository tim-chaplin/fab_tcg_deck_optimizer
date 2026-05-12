// Frontline Scout — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may look at the defending hero's hand. If Frontline Scout is played from arsenal, it
// gains **go again**."
//
// Modelling: hand-peek isn't modelled.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func frontlineScoutPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
}

// not implemented: opposing-hand-peek rider

func (FrontlineScoutRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	frontlineScoutPlay(g, l, self)
}

// not implemented: opposing-hand-peek rider

func (FrontlineScoutYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	frontlineScoutPlay(g, l, self)
}

// not implemented: opposing-hand-peek rider

func (FrontlineScoutBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	frontlineScoutPlay(g, l, self)
}
