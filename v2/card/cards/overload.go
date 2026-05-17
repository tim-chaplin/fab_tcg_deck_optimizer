// Overload — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Dominate** If Overload hits, it gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func overloadPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.LikelyToHit(self) {
		self.GrantedGoAgain = true
	}
}

func (OverloadRed) Dominate() {}
func (OverloadRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	overloadPlay(ge, l, self)
}

func (OverloadYellow) Dominate() {}
func (OverloadYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	overloadPlay(ge, l, self)
}

func (OverloadBlue) Dominate() {}
func (OverloadBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	overloadPlay(ge, l, self)
}
