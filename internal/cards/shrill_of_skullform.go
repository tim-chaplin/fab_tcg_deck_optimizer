// Shrill of Skullform — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Shrill of Skullform gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (ShrillOfSkullformRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	shrillPlay(ge, l, self)
}

func (ShrillOfSkullformYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	shrillPlay(ge, l, self)
}

func (ShrillOfSkullformBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	shrillPlay(ge, l, self)
}

// shrillPlay routes the +3{p} aura-in-play buff through self.BonusAttack.
func shrillPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.AuraCreated() {
		self.BonusAttack += 3
	}
}
