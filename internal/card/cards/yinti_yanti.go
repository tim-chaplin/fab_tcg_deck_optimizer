// Yinti Yanti: "While Yinti Yanti is attacking and you control an aura, it has +1{p}.
// While Yinti Yanti is defending and you control an aura, it has +1{d}." Both bonuses
// gate on ge.AuraCount() > 0 — any aura type qualifies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// yintiYantiBonus returns +1 when any aura is in play, else 0.
func yintiYantiBonus(ge card.GameEngine) int {
	if ge.AuraCount() > 0 {
		return 1
	}
	return 0
}

func yintiYantiPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += yintiYantiBonus(ge)
}

func yintiYantiBlock(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusDefense += yintiYantiBonus(ge)
}

func (YintiYantiRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(ge, l, self)
}
func (YintiYantiRed) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(ge, l, self)
}

func (YintiYantiYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(ge, l, self)
}
func (YintiYantiYellow) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(ge, l, self)
}

func (YintiYantiBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(ge, l, self)
}
func (YintiYantiBlue) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(ge, l, self)
}
