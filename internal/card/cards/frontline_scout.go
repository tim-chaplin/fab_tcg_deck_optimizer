// Frontline Scout — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may look at the defending hero's hand. If Frontline Scout is played from arsenal, it
// gains **go again**."
//
// not implemented: opposing-hand-peek rider — the optimizer plays only our hero, so peeking the
// defender's hand yields no tempo or damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func frontlineScoutPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
}

func (FrontlineScoutRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	frontlineScoutPlay(ge, l, self)
}

func (FrontlineScoutYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	frontlineScoutPlay(ge, l, self)
}

func (FrontlineScoutBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	frontlineScoutPlay(ge, l, self)
}
