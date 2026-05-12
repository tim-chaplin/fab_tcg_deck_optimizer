// Yinti Yanti: "While Yinti Yanti is attacking and you control an aura, it has +1{p}.
// While Yinti Yanti is defending and you control an aura, it has +1{d}." Both bonuses
// gate on s.AuraCount() > 0 — any aura type qualifies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// yintiYantiBonus returns +1 when any aura is in play, else 0.
func yintiYantiBonus(s card.GameEngine) int {
	if s.AuraCount() > 0 {
		return 1
	}
	return 0
}

func yintiYantiPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += yintiYantiBonus(s)
}

func yintiYantiBlock(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusDefense += yintiYantiBonus(s)
}

func (YintiYantiRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(s, l, self)
}
func (YintiYantiRed) Block(s card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(s, l, self)
}

func (YintiYantiYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(s, l, self)
}
func (YintiYantiYellow) Block(s card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(s, l, self)
}

func (YintiYantiBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(s, l, self)
}
func (YintiYantiBlue) Block(s card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(s, l, self)
}
