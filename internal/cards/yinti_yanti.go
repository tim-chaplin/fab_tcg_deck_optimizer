// Yinti Yanti: "While Yinti Yanti is attacking and you control an aura, it has +1{p}.
// While Yinti Yanti is defending and you control an aura, it has +1{d}." Both bonuses
// gate on g.AuraCount() > 0 — any aura type qualifies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// yintiYantiBonus returns +1 when any aura is in play, else 0.
func yintiYantiBonus(g card.GameEngine) int {
	if g.AuraCount() > 0 {
		return 1
	}
	return 0
}

func yintiYantiPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += yintiYantiBonus(g)
}

func yintiYantiBlock(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusDefense += yintiYantiBonus(g)
}

func (YintiYantiRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(g, l, self)
}
func (YintiYantiRed) Block(g card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(g, l, self)
}

func (YintiYantiYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(g, l, self)
}
func (YintiYantiYellow) Block(g card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(g, l, self)
}

func (YintiYantiBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiPlay(g, l, self)
}
func (YintiYantiBlue) Block(g card.GameEngine, l card.Logger, self *card.CardState) {
	yintiYantiBlock(g, l, self)
}
