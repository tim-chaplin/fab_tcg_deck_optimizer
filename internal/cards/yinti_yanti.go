// Yinti Yanti: "While Yinti Yanti is attacking and you control an aura, it has +1{p}.
// While Yinti Yanti is defending and you control an aura, it has +1{d}." Both bonuses
// gate on len(s.Auras) > 0 — any aura type qualifies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// yintiYantiBonus returns +1 when any aura is in play, else 0.
func yintiYantiBonus(s *sim.TurnState) int {
	if len(s.Auras) > 0 {
		return 1
	}
	return 0
}

func yintiYantiPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += yintiYantiBonus(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func yintiYantiBlock(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusDefense += yintiYantiBonus(s)
}

func (YintiYantiRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	yintiYantiPlay(s, l, self)
}
func (YintiYantiRed) Block(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	yintiYantiBlock(s, l, self)
}

func (YintiYantiYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	yintiYantiPlay(s, l, self)
}
func (YintiYantiYellow) Block(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	yintiYantiBlock(s, l, self)
}

func (YintiYantiBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	yintiYantiPlay(s, l, self)
}
func (YintiYantiBlue) Block(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	yintiYantiBlock(s, l, self)
}
