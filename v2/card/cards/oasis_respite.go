// Oasis Respite — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "Prevent the next 4 damage that would be dealt to target hero this turn by
// a source of your choice. If they have less life than each other hero, they may gain
// 1{h}." Yellow caps at 3, Blue at 2.
//
// The 1{h} life-gain rider fires for heroes opting into card.LowerHealthWanter; it folds
// into the chain-step "(+N)" by bumping self.BonusDefense so the sim's resolver caps it
// against IncomingDamage alongside the printed prevention.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func oasisRespitePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.HeroWantsLowerHealth() {
		self.BonusDefense += 1
	}
}

func (OasisRespiteRed) DefensiveInstant() {}
func (OasisRespiteRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	oasisRespitePlay(ge, l, self)
}

func (OasisRespiteYellow) DefensiveInstant() {}
func (OasisRespiteYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	oasisRespitePlay(ge, l, self)
}

func (OasisRespiteBlue) DefensiveInstant() {}
func (OasisRespiteBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	oasisRespitePlay(ge, l, self)
}
